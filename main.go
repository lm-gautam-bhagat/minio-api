package main

import (
	"fmt"

	"github.com/lm-gautam-bhagat/minio-server/config"
	"github.com/lm-gautam-bhagat/minio-server/log"
	"github.com/lm-gautam-bhagat/minio-server/server"
)

func main() {
	cfg, err := config.LoadConfig("/home/gautam/Desktop/DevHub/GoHub/Lemma/minio-server/resource/config/config.yml")
	if err != nil {
		fmt.Printf("err %v", err.Error())
		return
	}
	fmt.Println(config.GetConfig().MinioRoot.Pass)
	s := server.NewServer(cfg)
	err = s.Setup()
	if err != nil {
		return
	}
	log.Debug("Log Setup complete")
	s.StartServer()
}
