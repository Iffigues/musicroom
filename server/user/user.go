package user

import (
	"github.com/iffigues/musicroom/util"
	"github.com/iffigues/musicroom/postmail"
	"errors"
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
			name VARCHAR(100) NOT NULL,
			uuid  VARCHAR(255),
			email  VARCHAR(100) UNIQUE,
			password VARCHAR(255),
			creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			Fb_account_linked BOOLEAN DEFAULT FALSE,
			oauth BOOLEAN Default false,
			email_verif BOOLEAN DEFAULT false
		)`
	if _, err := db.Exec(user); err != nil {
		log.Fatal(err)
	}

	verif := `CREATE TABLE IF NOT EXISTS verif_user (
		id INT primary key auto_increment,
		user_id INT  NOT NULL,
		uuid VARCHAR(255),
		creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
	)`

	forgot := `CREATE TABLE IF NOT EXISTS forgot (
		id INT primary key auto_increment,
		user_id INT NOT NULL UNIQUE,
		uuid VARCHAR(255),
		creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
	)`

	if _, err := db.Exec(forgot); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(verif); err != nil {
		log.Fatal(err)
	}
	friend := `CREATE TABLE IF NOT EXISTS friends (
		id INT primary key auto_increment,
		user_id INT NOT NULL,
		friend_id INT NOT NULL,
		accept INT DEFAULT 0,
		date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
		FOREIGN KEY (friend_id) REFERENCES user(id) ON DELETE CASCADE
	)`
	friends := `CREATE TABLE IF NOT EXISTS friend (
		id INT primary key auto_increment,
		user_id INT NOT NULL,
		friend_id INT NOT NULL,
		ban boolean default false,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
		FOREIGN KEY (friend_id) REFERENCES user(id) ON DELETE CASCADE
	)`

	if _, err := db.Exec(friend);  err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(friends); err != nil {
		log.Fatal(err)
	}
	event := `CREATE EVENT IF NOT EXISTS delete_event
		ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 1 DAY
		ON COMPLETION PRESERVE
	
		DO BEGIN
			DELETE FROM user WHERE creation_date < DATE_SUB(NOW(), INTERVAL 7 DAY) AND email_verif = false;
			DELETE FROM friends WHERE date < DATE_SUB(NOW(), INTERVAL 40 DAY) AND accept = false;
			DELETE FROM verif_user WHERE date < DATE_SUB(NOW(), INTERVAL 40 DAY) AND accept = false;
			DELETE FROM forgot  WHERE date < DATE_SUB(NOW(), INTERVAL 40 DAY) AND accept = false;
		END
	`
	if _, err := db.Exec(event); err != nil {
		log.Fatal(err)
	}
}

func (a *UserUtils) SendMail(u, email string) (err error) {
	e := postmail.NewEmail(a.S.Data.Conf)
	e.AddTos(email)
	e.Html("./mailtemplate/register", "http://gopiko.fr:9000/user/verif/" + u)
	e.Auths()
	return e.Send()
}

func (a *UserUtils) AddUser(u *User) (err error){
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, errs := db.Prepare("INSERT INTO user (uuid, name, email, password) VALUES(?, ?, ?, ?)")
	if errs != nil {
		return errs
	}
	id, errt := stmt.Exec(util.Uid().String(), u.Name, u.Email, u.Password)
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
	return a.SendMail(st, u.Email)
}

func (a *UserUtils) GetUser(u *User) (err error){
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.QueryRow("SELECT uuid, password, email_verif  FROM user WHERE email = ? AND oauth = false", u.Email).Scan(&u.Uid, &u.Password, &u.MailVerif)
	if err != nil {
		return err
	}
	if !u.MailVerif {
		return errors.New("ee")
	}
	return nil
}

func (a *UserUtils) GetUseriVerif(uid string) (err error){

	println(uid)

	db, err := a.S.Data.Bdd.Connect()

	if err != nil {
		return err
	}
	defer db.Close()

	ex := `UPDATE user SET email_verif = true WHERE user.id = (select user_id FROM verif_user WHERE uuid = ?)`
	if row, err := db.Exec(ex,uid); err != nil {
		return err
	} else {
		if a, err := row.RowsAffected(); err != nil {
			return err
		} else {
			if a == 0 {
				return errors.New("No row affected");
			}
		}
	}
	_,err = db.Exec("DELETE FROM verif_user WHERE uuid = ?", uid)
	return err
}

