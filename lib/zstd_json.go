package steamreviewfetcher

import (
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
	"github.com/klauspost/compress/zstd"
)

func ReadZstdJson(filename string, v interface{}) error {
	rawReader, err := os.Open(filename)
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
	return json.Unmarshal(data, v)
}

func WriteZstdJsonTemp(directory string, v interface{}) (filename string, err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	rawWriter, err := os.CreateTemp(directory, "*.jsonz")
	if err != nil {
		return "", err
	}
	defer rawWriter.Close()
	writer, err := zstd.NewWriter(rawWriter)
	if err != nil {
		return "", err
	}
	defer writer.Close()
	_, err = writer.Write(data)
	return rawWriter.Name(), err
}

func WriteZstdJson(filename string, v interface{}) error {
	var dir = filepath.Dir(filename)
	tmpname, err := WriteZstdJsonTemp(dir, v)
	if err != nil {
		return err
	}
	return os.Rename(tmpname, filename)
}
