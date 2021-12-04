package room

import (

	"github.com/iffigues/musicroom/server"
	"github.com/iffigues/musicroom/user"
	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"

)

/*Sommaire NOTE: 
	struct : Struct server, room, song
	NewRoom : Create table in bdd
	Room  : Api handler Add, Get, Upd, Del song
	Song : 	Room  : Api handler Add, Get, Upd, Del song
	ROUTE : api routing
*/
//NOTE: Test Locale A retirer en production

type RoomUtils struct {
	S *server.Server
}

type Room struct {
	Id int `json:"id"`
	Name	string `json:"name"`
	CreatorId	int `json:"creator_id"`
	Private bool `json:"private"`
	Song []Song `json:"song"`
	Playlist string `json:"playlist"`
	Position string `json:"current_position"`
	
}

type Song struct {
	Id int `json:"id"`
	RoomId int `json:"roomid"`
	Name	string `json:"name"`
	Author	string `json:"author"`
	Ranking int `json:"ranking"`
	IsPlayed	bool   `json:"isplayed"`
}

type MoveSong struct {
	RoomId int `json:"room_id"`
	SongId string `json:"song_id"`
	Position int `json:"position"`
}

//NOTE: Create table in bdd
func NewRoom(s *server.Server) (r *RoomUtils) {
	r = new(RoomUtils)
	r.S = s
	r.InitRoom()
	return
}

//NOTE: Room Api handler
//NOTE: Api handler add a room
func (r *RoomUtils)RoomHandler(c *gin.Context) {
	var room Room
	c.BindJSON(&room)
	if err := r.AddRoom(&room); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	c.JSON(200, gin.H{"status": "OK"})
}

//NOTE: Api handler get all rooms
func (r *RoomUtils)GetRoomsHandler(c *gin.Context) {
	var rooms []Room
	/*e, err := user.ExtractTokenMetadata(c.Request)
	if err != nil {
		return
	}*/
	/*
	if err := r.GetAllPublicRoom(&rooms); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	if err := r.GetAllMyRoom(&rooms); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	if err := r.GetAllPublicRoom(&rooms); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}*/
	db, err := r.S.Data.Bdd.Connect()
	if err != nil {
		c.JSON(400, gin.H{"status": "badbdd"})
		return
	}
	defer db.Close()
	//NOTE : Select all public room AND all user's rooms
	rows, err := db.Query("SELECT name, creator_id From room WHERE private = false OR creator_id = (SELECT id FROM user WHERE uuid = ?)", "testkey") //e.UserId

	if err != nil {
		c.JSON(400, gin.H{"status": "badrequest"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.Name, &room.CreatorId); err != nil {
			c.JSON(400, gin.H{"status": "badscan"})
			return
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		c.JSON(400, gin.H{"status": "badrow"})
		return
	}
    c.IndentedJSON(http.StatusOK, rooms)
}



//NOTE: Api handler get room by id
func (r *RoomUtils)GetRoomHandler(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	var room Room
	room.Id = id
	if err := r.GetRoom(&room); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
    c.IndentedJSON(http.StatusOK, room)
}



//NOTE: Api handler del room by id
func (r *RoomUtils)DelRoomHandler(c *gin.Context) {
	//NOTE: Test Locale-comment this part
	e, ee := user.ExtractTokenMetadata(c.Request)
	if ee != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	var g bool
	var room Room
	c.BindJSON(&room)
	db, errd := r.S.Data.Bdd.Connect()
	if errd != nil {
		return
	}
	defer db.Close()
	eo := `SELECT IF(COUNT(*),'true','false') FROM room WHERE id = ? AND creator_id = (SELECT id FROM user WHERE uuid = ?)  LIMIT 1`
	errs := db.QueryRow(eo, room.Id, e.UserId).Scan(&g)
	//NOTE: Test Locale
	//errs := db.QueryRow(eo, room.Id, "testkey").Scan(&g)
	if errs != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	if !g {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	eo = `DELETE FROM room WHERE id = ? AND creator_id = (SELECT id FROM user WHERE uuid = ?)`
	_, errs = db.Exec(eo, room.Id, e.UserId)
	//NOTE: Test Locale
	//_, errs = db.Exec(eo, room.Id, "testkey")
	if errs != nil {
		c.JSON(400, gin.H{"status": "bad"})
	}
	c.JSON(http.StatusOK, "Successfully remove room")
}

//NOTE: Song Api handler

func TokenAuthMiddleware(c *gin.Context) {
	err := user.TokenValid(c.Request)
	if err != nil {
	   c.JSON(http.StatusUnauthorized, err.Error())
	   c.Abort()
	   return
	}
	c.Next()
}

//NOTE: Song handler add a song
func (r *RoomUtils)SongHandler(c *gin.Context) {
	var song Song
	c.BindJSON(&song)
	if err := r.AddSong(&song); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	c.JSON(200, gin.H{"status": "OK"})
}

func (r *RoomUtils)MoveSongHandler(c *gin.Context) {
	var song MoveSong
	c.BindJSON(&song)
	if err := r.MoveSong(&song); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	c.JSON(200, gin.H{"status": "OK"})
}


//NOTE: ROUTE
func (r *RoomUtils) WWW(s *server.Server) {
	s.NewR("/rooms", "rooms", "GET", 1, r.S.MakeMe(r.GetRoomsHandler))
	s.NewR("/rooms/:id", "roombyid", "GET", 1, r.S.MakeMe(r.GetRoomHandler))
	s.NewR("/rooms/add", "room", "POST", 1, r.S.MakeMe(r.RoomHandler))
	s.NewR("/rooms/delete", "delroombyid", "POST", 1, r.S.MakeMe(r.DelRoomHandler))
	//s.NewR("/rooms/edit", "editroombyid", "POST", 1, r.S.MakeMe(r.EditRoomHandler)) 

	//s.NewR("/rooms/song", "song", "GET", 1, r.S.MakeMe(r.GetSongsHandler))
	//s.NewR("/rooms/song/:id", "songbyid", "GET", 1, r.S.MakeMe(r.GetSongHandler))
	s.NewR("/rooms/song/add", "song", "POST", 1, r.S.MakeMe(r.SongHandler))
	s.NewR("/rooms/song/move", "movesong", "POST", 1, r.S.MakeMe(r.MoveSongHandler))
	//s.NewR("/rooms/song/remove", "removesong", "POST", 1, r.S.MakeMe(r.RemoveSongHandler))

}
