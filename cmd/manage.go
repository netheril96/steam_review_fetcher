/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	"github.com/go-resty/resty/v2"
	steamreviewfetcher "github.com/netheril96/steam_review_fetcher/lib"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

// manageCmd represents the manage command
var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Manage an app",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var restyClient = resty.New()
		var limiter = rate.NewLimiter(1, 1)
		restyClient.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			return limiter.Wait(context.Background())
		})
		var manager = steamreviewfetcher.AppManager{
			AppId:     appId,
			Directory: saveDir,
			ApiClient: steamreviewfetcher.NewSteamApiClient(restyClient),
		}
		var err = manager.Init()
		if err != nil {
			log.Fatal(err)
		}
		err = manager.ResumeFetch()
		if err != nil {
			log.Fatal(err)
		}
	},
}
var saveDir string

func init() {
	rootCmd.AddCommand(manageCmd)
	manageCmd.Flags().IntVar(&appId, "id", 0, "App Id")
	manageCmd.MarkFlagRequired("id")
	manageCmd.Flags().StringVar(&saveDir, "dir", "", "Where to save the data")
	manageCmd.MarkFlagRequired("dir")
}
