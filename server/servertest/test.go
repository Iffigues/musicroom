package servertest

import (
	"github.com/iffigues/musicroom/config"
	"github.com/iffigues/musicroom/init"
	"github.com/iffigues/musicroom/logger"
	"github.com/iffigues/musicroom/server"
	"github.com/iffigues/musicroom/util"
	"github.com/iffigues/musicroom/pk"
	"log"
)


func makeConf(ini *inits.Init) (conf *config.Conf) {
	conf = config.NewConf()
	conf.NewConfType("http", true)
	conf.NewConfType("bdd", true)
	conf.NewConfType("gin", true)
	conf.NewConfType("facebook", true)
	err := conf.AddState("http", "socket", ini.GetKey("http", "Socket"), true)

	if err != nil {
		log.Fatal(err)
	}
	err = conf.AddState("bdd", "host", ini.GetKey("bdd", "host"), true)
	if err != nil {
		log.Fatal(err)
	}
	err = conf.AddState("bdd", "user", ini.GetKey("bdd", "user"), true)
	if err != nil {
		log.Fatal(err)
	}

	err = conf.AddState("bdd", "pwd", ini.GetKey("bdd", "pwd"), true)
	if err != nil {
		log.Fatal(err)
	}

	err = conf.AddState("bdd", "dbname", ini.GetKey("bdd", "dbname-test"), true)
	if err != nil {
		log.Fatal(err)
	}

	conf.AddState("gin", "mode", ini.GetKey("gin-mode", "mode"), true)
	conf.AddState("facebook","id", ini.GetKey("facebook","id"), true)
	return conf
}

func Serves()  (*server.Server) {
	util.CreateDir("./logTest")
	logs := logger.NewLog("./logTest/music-room-test.log")
	ini, err := inits.NewInit("./conf/ini.ini")
	if err != nil {
		logs.Fatal(err.Error())
	}
	conf := makeConf(ini)
	ii := pk.NewPk(*conf)
	server := server.NewServer(conf)
	server.AddPk(ii)
	return server
}

func LanceServe(server *server.Server) {
	serve := server.Servers()
	err := serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

