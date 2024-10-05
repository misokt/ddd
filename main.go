package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const (
	MESSAGES_JSON_PATH string = "/messages.json"
	CHANNEL_JSON_PATH  string = "/channel.json"
	MESSAGES_PATH      string = "messages"
	DUMP_FILE          string = "messages.txt"
	FORMATINT_BASE     int    = 10
)

type MessagesJSON []struct {
	ID int64 `json:"ID,omitempty"`
}

type ChannelJSON struct {
	ID string `json:"id,omitempty"`
}

type FlagValues struct {
	Path             string
	AllMessages      bool
	ByYear           string
	ExcludeChannels  string
	SpecificChannels string
}

var (
	channelJson  ChannelJSON
	messagesJson MessagesJSON
	flagValues   FlagValues
	// wg       sync.WaitGroup
)

func parseFlags() {
	flag.StringVar(&flagValues.Path, "path", MESSAGES_PATH, "path to 'messages'")
	flag.BoolVar(&flagValues.AllMessages, "all", true, "dump every message")
	flag.StringVar(&flagValues.ByYear, "year", "", "dump every message from a specified year")
	flag.StringVar(&flagValues.ExcludeChannels, "exclude", "", "dump every message excluding specified channels")
	flag.StringVar(&flagValues.SpecificChannels, "specific", "", "dump every message only from specified channels")
	flag.Parse()
}

func pathWalk(dumpFile *os.File) error {
	defer dumpFile.Close()

	err := filepath.Walk(flagValues.Path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() != flagValues.Path && info.IsDir() {
			readDirs(path, dumpFile)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func readDirs(path string, dumpFile *os.File) {
	cFile, err := os.ReadFile(path + CHANNEL_JSON_PATH)
	if err != nil {
		log.Fatalf("error reading channel file. ERROR: %v\n", err)
	}
	if err := json.Unmarshal(cFile, &channelJson); err != nil {
		log.Fatalf("error unmarshalling channel json file. ERROR: %v\n", err)
	}

	mFile, err := os.ReadFile(path + MESSAGES_JSON_PATH)
	if err != nil {
		log.Fatalf("error reading messages file. ERROR: %v\n", err)
	}
	if err := json.Unmarshal(mFile, &messagesJson); err != nil {
		log.Fatalf("error unmarshalling messages json file. ERROR: %v\n", err)
	}

	messagesJsonLength := len(messagesJson)
	if messagesJsonLength == 0 {
		return
	}

	dumpToFile(dumpFile, (channelJson.ID + ":\n"))
	for i, m := range messagesJson {
		if i+1 == messagesJsonLength {
			dumpToFile(dumpFile, (strconv.FormatInt(m.ID, FORMATINT_BASE) + "\n\n"))
		} else {
			dumpToFile(dumpFile, (strconv.FormatInt(m.ID, FORMATINT_BASE) + ", "))
		}
	}
}

func dumpToFile(dumpFile *os.File, content string) {
	_, err := dumpFile.WriteString(content)
	if err != nil {
		log.Fatalf("error dumping into file. ERROR: %v\n", err)
	}
}

func createFile() (*os.File, error) {
	file, err := os.Create(DUMP_FILE)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func main() {
	parseFlags()

	dumpFile, err := createFile()
	if err != nil {
		log.Fatalf("error creating dump file. ERROR: %v\n", err)
	}

	if err := pathWalk(dumpFile); err != nil {
		log.Fatalf("error walking path. ERROR: %v\n", err)
	}

	fmt.Printf("dumped to '%s'\n", DUMP_FILE)
}
