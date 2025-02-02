package steamreviewfetcher

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/goccy/go-json"
	"github.com/klauspost/compress/zstd"
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
	rawReader, err := os.Open(filepath.Join(p.Directory, "details.jsonz"))
	if err != nil {
		return err
	}
	defer rawReader.Close()
	reader, err := zstd.NewReader(rawReader)
	if err != nil {
		return err
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &p.game)
}

func (p *AppManager) saveGameDetails() error {
	rawWriter, err := os.CreateTemp(p.Directory, "details*.jsonz")
	if err != nil {
		return err
	}
	defer rawWriter.Close()
	writer, err := zstd.NewWriter(rawWriter)
	if err != nil {
		return err
	}
	defer writer.Close()
	data, err := json.Marshal(p.game)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	err = rawWriter.Close()
	if err != nil {
		return err
	}
	return os.Rename(rawWriter.Name(), filepath.Join(p.Directory, "details.jsonz"))
}
