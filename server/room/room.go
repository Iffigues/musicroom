package room

import (
	"log"

)

/*Sommaire NOTE: 
	InitRoom : Init DB room and song
	Room methods : Add, Get, Upd, Del room methods
	Song methods : Add, Get, Upd, Del song methods
*/

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
		track_id INT,
		room_id INT,
		creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_modif TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		isplayed BOOLEAN DEFAULT false,
		CONSTRAINT musicroom_room
		FOREIGN KEY (room_id) 
			REFERENCES room(id)
			ON DELETE CASCADE
	)`
	if _, err := db.Exec(song); err != nil {
		log.Fatal(err)
	}
}

//NOTE: Room methods
//NOTE: Add a new room to BDD
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

//NOTE: get room by id
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
	//NOTE : Get playlist Songs
	rows, err := db.Query("SELECT id, name FROM song WHERE room_id = ?", r.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.Id, &song.Name); err != nil {
			return err
		}
		r.Song = append(r.Song, song)
	}
	/*if err = rows.Err(); err != nil {
		return err
	}*/
	return nil
}

//NOTE: get all rooms
func (a *RoomUtils) GetAllRoom(r *[]Room) (err error){
	
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query("SELECT name, creator_id FROM room")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.Name, &room.CreatorId); err != nil {
			return err
		}
		*r = append(*r, room)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
	
}


//NOTE: Delete room by id
func (a *RoomUtils) DeleteRoom(rid int) (err error){
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, errs := db.Prepare("DELETE FROM room WHERE id = ?")
	if errs != nil {
		return errs
	}
	_, err = stmt.Exec(rid)
	if err != nil {
		return err
	}
	return nil
}



//NOTE: Song methods

//NOTE: Add a new song to BDD
func (a *RoomUtils) AddSong(s *Song) (err error){
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, errs := db.Prepare("INSERT INTO song (name, room_id) VALUES(?, ?)")
	if errs != nil {
		return errs
	}
	_, err = stmt.Exec(s.Name, s.RoomId)
	if err != nil {
		return err
	}
	return nil
}






