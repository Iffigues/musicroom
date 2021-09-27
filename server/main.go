package main

import (
	"github.com/iffigues/musicroom/config"
	"github.com/iffigues/musicroom/init"
	"github.com/iffigues/musicroom/logger"
	"github.com/iffigues/musicroom/server"
	"github.com/iffigues/musicroom/user"
	"github.com/iffigues/musicroom/util"

	"github.com/sevlyar/go-daemon"
	"log"
)

func serves() {
	util.CreateDir("./log")
	logs := logger.NewLog("./log/music-room.log")
	ini, err := inits.NewInit("./conf/ini.ini")
	if err != nil {
		logs.Fatal(err.Error())
	}
	conf := config.NewConf()
	conf.NewConfType("http", true)
	conf.NewConfType("bdd", true)
	conf.NewConfType("gin", true)
	err = conf.AddState("http", "socket", ini.GetKey("http", "Socket"), true)
	if err != nil {
		log.Fatal(err)
	}
	conf.AddState("bdd", "host", "localhost", true)
	conf.AddState("bdd", "user", "root", true)
	conf.AddState("bdd", "pwd", "Petassia01", true)
	conf.AddState("bdd", "dbname", "musicroom", true)
	conf.AddState("gin", "mode", ini.GetKey("gin-mode", "mode"), true)
	server := server.NewServer(conf)
	user := user.NewUser(server)
	server.AddHH(user)
	serve := server.Servers()
	err = serve.ListenAndServe()
	if err != nil {
		logs.Info.Println(err)
	}
}

func main() {
	cntxt := &daemon.Context{
		PidFileName: "./log/taskmaster.pid",
		PidFilePerm: 0777,
		LogFileName: "./log/sample.log",
		LogFilePerm: 0777,
		WorkDir:     "./",
		Umask:       022,
		Args:        []string{"l"},
	}
	d, err := cntxt.Reborn()
	if err != nil {
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	serves()
}
