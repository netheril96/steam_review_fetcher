/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	steamreviewfetcher "github.com/netheril96/steam_review_fetcher/lib"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

// fetchallCmd represents the fetchall command
var fetchallCmd = &cobra.Command{
	Use:   "fetchall",
	Short: "Fetch all apps",
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
		var apiClient = steamreviewfetcher.NewSteamApiClient(restyClient)
		var appIds, err = apiClient.ListAppIds()
		if err != nil {
			log.Fatalf("Failed to query app IDs: %v", err)
		}
		for appid := range appIds {
			var directory = filepath.Join(saveDir, fmt.Sprintf("app.%d", appid))
			err = os.MkdirAll(directory, 0755)
			if err != nil {
				log.Fatalf("Failed to make directory: %v", err)
			}
			var manager = steamreviewfetcher.AppManager{ApiClient: apiClient, Directory: directory, AppId: appid}
			err = manager.Init()
			if err != nil {
				continue
			}
			if manager.ShouldSkip() {
				continue
			}
			err = manager.ResumeFetch()
			if err != nil {
				continue
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchallCmd)
	fetchallCmd.Flags().StringVar(&saveDir, "dir", "", "Where to save the data")
	fetchallCmd.MarkFlagRequired("dir")
}
