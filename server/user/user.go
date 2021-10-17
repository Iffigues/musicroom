package user

import (
	"github.com/iffigues/musicroom/util"
	"github.com/iffigues/musicroom/postmail"
	"errors"
	"log"
	"fmt"
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
			Fb_account_linked BOOLEAN DEFAULT FALSE,
			email_verif BOOLEAN DEFAULT false
		)`
	if _, err := db.Exec(user); err != nil {
		log.Fatal(err)
	}

	verif := `CREATE TABLE IF NOT EXISTS verif_user (
		id INT primary key auto_increment,
		user_id INT  NOT NULL,
		uuid VARCHAR(255),
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
	)`
	if _, err := db.Exec(verif); err != nil {
		log.Fatal(err)
	}
}

func (a *UserUtils) SendMail(u string) (err error) {
	e := postmail.NewEmail(a.S.Data.Conf)
	e.AddTos("denoyelle.boris@gmail.com")
	fmt.Println(e.Html("./mailtemplate/register", u))
	e.Auths()
	return e.Send()
}

func (a *UserUtils) AddUser(u *User) (err error){
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, errs := db.Prepare("INSERT INTO user (uuid, email, password) VALUES(?, ?, ?)")
	if errs != nil {
		return errs
	}
	id, errt := stmt.Exec(util.Uid().String(), u.Email, u.Password)
	if errt != nil {
		return errt
	}
	lid, errno := id.LastInsertId()
	if errno != nil {
		return errno
	}
	stmt, err = db.Prepare("INSERT INTO verif_user (user_id, uuid) VALUES(?, ?)")
	if err != nil {
		return err
	}
	st := util.Uid().String()
	_, err = stmt.Exec(lid, st)
	if err != nil {
		return err
	}
	return a.SendMail(st)
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

func (a *UserUtils) GetUseriVerif(uid string) (err error){

	db, err := a.S.Data.Bdd.Connect()

	if err != nil {
		return err
	}

	defer db.Close()

	ee := ""

	err = db.QueryRow("SELECT user_id, FROM verif_user WHERE uuid = ?", uid).Scan(&ee)

	if err != nil {
		return err
	}

	if ee == "" {
		return errors.New("empty user")
	}

	_, err = db.Exec("UPDATE user SET email_verif = TRUE WHERE uuid = ?", ee)

	return err
 }
