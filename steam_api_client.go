package steamreviewfetcher

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-json"
)

type SteamApiClient struct {
	httpClient *resty.Client
	appListUrl string
}

func NewSteamApiClient(httpClient *resty.Client) *SteamApiClient {
	return &SteamApiClient{httpClient: httpClient, appListUrl: "http://api.steampowered.com/ISteamApps/GetAppList/v2"}
}

func (p *SteamApiClient) ListAppIds() ([]int, error) {
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

	resp, err := p.httpClient.R().Get(p.appListUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode())
	}
	var apiRes apiResponse
	err = json.Unmarshal(resp.Body(), &apiRes)
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
