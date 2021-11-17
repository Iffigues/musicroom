package room

import (

	"github.com/iffigues/musicroom/server"
	"github.com/iffigues/musicroom/user"
	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"
	"fmt"

)

/*Sommaire NOTE: 
	struct : Struct server, room, song
	NewRoom : Create table in bdd
	Room  : Api handler Add, Get, Upd, Del song
	Song : 	Room  : Api handler Add, Get, Upd, Del song
	ROUTE : api routing
*/

type RoomUtils struct {
	S *server.Server
}

type Room struct {
	Id int `json:"id"`
	Name	string `json:"name"`
	CreatorId	int `json:"creator_id"`
	Song []Song `json:"song"`
}

type Song struct {
	Id int `json:"id"`
	RoomId int `json:"roomid"`
	Name	string `json:"name"`
	Author	string `json:"author"`
	Ranking int `json:"ranking"`
	IsPlayed	bool   `json:"isplayed"`
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
	if err := r.GetAllRoom(&rooms); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
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
	//NOTE: Verif client Rights
	e, ee := user.ExtractTokenMetadata(c.Request)
	if ee != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	var room Room
	c.BindJSON(&room)
	db, err := r.S.Data.Bdd.Connect()
	if err != nil {
		return
	}
	defer db.Close()
	eo := `SELECT creator_id FROM room WHERE id = ?`
	errs := db.QueryRow(eo, room.Id).Scan(&room.CreatorId)
	if errs != nil {
		fmt.Println(err)
	}
	id, err := strconv.Atoi(e.UserId)
	if err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	if room.CreatorId != id {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//NOTE: Run Delete Action
	if err := r.DeleteRoom(room.Id); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
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

//NOTE: ROUTE
func (r *RoomUtils) WWW(s *server.Server) {
	s.NewR("/room/add", "room", "POST", 1, r.S.MakeMe(r.RoomHandler))
	s.NewR("/rooms", "rooms", "GET", 1, r.S.MakeMe(r.GetRoomsHandler))
	s.NewR("/rooms/:id", "roombyid", "GET", 1, r.S.MakeMe(r.GetRoomHandler))
	s.NewR("/rooms/delete", "delroombyid", "POST", 1, r.S.MakeMe(r.DelRoomHandler))

	s.NewR("/song/add", "song", "POST", 1, r.S.MakeMe(r.SongHandler))

}
