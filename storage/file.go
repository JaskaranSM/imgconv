package storage

import (
	"io/ioutil"
	"os"
	"path"
)

func NewFileStorage(baseDir string) *FileStorage {
	if baseDir == "" {
		baseDir = "."
	}
	return &FileStorage{
		baseDir: baseDir,
	}
}

type FileStorage struct {
	baseDir string
}

func (f *FileStorage) SaveImageBytesId(id string, data []byte) error {
	file, err := os.Create(path.Join(f.baseDir, id))
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil
}

func (f *FileStorage) GetImageBytesById(id string) ([]byte, error) {
	var out []byte
	var err error
	if out, err = ioutil.ReadFile(path.Join(f.baseDir, id)); err != nil {
		return nil, err
	}
	return out, nil
}
