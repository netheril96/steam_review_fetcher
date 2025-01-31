package steamreviewfetcher

import (
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/goccy/go-json"
)

const APP_LIST_URL = "http://api.steampowered.com/ISteamApps/GetAppList/v2"

type app struct {
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}

type appList struct {
	Apps []app `json:"apps"`
}

type apiResponse struct {
	AppList appList `json:"applist"`
}

func ListAppIds(url string) ([]int, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var apiRes apiResponse
	err = json.Unmarshal(body, &apiRes)
	if err != nil {
		return nil, err
	}
	result := make([]int, len(apiRes.AppList.Apps))
	for i := range result {
		result[i] = apiRes.AppList.Apps[i].AppID
	}
	slices.Sort(result)
	return result, nil
}
