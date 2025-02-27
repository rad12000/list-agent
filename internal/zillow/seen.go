package zillow

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	writeIgnoreCh = make(chan string)
	writeMatchCh  = make(chan string)
	readSeenCh    = make(chan pair[chan bool, string])
	persistence   *os.File
)

func startListening(file string) {
	persistence = Must(os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644))
	go listenSeen()
}

func listenSeen() {
	m := make(map[string]bool)
	reader := bufio.NewReader(persistence)
	for {
		line, _, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			panic(err)
		}

		parts := strings.Split(string(line), " ")
		m[parts[0]] = parts[1] == "true"
	}

	for {
		var err error
		select {
		case v := <-writeIgnoreCh:
			m[v] = false
			_, err = persistence.WriteString(fmt.Sprintf("%s false\n", v))
		case v := <-writeMatchCh:
			m[v] = true
			_, err = persistence.WriteString(fmt.Sprintf("%s true\n", v))
			openURL(v)
		case tuple := <-readSeenCh:
			_, ok := m[tuple.Second]
			tuple.First <- ok
		}

		if err != nil {
			panic(err)
		}
	}
}

func openURL(url string) {
	slog.Debug("Attempting to open url in browser", "url", url)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
	}

	if cmd == nil {
		return
	}

	if err := cmd.Run(); err != nil {
		slog.Error("failed to open url in default browser", "url", url, "err", err)
	}
}

func markAsSeen(uri string, matched bool) {
	if matched {
		writeMatchCh <- uri
	} else {
		writeIgnoreCh <- uri
	}
}

func hasBeenSeen(uri string) bool {
	ch := make(chan bool)
	readSeenCh <- pair[chan bool, string]{ch, uri}
	return <-ch
}
