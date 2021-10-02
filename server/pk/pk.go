package pk

import (
	"database/sql"
	"log"
	"github.com/iffigues/musicroom/config"
	_ "github.com/go-sql-driver/mysql"
)

type Pk struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

func NewPk(conf config.Conf) (a *Pk) {
	a = &Pk{}
	if Host := conf.GetValue("bdd", "host"); Host != nil {
		a.Host = Host.(string)
	} else {
	}

	if Port := conf.GetValue("bdd", "port"); Port != nil {
		a.Port = Port.(string)
	} else {
	}

	if  User := conf.GetValue("bdd", "user"); User != nil {
		a.User = User.(string)
	} else {
	}

	if Password := conf.GetValue("bdd", "pwd"); Password != nil {
		a.Password = Password.(string)
	} else {
	}

	if Dbname := conf.GetValue("bdd", "dbname"); Dbname != nil {
		a.Dbname = Dbname.(string)
	} else {
	}
	a.Starter()
	return

}

func (a *Pk) Init() {
	 db, err := sql.Open("mysql", a.User+":"+a.Password+"@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Fatal(err)
	}
	if err  := db.Ping(); err != nil {
		 log.Fatal(err)
	 }
	 _, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + a.Dbname)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *Pk) Reset() {
}

func (a *Pk) Tables(db *sql.DB) {
	haversin := `CREATE FUNCTION haversine(
			        lat1 FLOAT, lon1 FLOAT,
			        lat2 FLOAT, lon2 FLOAT
			     ) RETURNS FLOAT
			    NO SQL DETERMINISTIC
			BEGIN
			    RETURN DEGREES(ACOS(
					COS(RADIANS(lat1)) *
				                  COS(RADIANS(lat2)) *
				                  COS(RADIANS(lon2) - RADIANS(lon1)) +
				                  SIN(RADIANS(lat1)) * SIN(RADIANS(lat2))
				                ));
			END`
	if _, err := db.Exec(haversin); err != nil {
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
	a.Init();
	db, err := a.Connect()
	if err != nil {
		log.Fatal(err)
	}
	a.Tables(db)
	defer db.Close()
}

func (a *Pk) Connect() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", a.User+":"+a.Password+"@/"+a.Dbname+"?charset=utf8mb4")
	if err != nil {
		return nil, err
	}
	if err  = db.Ping(); err != nil {
		return nil, err
	}
	return

}
