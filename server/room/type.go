package room

import (

	"github.com/iffigues/musicroom/server"
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

type RoomUtils struct {
	S *server.Server
}

type Room struct {
	Id int `json:"id"`
	Name	string `json:"name"`
	CreatorId	int `json:"creator_id"`
}

type Song struct {
	Id int `json:"id"`
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
func (r *RoomUtils)GetAllRooms(c *gin.Context) {
	var rooms []Room
	if err := r.GetAllRoom(&rooms); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
    c.IndentedJSON(http.StatusOK, rooms)
}

//NOTE: Api handler get room by id
func (r *RoomUtils)GetRooms(c *gin.Context) {
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


//NOTE: Song Api handler



//NOTE: ROUTE
func (r *RoomUtils) WWW(s *server.Server) {
	s.NewR("/room/add", "room", "POST", 1, r.S.MakeMe(r.RoomHandler))
	s.NewR("/rooms", "rooms", "GET", 1, r.S.MakeMe(r.GetAllRooms))
	s.NewR("/rooms/:id", "roombyid", "GET", 1, r.S.MakeMe(r.GetRooms))
}
