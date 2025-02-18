/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	steamreviewfetcher "github.com/netheril96/steam_review_fetcher/lib"
	"github.com/spf13/cobra"
)

// fetchreviewsCmd represents the fetchapp command
var fetchreviewsCmd = &cobra.Command{
	Use:   "fetchreviews",
	Short: "Fetch all reviews for a single app",
	Run: func(cmd *cobra.Command, args []string) {
		var httpClient = resty.New()
		httpClient.SetRetryCount(10)
		var client = steamreviewfetcher.NewSteamApiClient(httpClient)
		cursor := "*"
		for {
			raw, newCursor, err := client.QueryAppReview(appId, cursor)
			if err != nil {
				if errors.Is(err, &steamreviewfetcher.EndOfReview{}) {
					break
				} else {
					log.Fatalf("Failed to query app review %v", err)
				}
			}
			hash := sha256.Sum224([]byte(cursor))
			var filename = fmt.Sprintf("%d.%s.json", appId, hex.EncodeToString(hash[:]))
			err = os.WriteFile(filename, raw, 0644)
			if err != nil {
				log.Fatalf("Failed to write to %s: %v", filename, err)
			}
			cursor = newCursor
		}
	},
}
var appId int

func init() {
	rootCmd.AddCommand(fetchreviewsCmd)
	fetchreviewsCmd.Flags().IntVar(&appId, "id", 0, "The ID of the Steam app")
	fetchreviewsCmd.MarkFlagRequired("id")
}
