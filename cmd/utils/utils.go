package utils

import (
	"cmdboard/typefile"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const STORED_FILE_NAME = ".cmdboard.json"

func LoadCommands() (commands map[int]typefile.Command, err error) {
	bytes, err := ioutil.ReadFile(StoredFilePath())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := json.Unmarshal(bytes, &commands); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return commands, nil
}

func WriteCommand(commands map[int]typefile.Command) error {
	commandsJson, err := json.Marshal(commands)
	if err != nil {
		return err
	}
	ioutil.WriteFile(StoredFilePath(), commandsJson, 0664)
	return nil
}

func StoredFilePath() string {
	return GetFilePath(STORED_FILE_NAME)
}

func GetFilePath(fileName string) string {
	dir := StoredFileDir()
	return filepath.Join(dir, fileName)
}

func StoredFileDir() (filePath string) {
	if os.Getenv("CMDBOARD_STORED_FILE_PATH") != "" {
		filePath = os.Getenv("CMDBOARD_STORED_FILE_PATH")
	} else if os.Getenv("HOME") != "" {
		filePath = os.Getenv("HOME")
	} else if runtime.GOOS == "windows" {
		filePath = os.Getenv("APPDATA")
	} else {
		log.Fatal("Can't find work dir")
	}
	return
}
