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
	"github.com/stretchr/testify/assert"
)

func TestSteamApiClient_ListAppIds(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Could not get current filename")
	}
	testDir := filepath.Dir(filename)

	testFilePath := filepath.Join(testDir, "testdata", "applist.json")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a successful response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		file, err := os.Open(testFilePath)
		if err != nil {
			t.Fatal(err.Error())
		}
		defer file.Close()
		io.Copy(w, file)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	var client = SteamApiClient{httpClient: resty.New(), appListUrl: server.URL}

	appListResult, err := client.ListAppIds()
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Contains(t, appListResult, 1835850)
	assert.Contains(t, appListResult, 1835930)
	assert.Contains(t, appListResult, 1835570)
	assert.NotContains(t, appListResult, 1111)
}
