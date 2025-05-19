package memory

import (
	"os"
	"path/filepath"
)

func GetMemoriesFilePath(prefix string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	confDir := filepath.Join(homeDir, prefix)

	err = os.MkdirAll(confDir, 0755)
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(confDir, "memories.bin")
	return fp
}
