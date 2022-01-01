package room

import (
	"strings"
	"log"
    "fmt"
	"errors"
	"strconv"
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
			private BOOLEAN DEFAULT false,
			playlist VARCHAR(255),
			current_position INT DEFAULT 0,
			creator_id INT NOT NULL,
			creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_modif TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			playlist_type VARCHAR(255) NOT NULL
		)`
	if _, err := db.Exec(room); err != nil {
		log.Fatal(err)
	}

	invite := `CREATE TABLE IF NOT EXISTS invite (
		id INT primary key auto_increment,
		room_id INT NOT NULL,
		user_id INT NOT NULL,
		date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
		FOREIGN KEY (room_id) REFERENCES room(id) ON DELETE CASCADE
	)`
	if _, err := db.Exec(invite); err != nil {
		log.Fatal(err)
	}

	song := `CREATE TABLE IF NOT EXISTS song(
		id INT primary key auto_increment,
		name  VARCHAR(255),
		author VARCHAR(255),
		ranking INT DEFAULT 0,
		track_id VARCHAR(255),
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

	vote := `CREATE TABLE IF NOT EXISTS vote (
		id INT primary key auto_increment,
		user_id INT NOT NULL,
		song_id INT NOT NULL,
		vote INT DEFAULT 0,
		date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
		FOREIGN KEY (song_id) REFERENCES song(id) ON DELETE CASCADE
	)`
	if _, err := db.Exec(vote); err != nil {
		log.Fatal(err)
	}
}

//NOTE: Room methods
//NOTE: Add a new room to BDD
func (a *RoomUtils) AddRoom(r *Room) (err error){
	r.Playlist = ""
	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, errs := db.Prepare("INSERT INTO room (name, creator_id, playlist_type, playlist) VALUES(?, ?, ?, ?)")
	if errs != nil {
		return errs
	}
	_, err = stmt.Exec(r.Name, r.CreatorId, r.PlaylistType, r.Playlist)
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
	err = db.QueryRow("SELECT name, creator_id, private, playlist, playlist_type, current_position FROM room WHERE id = ?", r.Id).Scan( &r.Name, &r.CreatorId, &r.Private, &r.Playlist, &r.PlaylistType, &r.Position)
	if err != nil {
		return err
	}

	//NOTE : Get playlist Songs
	rows, err := db.Query("SELECT id, name, author, ranking, track_id, room_id, isplayed FROM song WHERE room_id = ?", r.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.Id, &song.Name, &song.Author, &song.Ranking, &song.TrackId, &song.RoomId, &song.IsPlayed); err != nil {
			return err
		}
		r.Song = append(r.Song, song)
	}
	/*if err = rows.Err(); err != nil {
		return err
	}*/
	return nil
}




//NOTE: Song methods

/*id INT primary key auto_increment,
name  VARCHAR(255),
author VARCHAR(255),
ranking INT DEFAULT 0,
track_id VARCHAR(255),
room_id INT,
creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
last_modif TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
isplayed BOOLEAN DEFAULT false,
CONSTRAINT musicroom_room*/

//NOTE: Add a new song to BDD
func (a *RoomUtils) AddSong(s *Song) (err error){
	var p string
	var g bool

	db, err := a.S.Data.Bdd.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	eo := `SELECT IF(COUNT(*),'true','false') FROM song WHERE track_id = ? AND room_id = ?`
	errs := db.QueryRow(eo, s.TrackId, s.RoomId).Scan(&g)
	if errs != nil {
		return errs
	}
	if g {
		return (errors.New("Song already in playlist"))
	}
	eo = `INSERT INTO song (name, room_id, author, track_id) VALUES(?, ?, ?, ?)`
	_, errs = db.Exec(eo, s.Name, s.RoomId, s.Author, s.TrackId)
	if errs != nil {
		return errs
	}

	eo = `SELECT s.id, r.playlist FROM song as s
	JOIN room as r ON s.room_id = r.id
	WHERE s.room_id = ? AND track_id = ?
	`
	errs = db.QueryRow(eo, s.RoomId, s.TrackId).Scan(&s.Id, &p)
	if errs != nil {
		return errs
	}
	if p == "" {
		p = strconv.Itoa(s.Id)
	} else {
		p = p + "," + strconv.Itoa(s.Id)
	}

	eo = `UPDATE room SET playlist = ? WHERE id = ?`
	_, errs = db.Exec(eo, p, s.RoomId)
	if errs != nil {
		return errs
	}
	return
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