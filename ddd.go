package main

import (
    "encoding/csv"
    "encoding/json"
    "flag"
    "fmt"
    "io/fs"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    MESSAGES_DIR_PATH  string = "Messages"
    MESSAGES_JSON_PATH string = "messages.json"
    CHANNEL_JSON_PATH  string = "channel.json"
    DUMP_FILE_CSV      string = "messages.csv"

    TIMESTAMP_LAYOUT string = "2006-01-02 15:04:05"
    BASE_10          int    = 10
)

type MessagesJSON []struct {
    ID        int64  `json:"ID,omitempty"`
    Timestamp string `json:"Timestamp,omitempty"`
}

type ChannelJSON struct {
    ID string `json:"id,omitempty"`
}

type FlagValues struct {
    Exclude     bool
    Include     bool
    AllMessages bool
    ByYear      string
    ByChannels  string
}

var (
    flagValues      FlagValues
    channelJson     ChannelJSON
    messagesJson    MessagesJSON
    channelsList    map[string]struct{} = make(map[string]struct{})
    messagesCSVData [][]string          = [][]string{
        {"channelid", "messageid"},
    }
)

func parseFlags() {
    flag.BoolVar(&flagValues.AllMessages, "all", true, "dump every message")
    flag.StringVar(&flagValues.ByYear, "year", "", "dump every message from a specified year")
    flag.StringVar(&flagValues.ByChannels, "channels", "", "channels to exclude or include from the dump. [comma,separate,the,input]")
    flag.BoolVar(&flagValues.Exclude, "exclude", false, "exclude specified channels from the dump")
    flag.BoolVar(&flagValues.Include, "include", false, "only include specified channels from the dump")
    flag.Parse()
}

func pathWalk(pathToWalk string) error {
    err := filepath.Walk(pathToWalk, func(path string, info fs.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() && info.Name() != pathToWalk {
            readDirs(path)
        }
        return nil
    })
    if err != nil {
        return err
    }

    return nil
}

func readDirs(path string) {
    cFile, err := os.ReadFile(filepath.Join(path, CHANNEL_JSON_PATH))
    if err != nil {
        log.Fatalln("ERROR: could not read channel file:", err)
    }
    if err := json.Unmarshal(cFile, &channelJson); err != nil {
        log.Fatalln("ERROR: could not unmarshal channel json file:", err)
    }

    mFile, err := os.ReadFile(filepath.Join(path, MESSAGES_JSON_PATH))
    if err != nil {
        log.Fatalln("ERROR: could not read messages file:", err)
    }
    if err := json.Unmarshal(mFile, &messagesJson); err != nil {
        log.Fatalln("ERROR: could not unmarshal messages json file:", err)
    }

    messagesJsonLength := len(messagesJson)
    if messagesJsonLength == 0 {
        return
    }

    if flagValues.ByYear != "" {
        dumpByYear()
    } else if flagValues.ByChannels != "" {
        dumpByChannels()
    } else {
        dumpAllMessages()
    }
}

func dumpAllMessages() {
    for _, m := range messagesJson {
        messagesCSVData = append(messagesCSVData, []string{channelJson.ID, strconv.FormatInt(m.ID, BASE_10)})
    }
}

func dumpByYear() {
    for _, m := range messagesJson {
        parsedTimestamp, err := time.Parse(TIMESTAMP_LAYOUT, m.Timestamp)
        if err != nil {
             log.Fatalln("ERROR: failed parsing timestamp:", err)
        }

        flagYearInt, err := strconv.Atoi(flagValues.ByYear)
        if err != nil {
            log.Fatalln("ERROR: invalid year input:", err)
        }

        if parsedTimestamp.Year() != flagYearInt {
            continue
        }

        messagesCSVData = append(messagesCSVData, []string{channelJson.ID, strconv.FormatInt(m.ID, BASE_10)})
    }
}

func dumpByChannels() {
    _, channelIDExists := channelsList[channelJson.ID]
    if flagValues.Exclude && channelIDExists {
        return
    }
    if flagValues.Include && !channelIDExists {
        return
    }

    dumpAllMessages()
}

func dumpToFile(dumpFile *os.File) {
    writer := csv.NewWriter(dumpFile)
    if err := writer.WriteAll(messagesCSVData); err != nil {
        log.Fatalln("ERROR: failed writing to file:", err)
    }
}

func createFile() (*os.File, error) {
    file, err := os.Create(DUMP_FILE_CSV)
    if err != nil {
        return nil, err
    }
    return file, nil
}

func fillChannelsList() {
    if flagValues.ByChannels != "" {
        excludedChannels := strings.Split(flagValues.ByChannels, ",")
        for _, e := range excludedChannels {
            channelsList[e] = struct{}{}
        }
    }
}

func main() {
    log.SetFlags(log.Lshortfile)

    parseFlags()
    fillChannelsList()

    dumpFile, err := createFile()
    if err != nil {
        log.Fatalln("ERROR: failed creating dump file:", err)
    }
    defer dumpFile.Close()

    if err := pathWalk(MESSAGES_DIR_PATH); err != nil {
        log.Fatalf("ERROR: failed to recursively visit directories at '%s': %v\n", MESSAGES_DIR_PATH, err)
    }

    // Exit if CSV only has headers, i.e. no message was found to export
    if len(messagesCSVData) == 1 {
        fmt.Println("INFO: Didn't find any message to export")
        return
    }

    dumpToFile(dumpFile)

    fmt.Printf("INFO: Messages dumped to '%s'\n", dumpFile.Name())
}
