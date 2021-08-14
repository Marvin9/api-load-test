package cmd

import (
	"fmt"
	"os"

	"github.com/Marvin9/api-load-test/pkg"
	"github.com/spf13/cobra"
)

func Execute() {
	var session pkg.Session

	var rootCmd = &cobra.Command{
		Use:   "loadtest",
		Short: "Simple and elegant load testing tool for your API",
		Run: func(cmd *cobra.Command, args []string) {
			metadata := session.GenerateMetadata()

			session.LoadTest(metadata)
			session.Success()

			pkg.Analysis(session.Data).Display()
		},
	}
	rootCmd.Flags().StringVarP(&session.TargetEndpoint, "endpoint", "e", "", "target endpoint. eg: http://locahost:8000/")
	rootCmd.MarkFlagRequired("endpoint")
	rootCmd.Flags().StringVarP(&session.Method, "method", "m", "GET", "method of target endpoint [GET/POST/PUT...]")
	rootCmd.Flags().IntVarP(&session.Rate, "rate", "r", 100, "load of requets per second")
	rootCmd.Flags().IntVarP(&session.Until, "until", "u", 10, "duration of load in seconds")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
