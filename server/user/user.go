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
			name VARCHAR(100),
			uuid  VARCHAR(255),
			email  VARCHAR(100),
			password VARCHAR(255),
			creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			Fb_account_linked BOoLEAN,
			email_verif BOOLEAN
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
	stmt, errs := db.Prepare("INSERT INTO user (uuid, email, password, email_verif) VALUES(?, ?, ?, ?)")
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
	err = db.QueryRow("SELECT uuid, password, email_verif  FROM user WHERE email = ?", u.Email).Scan(&u.Uid, &u.Password, &u.MailVerif)
	if err != nil {
		return err
	}
	return nil
}
