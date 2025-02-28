package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rad12000/list-agent/internal/github/dto"
	"io"
	"net/http"
	"os"
)

const (
	repoOwner = "rad12000"
	repoName  = "list-agent"
)

func NewService(accessToken string) Service {
	return Service{
		client: newHTTPClient(accessToken),
	}
}

type Service struct {
	client *http.Client
}

// DownloadAsset downloads the given asset and writes the data to the provided outputFile.
func (s Service) DownloadAsset(asset dto.Asset, outputFile *os.File) error {
	req, err := http.NewRequest(http.MethodGet, asset.Url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("accept", asset.ContentType)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		resBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(resBytes))
		return errors.New(fmt.Sprintf("failed to download asset. Got response code %d", resp.StatusCode))
	}

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// GetLatestRelease returns the latest release for the snap-cli repository.
func (s Service) GetLatestRelease() (dto.ReleaseResponse, error) {
	var result dto.ReleaseResponse

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/repos/%s/%s/releases/latest", repoOwner, repoName), nil)
	if err != nil {
		return result, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		resBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(resBytes))
		return result, errors.New(fmt.Sprintf("failed to get latest release. Got response code %d", resp.StatusCode))
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	return result, err
}

// GetReleaseByTagName gets the release matching the provided tag name.
func (s Service) GetReleaseByTagName(tagName string) (dto.ReleaseResponse, error) {
	var result dto.ReleaseResponse

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/repos/%s/%s/releases/tags/%s", repoOwner, repoName, tagName), nil)
	if err != nil {
		return result, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		resBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(resBytes))
		return result, errors.New(fmt.Sprintf("failed to get latest release. Got response code %d", resp.StatusCode))
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	return result, err
}
