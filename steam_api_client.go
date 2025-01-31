package steamreviewfetcher

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-json"
)

type SteamApiClient struct {
	httpClient    *resty.Client
	appListUrl    string
	appDetailsUrl string
	appReviewUrl  string
}

func NewSteamApiClient(httpClient *resty.Client) *SteamApiClient {
	return &SteamApiClient{
		httpClient:    httpClient,
		appListUrl:    "http://api.steampowered.com/ISteamApps/GetAppList/v2",
		appDetailsUrl: "https://store.steampowered.com/api/appdetails",
		appReviewUrl:  "https://store.steampowered.com/appreviews",
	}
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

func (p *SteamApiClient) QueryAppDetails(appid int) (raw []byte, err error) {
	resp, err := p.httpClient.R().SetQueryParam("appids", strconv.Itoa(appid)).Get(p.appDetailsUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("wrong status %s", resp.Status())
	}

	type GameData struct {
		Success bool `json:"success"`
	}

	var data map[string]GameData
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, err
	}
	if data[strconv.Itoa(appid)].Success {
		return resp.Body(), nil
	}
	return nil, fmt.Errorf("the response is a failure:\n%s", string(resp.Body()))
}

func (p *SteamApiClient) QueryAppReview(appid int, cursor string) (raw []byte, newCursor string, err error) {
	resp, err := p.httpClient.R().
		SetQueryParam("json", "1").
		SetQueryParam("filter", "recent").
		SetQueryParam("num_per_page", "100").
		SetQueryParam("cursor", cursor).
		Get(fmt.Sprintf("%s/%d", p.appReviewUrl, appid))
	if err != nil {
		return
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("wrong status %s", resp.Status())
		return
	}

	type ReviewData struct {
		Success int    `json:"success"`
		Cursor  string `json:"cursor"`
	}
	var data ReviewData
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return
	}
	if data.Success == 0 {
		err = fmt.Errorf("the response is a failure:\n%s", string(resp.Body()))
		return
	}
	return resp.Body(), data.Cursor, nil
}
