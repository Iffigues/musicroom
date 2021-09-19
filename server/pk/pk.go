package pk

import (
	"database/sql"
	"log"
	"github.com/iffigues/musicroom/config"
	_ "github.com/lib/pq"
)

type Pk struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

func NewPk(conf config.Conf) (a *Pk) {
	a = &Pk{
		Host:     conf.GetValue("bdd", "host").(string),
		Port:     conf.GetValue("bdd", "port").(string),
		User:     conf.GetValue("bdd", "user").(string),
		Password: conf.GetValue("bdd", "pwd").(string),
		Dbname:   conf.GetValue("bdd", "dbname").(string),
	}
	a.Starter()
	return

}

func (a *Pk) Init(db *sql.DB) {
	if err  := db.Ping(); err != nil {
		 log.Fatal(err)
	 }
	 _, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + a.Dbname)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *Pk) IsUsers() (ok bool) {
	db, err := a.Connect()
	if err != nil {
		return
	}
	defer db.Close()
	return
}

func (a *Pk) Starter() {

	db, err := a.Connect()
	if err != nil {
		log.Fatal(err)
	}
	a.Init(db)
	defer db.Close()
}

func (a *Pk) Connect() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", a.User+a.Password+"@/"+a.Dbname+"?charset=utf8mb4")
	if err != nil {
		return nil, err
	}
	if err  = db.Ping(); err != nil {
		return nil, err
	}
	return

}
