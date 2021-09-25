package main

import (
	"github.com/iffigues/musicroom/init"
	"github.com/iffigues/musicroom/config"
	"github.com/iffigues/musicroom/logger"
	"github.com/iffigues/musicroom/server"
	"github.com/iffigues/musicroom/user"
	"log"
)

func main() {
	logs := logger.NewLog("./log/music-room.log");
	ini, err := inits.NewInit("./conf/ini.ini")

	if err != nil {
		logs.Fatal(err.Error());
	}

	conf := config.NewConf()
	conf.NewConfType("http", true)
	conf.NewConfType("bdd", true)

	err = conf.AddState("http", "socket", ini.GetKey("http", "Socket"), true)

	if err != nil {
		log.Fatal(err);
	}

	conf.AddState("bdd","host","localhost", true)
	conf.AddState("bdd","user","root", true)
	conf.AddState("bdd","pwd","Petassia01", true)
	conf.AddState("bdd","dbname","musicroom", true)
	server := server.NewServer(conf)
	user := user.NewUser(server)
	server.AddHH(user)
	serve := server.Servers();
	err = serve.ListenAndServe()
	if err != nil {
		logs.Info.Println(err)
	}
}

