package steamreviewfetcher

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func createTestServer(testResponseFile string) *httptest.Server {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Could not get current filename")
	}
	testDir := filepath.Dir(filename)
	testFilePath := filepath.Join(testDir, "testdata", testResponseFile)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a successful response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		file, err := os.Open(testFilePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		io.Copy(w, file)
	})

	return httptest.NewServer(handler)
}

func TestSteamApiClient_ListAppIds(t *testing.T) {
	server := createTestServer("applist.json")
	defer server.Close()

	var client = SteamApiClient{httpClient: resty.New(), appListUrl: server.URL}

	appListResult, err := client.ListAppIds()
	if err != nil {
		t.Fatal(err.Error())
	}
	require.Contains(t, appListResult, 1835850)
	require.Contains(t, appListResult, 1835930)
	require.Contains(t, appListResult, 1835570)
	require.NotContains(t, appListResult, 1111)
}

func TestSteamApiClient_QueryAppDetails(t *testing.T) {
	server := createTestServer("appdetails.json")
	defer server.Close()

	var client = SteamApiClient{httpClient: resty.New(), appDetailsUrl: server.URL}

	appListResult, err := client.QueryAppDetails(1997660)
	require.Nil(t, err)
	require.Contains(t, string(appListResult), "Early Access")
}
