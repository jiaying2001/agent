package myzk

import (
	"encoding/json"
	"github.com/jiaying2001/agent/harvester"
	"github.com/jiaying2001/agent/launcher"
	"github.com/jiaying2001/agent/log"
	"github.com/jiaying2001/agent/store"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"time"
)

var conn *zk.Conn

func init() {
	// 创建zk连接地址
	hosts := []string{store.C.Zookeeper.Hostname + `:` + store.C.Zookeeper.Port}
	// 连接zk
	var err error
	conn, _, err = zk.Connect(hosts, time.Second*5)
	if err != nil {
		log.Logger.Error("Error connecting to zookeeper" + err.Error())
		os.Exit(1)
	}
}

func Listen(path string) {
	exists, _, err := conn.Exists(path)
	if err != nil {
		log.Logger.Error("Error checking if a zookeeper node exists")
		os.Exit(1)
	}

	if exists {
		// If exist delete the node
		_, dStat, _ := conn.Get(path)
		err = conn.Delete(path, dStat.Version)
		if err != nil {
			log.Logger.Error("Error deleting a zookeeper node: " + err.Error())
			os.Exit(1)
		}
	}

	// Create the node
	data := []byte("")
	_, err = conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		os.Exit(1)
	}

	go func() {
		for {
			data, _, ch, err := conn.GetW(path)
			if len(data) == 0 {
				continue
			}
			log.Logger.Info("Node data at node " + path + ": " + string(data))
			if err != nil {
				log.Logger.Error("Error watching children for path " + path + " " + err.Error())
				return
			}
			// Wait for changes
			event := <-ch
			handleEvent(event, data)
		}
	}()
}

func LoadIdsNodes() {
	path := store.C.Zookeeper.Ids.Config
	data, _, err := conn.Get(path)
	if len(data) == 0 {
		return
	}
	log.Logger.Info("Node data at node " + path + ": " + string(data))
	if err != nil {
		log.Logger.Error("Error watching children for path " + path + " " + err.Error())
		return
	}
	err = json.Unmarshal(data, &store.Ids)
	if err != nil {
		log.Logger.Error("Error parsing: " + err.Error())
		return
	}
}

func ListenIdsNodes() {
	go func() {
		path := store.C.Zookeeper.Ids.Config
		for {
			data, _, ch, err := conn.GetW(path)
			if len(data) == 0 {
				continue
			}
			log.Logger.Info("Node data at node " + path + ": " + string(data))
			if err != nil {
				log.Logger.Error("Error watching children for path " + path + " " + err.Error())
				return
			}
			// Wait for changes
			event := <-ch
			handleIdsConfigsEven(event, data)
		}
	}()
}

func handleIdsConfigsEven(event zk.Event, data []byte) {
	switch event.Type {
	case zk.EventNodeDataChanged:
		log.Logger.Info("Receive event type: " + event.Type.String())
		err := json.Unmarshal(data, &store.Ids)
		if err != nil {
			log.Logger.Error("Error parsing: " + err.Error())
			return
		}
	default:
		log.Logger.Info("Event Not Handled: " + event.Type.String())
	}
}

func handleEvent(event zk.Event, data []byte) {
	switch event.Type {
	case zk.EventNodeDataChanged:
		log.Logger.Info("Receive event type: " + event.Type.String())
		var configs []harvester.Harvester
		err := json.Unmarshal(data, &configs)
		if err != nil {
			log.Logger.Error("Error parsing: " + err.Error())
			return
		}
		launcher.L.Load(&configs)
		launcher.L.StartWorkers()
	default:
		log.Logger.Info("Event Not Handled: " + event.Type.String())
	}
}
