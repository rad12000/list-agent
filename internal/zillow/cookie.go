package zillow

import (
	_ "embed"
	"github.com/spf13/cobra"
	"net/http"
	"net/http/cookiejar"
)

var (
	cookieJar, _ = cookiejar.New(nil)
)

func init() {
	http.DefaultClient.Jar = cookieJar
}

func Must[T any](t T, err error) T {
	cobra.CheckErr(err)
	return t
}
