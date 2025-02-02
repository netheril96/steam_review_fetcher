package steamreviewfetcher

import (
	"path/filepath"
	"time"
)

type AppManager struct {
	ApiClient *SteamApiClient
	Directory string

	game               Game
	appManagerMetadata appManagerMetadata
}

type appManagerMetadata struct {
	Cursor     string    `json:"cursor"`
	FinishTime time.Time `json:"finish_time"`
}

func (p *AppManager) loadGameDetails() error {
	filepath := (filepath.Join(p.Directory, "details.jsonz"))
	return ReadZstdJson(filepath, &p.game)
}

func (p *AppManager) saveGameDetails() error {
	filepath := (filepath.Join(p.Directory, "details.jsonz"))
	return WriteZstdJson(filepath, &p.game)
}
