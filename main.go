package main

import (
	"fmt"

	"github.com/lm-gautam-bhagat/minio-server/config"
	"github.com/lm-gautam-bhagat/minio-server/log"
	"github.com/lm-gautam-bhagat/minio-server/server"
	"github.com/lm-gautam-bhagat/minio-server/storage"
)

func main() {
	cfg, err := config.LoadConfig("/home/gautam/Desktop/DevHub/GoHub/Lemma/minio-server/resource/config/config.yml")
	if err != nil {
		fmt.Printf("err %v", err.Error())
		return
	}

	endpoint := "127.0.0.1:9000"

	strClient, err := storage.NewStorageClient(endpoint, cfg.MinioRoot.User, cfg.MinioRoot.Pass, false)
	if err != nil {
		fmt.Printf("error while initializing storage client: %+s", err.Error())
		return
	}
	strAdmin, err := storage.NewStorageAdmin(endpoint, cfg.MinioRoot.User, cfg.MinioRoot.Pass, false)
	if err != nil {
		fmt.Printf("error while initializing storage admin: %+s", err.Error())
		return
	}
	s := server.NewServer(cfg, strClient, strAdmin)
	err = s.Setup()
	if err != nil {
		return
	}
	log.Debug("Log Setup")
	s.StartServer()
}
