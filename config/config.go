package config

import (
	"net/url"

	"github.com/spf13/cobra"
)

var (
	// HTTPEndpoint api endpoint
	HTTPEndpoint = "0.0.0.0:1234"
	// GRPCEndpoint api endpoint
	GRPCEndpoint = "0.0.0.0:1235"
	configURI    *url.URL
	loggerURI    *url.URL
	rootCmd      *cobra.Command
)

// SetConfigURI set configURI
func SetConfigURI(uri string) error {
	var err error
	configURI, err = url.Parse(uri)
	return err
}

// SetLoggerURI set loggerURI
func SetLoggerURI(uri string) error {
	var err error
	loggerURI, err = url.Parse(uri)
	return err
}

func init() {

}
