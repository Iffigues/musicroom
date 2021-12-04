package room

import (
	"strings"
	"log"
    "fmt"
	"errors"
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
			private BOOLEAN DEFAULT true,
			playlist VARCHAR(255),
			current_position INT,
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

//NOTE: get all public rooms
func (a *RoomUtils) GetAllPublicRoom(r *[]Room) (err error){
	
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query("SELECT name, creator_id From room WHERE private = false")
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

//NOTE: get all user's rooms
func (a *RoomUtils) GetAllMyRoom(r *[]Room) (err error){
	
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, creator_id From room WHERE creator_id = false")
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

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}


func (a *RoomUtils) MoveSong(s *MoveSong) (err error){

	db, errd := a.S.Data.Bdd.Connect()
	if errd != nil {
		return errd
	}
	var pstring string
	var currentPos int
	var index int
	var g bool = false
	var p []string

	defer db.Close()
	eo := `SELECT current_position, playlist FROM room WHERE id = ?`
	errs := db.QueryRow(eo, s.RoomId).Scan(&currentPos, &pstring)
	if errs != nil {
		fmt.Println(errs)
		return errs
	}

	//Transforme la playlist en array
	oldPlaylist := strings.Split(pstring, ",")

	err1 := errors.New("error: invalid input")
	//Check if position input is valid
	if s.Position >= len(oldPlaylist) {
		return err1
	}
	//Sauvegarder le song selectionné
	curseurId := oldPlaylist[currentPos]

	//Aller chercher l'index du song à déplacer
	for i, element := range oldPlaylist {
		if element == s.SongId {
			index = i
			g = true
		}
	}
	if !g {
		return err1
	}
	//Supprimer le song à déplacer
	oldPlaylist = remove(oldPlaylist, index)
	//Boucler sur 
	for i, element := range oldPlaylist {

		if (i == s.Position){
			p = append(p, s.SongId)
		}
		p = append(p, element)

	}
	if p[currentPos] != curseurId {
		for i, element := range p {
			if element == curseurId {
				currentPos = i
			}
		}
	}

	newString := strings.Join(p, ",")

	eo = `UPDATE room SET current_position = ?, playlist = ? WHERE id = ?`
	_, err = db.Exec(eo, currentPos, newString, s.RoomId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return
}


func (a *RoomUtils) DeleteSong(s *Song) (err error){
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, errs := db.Prepare("DELETE FROM song WHERE id = ?")
	if errs != nil {
		return errs
	}
	_, err = stmt.Exec(s.Id)
	if err != nil {
		return err
	}
	return nil
}