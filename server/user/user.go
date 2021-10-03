package user

import (
	"github.com/iffigues/musicroom/util"

	"time"
	"log"
)

func (a *UserUtils) InitUser() {
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	user := `CREATE TABLE IF NOT EXISTS user(
			id INT primary key auto_increment,
			uuid  VARCHAR(255),
			authuuid VARCHAR(255),
			email  VARCHAR(100),
			password VARCHAR(100) NOT NULL,
			types BOOLEAN DEFAULT false,
			buy BOOLEAN DEFAULT false,
			emailverif BOOLEAN DEFAULT false,
			tokenVerif varchar(255),
			tokenexp DATE DEFAULT NULL,
			emailtype VARCHAR(255) UNIQUE NOT NULL
		)`
	if _, err := db.Exec(user); err != nil {
		log.Fatal(err)
	}
}

func (a *UserUtils) AddUser(u *User) (err error){
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, errs := db.Prepare("INSERT INTO user (uuid, email, password, tokenverif, tokenexp, emailtype) VALUES(?, ?, ?, ?, ?, ?)")
	if errs != nil {
		return errs
	}
	_, err = stmt.Exec(util.Uid().String(), u.Email, u.Password, util.Uid().String(), time.Now().AddDate(0, 1, 0), "0" + u.Email)
	if err != nil {
		return err
	}
	return nil
}

func (a *UserUtils) GetUser(u *User) (err error){
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.QueryRow("SELECT uuid, password, emailverif  FROM user WHERE email = ?", u.Email).Scan(&u.Uid, &u.Password, &u.MailVerif)
	if err != nil {
		return err
	}
	return nil
}
