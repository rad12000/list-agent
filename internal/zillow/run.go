package zillow

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	baseRequest = request{
		SearchQueryState: searchState{
			IsEntirePlaceForRent: true,
			IsRoomForRent:        false,
			IsListVisible:        true,
			MapZoom:              8,
		},
		Wants: rWants{
			Cat1: []string{"listResults"},
			Cat2: []string{"total"},
		},
		RequestId:      21,
		IsDebugRequest: false,
	}

	searchRegex *regexp.Regexp
)

type RunData struct {
	SearchTerms []string
	FilePath    string
	MapBounds   Bounds
	FilterState FilterState
	UserAgent   string
}

func parseSearchTerms(terms []string) {
	if len(terms) == 0 {
		return
	}
	var sb strings.Builder
	sb.WriteString(`(?i)`)
	for i, term := range terms {
		if i > 0 {
			sb.WriteString(`|`)
		}

		sb.WriteByte('(')
		sb.WriteString(regexp.QuoteMeta(term))
		sb.WriteByte(')')
	}

	searchRegex = regexp.MustCompile(sb.String())
}

func Run(data RunData) {
	startListening(data.FilePath)
	parseSearchTerms(data.SearchTerms)
	baseRequest.SearchQueryState.MapBounds = data.MapBounds
	baseRequest.SearchQueryState.FilterState = data.FilterState
	page := 1

	for {
		reqBody, err := json.Marshal(Copy(baseRequest, func(r request) request {
			r.SearchQueryState.Pagination.CurrentPage = page
			return r
		}))
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest(http.MethodPut, "https://www.zillow.com/async-create-search-page-state", bytes.NewBuffer(reqBody))
		if err != nil {
			panic(err)
		}
		req.Header = getHeaders(data)

		slog.Debug("sending new page state request")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		if res.StatusCode > 299 || res.StatusCode < 200 {
			io.Copy(os.Stdout, res.Body)
			res.Body.Close()
			panic(fmt.Sprintf("got status code %d", res.StatusCode))
		}

		reader := res.Body
		if res.Header.Get("Content-Encoding") == "gzip" {
			reader, err = gzip.NewReader(res.Body)
			if err != nil {
				panic(err)
			}
		}

		var resp Response
		if err = json.NewDecoder(reader).Decode(&resp); err != nil {
			panic(err)
		}
		_ = reader.Close()
		cookieJar.SetCookies(req.URL, res.Cookies())

		slog.Debug("handling results", "resultCount", len(resp.Cat1.SearchResults.ListResults))
		var wg sync.WaitGroup
		for _, result := range resp.Cat1.SearchResults.ListResults {
			wg.Add(1)
			go func() {
				defer wg.Done()
				handleResult(data, result.DetailUrl)
			}()
		}
		wg.Wait()

		if resp.Cat1.SearchList.Pagination.NextUrl != "" {
			page++
		} else {
			page = 0
			slog.Info("waiting three hours")
			doneCh := make(chan struct{})
			go func(uri string) {
				for {
					select {
					case <-doneCh:
						fmt.Println("exiting pinging")
						return
					default:
					}

					fmt.Println("pinging")
					req := Must(http.NewRequest(http.MethodGet, uri, nil))
					req.Header = getHeaders(data)
					res, err := http.DefaultClient.Do(req)
					if err != nil {
						slog.Error("failed to ping", "uri", uri, "err", err)
						return
					}
					if res.StatusCode > 299 || res.StatusCode < 200 {
						panic(fmt.Sprintf("got status code %d", res.StatusCode))
					}
					Must(io.Copy(io.Discard, res.Body))
					cobra.CheckErr(res.Body.Close())
					cookieJar.SetCookies(req.URL, res.Cookies())
					<-time.NewTimer(time.Minute * 3).C
				}
			}(resp.Cat1.SearchResults.ListResults[len(resp.Cat1.SearchResults.ListResults)-1].DetailUrl)
			<-time.NewTimer(time.Hour * 3).C
			continue
		}

		slog.Info("waiting a minute for next page", "page", page)
		<-time.NewTimer(time.Minute * 1).C
	}
}

func getHeaders(data RunData) map[string][]string {
	m := map[string]string{
		"Content-Type":    "application/json",
		"User-Agent":      data.UserAgent,
		"Accept":          "*/*",
		"Accept-Encoding": "gzip",
		"Cache-Control":   "no-cache",
	}

	out := make(map[string][]string)
	for k, v := range m {
		out[k] = []string{v}
	}

	return out
}

func handleResult(data RunData, uri string) (didSendRequest bool) {
	if hasBeenSeen(uri) {
		slog.Debug("skipping seen url", "uri", uri)
		return
	}

	if searchRegex == nil {
		markAsSeen(uri, true)
		slog.Info("new match", "url", uri)
		return
	}

	req := Must(http.NewRequest(http.MethodGet, uri, nil))
	req.Header = getHeaders(data)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("failed to handle result", "uri", uri, "err", err)
		return
	}
	if res.StatusCode > 299 || res.StatusCode < 200 {
		slog.Error("got error response for result", "uri", uri, "headers", req.Header)
		panic(fmt.Sprintf("got status code %d", res.StatusCode))
	}

	reader := res.Body
	if res.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			panic(err)
		}
	}
	defer func() {
		cobra.CheckErr(reader.Close())
	}()
	cookieJar.SetCookies(req.URL, res.Cookies())

	b, err := io.ReadAll(reader)
	if err != nil {
		slog.Error("failed to parse result response", "uri", uri, "err", err)
	}

	matched := searchRegex.Match(b)
	markAsSeen(uri, matched)
	if matched {
		slog.Info("new match", "url", uri)
	} else {
		slog.Debug("skipping non-match", "url", uri)
	}

	return true
}
