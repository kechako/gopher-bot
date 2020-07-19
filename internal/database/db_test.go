package database

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func makeTestDir(name string) (string, func(), error) {
	dir, err := ioutil.TempDir(os.TempDir(), "gopher-bot-database-test")
	if err != nil {
		return "", nil, err
	}

	path := filepath.Join(dir, name)

	return path, func() {
		os.RemoveAll(dir)
	}, nil
}
