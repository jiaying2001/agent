package store

import (
	"github.com/jiaying2001/agent/log"
	"gopkg.in/yaml.v3"
	"os"
	"runtime"
)

var (
	Pass      Identity
	FileStore Files
	OsI       OsInfo
	Ids       map[string]IdsType
	C         Config
)

type Identity struct {
	AuthToken string
	UserName  string
}

type Files struct {
	Content map[string]*File
}

type File struct {
	Offset int64
}

type OsInfo struct {
	Version string
}

type IdsType struct {
	Type string `json:"type"`
	Ref  int    `json:"ref"`
}

type Config struct {
	Server    ServerConfig
	Kafka     KafkaConfig
	Zookeeper ZookeeperConfig
}

type ServerConfig struct {
	Hostname string
	Port     string
}

type KafkaConfig struct {
	Topic    string
	HostName string
	Port     string
	Client   KafkaClientConfig
}

type KafkaClientConfig struct {
	Id string
}

type ZookeeperConfig struct {
	Hostname string
	Port     string
	Ids      IdsConfig
}

type IdsConfig struct {
	Config string
}

func init() {
	FileStore.Content = make(map[string]*File)
	OsI.Version = runtime.GOOS
	loadConfig("config.prod.yml")
}

func loadConfig(fileName string) {
	bytes, err := os.ReadFile(fileName)
	log.Logger.Info("Loading configuration file " + fileName)
	if err != nil {
		log.Logger.Error("Error reading configuration file: " + err.Error())
		os.Exit(1)
	}
	err = yaml.Unmarshal(bytes, &C)
	if err != nil {
		log.Logger.Error("Error parsing configuration file: " + err.Error())
		os.Exit(1)
	}
	log.Logger.Info("Loaded configuration file " + fileName)
}

func (f *Files) GetOffset(path string) int64 {
	if f.Content[path] == nil {
		f.Content[path] = &File{
			0,
		}
	}
	return f.Content[path].Offset
}

func (f *Files) SetOffset(path string, offset int64) {
	f.Content[path].Offset = offset
}
