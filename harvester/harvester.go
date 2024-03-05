package harvester

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jiaying2001/agent/dto"
	"github.com/jiaying2001/agent/kafka"
	"github.com/jiaying2001/agent/log"
	"github.com/jiaying2001/agent/parser"
	"github.com/jiaying2001/agent/store"
	"github.com/satori/go.uuid"
	"os"
	"strconv"
	"sync"
	"time"
)

type Harvester struct {
	Offset     int64  `json:"offset"`
	Path       string `json:"path"`
	FileFormat string `json:"file_format"`
	Uuid       string
	Shutdown   *sync.WaitGroup
	Interrupt  bool
	State      int
	UserId     int `json:"user_id"`
}

const (
	Created = iota
	Running
)

func CreateHarvester(fileName string, fileFormat string) *Harvester {
	u := uuid.NewV4()
	return &Harvester{
		store.FileStore.GetOffset(fileName),
		fileName,
		fileFormat,
		u.String(),
		&sync.WaitGroup{},
		false,
		Created,
		store.Pass.UserID,
	}
}

func (h *Harvester) read() *[]string {
	// 打开文件
	file, err := os.Open(h.Path)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}(file)

	_, err = file.Seek(h.Offset, 0)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// 创建扫描器
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	// 逐行读取文件内容
	for scanner.Scan() {
		line := scanner.Text()
		h.Offset += int64(len(line) + 1) // 加上换行符的长度
		lines = append(lines, line)
	}

	// 检查扫描时是否出现错误
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}
	return &lines
}

func (h *Harvester) HandleInterrupt() {
	log.Logger.Info("Harvester for path " + h.Path + " received an interrupt")
	h.Interrupt = false
	// Save all the config
	log.Logger.Info("Saving offset " + strconv.FormatInt(h.Offset, 10) + " to " + h.Path)
	store.FileStore.Content[h.Path] = &store.File{
		Offset: h.Offset,
	}
}

func (h *Harvester) Run(parser parser.Parser) {
	for {
		lines := h.read()
		for _, value := range *lines {
			appName, pid := parser.Parse(value)
			key1 := store.OsI.Version + `_` + appName + `_private`
			key2 := store.OsI.Version + `_` + appName + `_public`
			_, ok1 := store.Ids[key1]
			_, ok2 := store.Ids[key2]
			if !(ok1 || ok2) {
				log.Logger.Info("No ids node for " + key1 + " & " + key2)
				continue
			}
			msg := dto.Message{
				Header: dto.Header{
					TraceId:   uuid.NewV4().String(),
					AuthToken: store.Pass.AuthToken,
					AppName:   appName,
					Pid:       pid,
					Os:        store.OsI.Version,
					UserID:    store.Pass.UserID,
					Path:      h.Path,
				},
				Body: dto.Body{
					Content: value,
				},
				Property: dto.Property{},
			}
			bytes, _ := json.Marshal(msg)
			kafka.Send(store.C.Kafka.Topic, bytes)
		}
		time.Sleep(1 * time.Second)
	}
}

func (h *Harvester) RunOutputToFile() {
	// 打开文件
	file, err := os.OpenFile("output.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}(file)

	lines := h.read()
	for _, value := range *lines {
		_, err := file.WriteString(value + "\n")
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
