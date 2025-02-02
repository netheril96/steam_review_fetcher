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
	"github.com/vbauerster/mpb/v8"
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
		var progressContainer = mpb.New()
		var bar = progressContainer.AddBar(int64(len(appIds)))
		for appid := range appIds {
			func() {
				defer bar.Increment()
				var directory = filepath.Join(saveDir, fmt.Sprintf("app.%d", appid))
				var err = os.MkdirAll(directory, 0755)
				if err != nil {
					log.Fatalf("Failed to make directory: %v", err)
				}
				var manager = steamreviewfetcher.AppManager{ApiClient: apiClient, Directory: directory, AppId: appid}
				err = manager.Init()
				if err != nil {
					log.Printf("Failed to init manager for app %d: %v", appid, err)
					return
				}
				if manager.ShouldSkip() {
					return
				}
				err = manager.ResumeFetch()
				if err != nil {
					log.Printf("Failed to fetch for app %d: %v", appid, err)
					return
				}
			}()
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchallCmd)
	fetchallCmd.Flags().StringVar(&saveDir, "dir", "", "Where to save the data")
	fetchallCmd.MarkFlagRequired("dir")
}
