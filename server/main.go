package main

import (
	"github.com/iffigues/musicroom/init"
	"github.com/iffigues/musicroom/config"
	"github.com/iffigues/musicroom/logger"
	"github.com/iffigues/musicroom/server"
	"log"
)

func main() {

	logs := logger.NewLog("./log/music-room.log");
	ini, err := inits.NewInit("./conf/ini.ini")

	if err != nil {
		log.Fatal(err);
	}

	conf := config.NewConf()
	conf.NewConfType("http", true)
	err = conf.AddState("http", "socket", ini.GetKey("http", "Socket"), true)

	if err != nil {
		log.Fatal(err);
	}

	server := server.NewServer(conf)
	serve := server.Servers();
	err = serve.ListenAndServe()
	if err != nil {
		logs.Info.Println(err)
	}
}

