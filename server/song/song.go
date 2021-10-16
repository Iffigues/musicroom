package song

import (
	"github.com/iffigues/musicroom/util"

	"time"
	"log"
)

func (a *SongUtils) InitSong() {
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	song := `CREATE TABLE IF NOT EXISTS song(
			id INT primary key auto_increment,
			name  VARCHAR(255),
			author VARCHAR(255),
			ranking INT,
			datetime DATETIME,
			lasttimestamp TIMESTAMP,
			isplayed BOOLEAN DEFAULT false,
		)`
	if _, err := db.Exec(song); err != nil {
		log.Fatal(err)
	}
}

func (a *SongUtils) AddSong(u *Song) (err error){
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, errs := db.Prepare("INSERT INTO song (name, author, ranking, datetime, lasttimestamp, isplayed) VALUES(?, ?, ?, ?, ?, ?)")
	if errs != nil {
		return errs
	}
	_, err = stmt.Exec(u.Name, u.Author, u.Ranking, time.Now().AddDate(0, 1, 0), time.Now().AddDate(0, 1, 0) + u.IsPlayed)
	if err != nil {
		return err
	}
	return nil
}


func (a *SongUtils) GetSong(u *Song) (err error){
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.QueryRow("SELECT name, author, ranking, isplayed FROM user WHERE name = ?", u.Name).Scan(&u.Name, &u.Author, &u.Ranking, &u.IsPlayed)
	if err != nil {
		return err
	}
	return nil
}
