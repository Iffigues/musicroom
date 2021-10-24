package room

import (
	"log"

)



func (a *RoomUtils) InitRoom() {
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	room := `CREATE TABLE IF NOT EXISTS room(
			id INT primary key auto_increment,
			name  VARCHAR(255),
			creator_id INT NOT NULL,
			creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_modif TIMESTAMP ON UPDATE CURRENT_TIMESTAMP

		)`
	if _, err := db.Exec(room); err != nil {
		log.Fatal(err)
	}

	song := `CREATE TABLE IF NOT EXISTS song(
		id INT primary key auto_increment,
		name  VARCHAR(255),
		author VARCHAR(255),
		ranking INT,
		creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_modif TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		isplayed BOOLEAN DEFAULT false
	)`
	if _, err := db.Exec(song); err != nil {
		log.Fatal(err)
	}
}

func (a *RoomUtils) AddRoom(r *Room) (err error){
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, errs := db.Prepare("INSERT INTO room (name, creator_id) VALUES(?, ?)")
	if errs != nil {
		return errs
	}
	_, err = stmt.Exec(r.Name, r.CreatorId)
	if err != nil {
		return err
	}
	return nil
}


func (a *RoomUtils) GetRoom(r *Room) (err error){
	
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.QueryRow("SELECT name, creator_id FROM room WHERE id = ?", r.Id).Scan(&r.Name, &r.CreatorId)
	if err != nil {
		return err
	}
	return nil
}


