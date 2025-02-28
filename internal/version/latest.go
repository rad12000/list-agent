package version

import "github.com/rad12000/list-agent/internal/github"

func Latest() string {
	const unknown = "UNKNOWN"

	ghService := github.NewService("")
	rls, err := ghService.GetLatestRelease()
	if err != nil {
		return unknown
	}

	return rls.TagName
}
