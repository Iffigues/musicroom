package main

import (
	"github.com/iffigues/musicroom/config"
	"github.com/iffigues/musicroom/init"
	"github.com/iffigues/musicroom/logger"
	"github.com/iffigues/musicroom/server"
	"github.com/iffigues/musicroom/user"
	"github.com/iffigues/musicroom/util"
	"github.com/iffigues/musicroom/pk"
	"github.com/sevlyar/go-daemon"
	"os"
	"log"
)


func makeConf(ini *inits.Init) (conf *config.Conf) {
	conf = config.NewConf()
	conf.NewConfType("http", true)
	conf.NewConfType("bdd", true)
	conf.NewConfType("gin", true)
	conf.NewConfType("facebook", true)
	conf.NewConfType("email", true)
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

	err = conf.AddState("bdd", "dbname", ini.GetKey("bdd", "dbname"), true)
	if err != nil {
		log.Fatal(err)
	}

	conf.AddState("gin", "mode", ini.GetKey("gin-mode", "mode"), true)
	conf.AddState("facebook","id", ini.GetKey("facebook","id"), true)
	conf.AddState("facebook", "id", ini.GetKey("facebook", "key"), true)
	conf.AddState("email","user", ini.GetKey("email","user"), true)
	conf.AddState("email","from", ini.GetKey("email","from"), true)
	conf.AddState("email","smtphost", ini.GetKey("email","smtphost"), true)
	conf.AddState("email","smtpport", ini.GetKey("email","smtpport"), true)
	conf.AddState("email","pwd", ini.GetKey("email","pwd"), true)
	conf.AddState("email","mime", ini.GetKey("email","mime"), true)
	return conf
}

func serves() {
	util.CreateDir("./log")
	logs := logger.NewLog("./log/music-room.log")
	ini, err := inits.NewInit("./conf/ini.ini")
	if err != nil {
		logs.Fatal(err.Error())
	}
	conf := makeConf(ini)
	ii := pk.NewPk(*conf)
	server := server.NewServer(conf)
	server.AddPk(ii)
	user := user.NewUser(server)
	server.AddHH(user)
	w := server.NewWw(server)
	server.AddHH(w)
	serve := server.Servers()
	err = serve.ListenAndServe()
	if err != nil {
		logs.Info.Println(err)
	}
}

func init() {
}

func main() {
	t := false
	if len(os.Args) > 1 {
		if os.Args[1] ==  "reset" {
			ini, _ := inits.NewInit("./conf/ini.ini")
			conf := makeConf(ini)
			pk.NewPk(*conf).Reset()
			return
		} else if os.Args[1] == "daemon" {
			t = true
		} else {
			return
		}
	}
	if t {
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
	}
	serves()
}
