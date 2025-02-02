package steamreviewfetcher

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/goccy/go-json"
)

type AppManager struct {
	AppId     int
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

func (p *AppManager) Init() error {
	err := ReadZstdJson(filepath.Join(p.Directory, "metadata.jsonz"), &p.appManagerMetadata)
	if errors.Is(err, fs.ErrNotExist) {
		p.appManagerMetadata = appManagerMetadata{
			Cursor:     "*",
			FinishTime: time.Time{},
		}
		err = WriteZstdJson(filepath.Join(p.Directory, "metadata.jsonz"), &p.appManagerMetadata)
	}
	if err != nil {
		return err
	}
	err = ReadZstdJson(filepath.Join(p.Directory, "details.jsonz"), &p.game)
	if errors.Is(err, fs.ErrNotExist) {
		var details []byte
		details, err = p.ApiClient.QueryAppDetails(p.AppId)
		if err != nil {
			return err
		}
		err = json.Unmarshal(details, &p.game)
		if err != nil {
			return err
		}
		err = WriteRawZstd(filepath.Join(p.Directory, "details.jsonz"), details)
	}
	return err
}

func (p *AppManager) Save() error {
	return WriteZstdJson(filepath.Join(p.Directory, "metadata.jsonz"), &p.appManagerMetadata)
}

func (p *AppManager) ResumeFetch() error {
	if !p.appManagerMetadata.FinishTime.IsZero() {
		return nil
	}
	raw, cursor, err := p.ApiClient.QueryAppReview(p.AppId, p.appManagerMetadata.Cursor)
	if err != nil {
		if errors.Is(err, &EndOfReview{}) {
			p.appManagerMetadata.FinishTime = time.Now()
			p.appManagerMetadata.Cursor = "*"
			return p.Save()
		}
		return err
	}
	cursorHash := sha256.Sum224([]byte(cursor))
	err = WriteRawZstd(filepath.Join(p.Directory, fmt.Sprintf("review.%s.jsonz", hex.EncodeToString(cursorHash[:]))), raw)
	if err != nil {
		return err
	}
	p.appManagerMetadata.Cursor = cursor
	return p.Save()
}
