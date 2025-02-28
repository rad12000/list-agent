package upgrade

import (
	"errors"
	"fmt"
	"github.com/rad12000/list-agent/internal/file"
	"github.com/rad12000/list-agent/internal/github"
	"github.com/rad12000/list-agent/internal/github/dto"
	"github.com/rad12000/list-agent/internal/version"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

type releaseService interface {
	DownloadAsset(asset dto.Asset, outputFile *os.File) error
	GetLatestRelease() (dto.ReleaseResponse, error)
	GetReleaseByTagName(tagName string) (dto.ReleaseResponse, error)
}

func run() error {
	var rls dto.ReleaseResponse
	var err error

	rlsService := github.NewService("")
	if flagVersion == "" {
		// Get Latest Release
		rls, err = rlsService.GetLatestRelease()
	} else {
		// Get release by tag
		rls, err = rlsService.GetReleaseByTagName(cleanTagVersion(flagVersion))
	}

	if err != nil {
		return err
	}

	// Download Release
	tarFile, err := downloadOSSpecificAsset(rls, rlsService)
	if err != nil {
		return err
	}

	// Navigate to temp dir
	tmpDir := filepath.Dir(tarFile.Name())
	defer func(file *os.File) {
		_ = os.RemoveAll(tmpDir)
	}(tarFile)

	if err := os.Chdir(tmpDir); err != nil {
		return err
	}

	// Extract tar file
	execFilename, err := extractReleaseTARFile(tarFile.Name())
	if err != nil {
		return err
	}

	// Move Executable to $PATH
	ex, err := os.Executable()
	if err != nil {
		return err
	}

	if err := os.Rename(execFilename, ex); err != nil {
		return err
	}

	fmt.Println("updated successfully!")
	currentVer := version.Version()
	fmt.Println(fmt.Sprintf("%s -> %s", currentVer, rls.TagName))

	return nil
}

func extractReleaseTARFile(tarFile string) (string, error) {
	files, err := file.DecompressTAR(tarFile)
	if err != nil {
		return "", err
	}

	execFilename := files[slices.IndexFunc(files, func(s string) bool {
		const listAgent = "listagent"
		return strings.HasSuffix(s, listAgent)
	})]

	// Grant Execute Permissions
	return execFilename, os.Chmod(execFilename, 0755)
}

// downloadOSSpecificAsset downloads the appropriate release asset for the OS and returns the tar file.
// NOTE: the *os.File will be closed and is unfit for I/O.
func downloadOSSpecificAsset(rls dto.ReleaseResponse, releaseService releaseService) (*os.File, error) {
	releaseFile := fmt.Sprintf("listagent_%s_%s", runtime.GOOS, runtime.GOARCH)
	assetToDownload, err := getAssetToDownload(releaseFile, rls.Assets)
	if err != nil {
		return nil, err
	}

	assetDir := filepath.Join(os.TempDir(), "listagent")
	if err := os.MkdirAll(assetDir, 0755); err != nil {
		return nil, err
	}

	tarFile, err := os.CreateTemp(assetDir, assetToDownload.Name)
	if err != nil {
		return nil, err
	}

	defer tarFile.Close()

	if err = releaseService.DownloadAsset(assetToDownload, tarFile); err != nil {
		return nil, err
	}

	return tarFile, nil
}

func getAssetToDownload(targetFileName string, assets []dto.Asset) (dto.Asset, error) {
	fullName := targetFileName + ".tgz"
	assetIndex := slices.IndexFunc(assets, func(asset dto.Asset) bool {
		return asset.Name == fullName
	})

	if assetIndex == -1 {
		return dto.Asset{}, errors.New(fmt.Sprintf("could not find a release artifact named %s", fullName))
	}

	return assets[assetIndex], nil
}

func cleanTagVersion(tag string) string {
	return "v" + strings.Trim(tag, "v")
}
