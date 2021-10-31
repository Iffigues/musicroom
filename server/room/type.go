package room

import (

	"github.com/iffigues/musicroom/server"
	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"

)

type RoomUtils struct {
	S *server.Server
}

type Room struct {
	Id int `json:"id"`
	Name	string `json:"name"`
	CreatorId	int `json:"creator_id"`
}

// albums slice to seed record album data.
var Albums = []Room{
    {Id: 1, Name: "Blue Train", CreatorId: 1},
    {Id: 2, Name: "Totoro", CreatorId: 1},
    {Id: 3, Name: "Titi", CreatorId: 1},
}


func NewRoom(s *server.Server) (r *RoomUtils) {
	r = new(RoomUtils)
	r.S = s
	r.InitRoom()
	return
}



func (r *RoomUtils)RoomHandler(c *gin.Context) {
	var room Room
	c.BindJSON(&room)
	if err := r.AddRoom(&room); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	c.JSON(200, gin.H{"status": "OK"})
}

func (r *RoomUtils)GetAllRooms(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, Albums)

}

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

func (r *RoomUtils) WWW(s *server.Server) {
	s.NewR("/room/add", "room", "POST", 1, r.S.MakeMe(r.RoomHandler))
	s.NewR("/rooms", "rooms", "GET", 1, r.S.MakeMe(r.GetAllRooms))
	s.NewR("/rooms/:id", "roombyid", "GET", 1, r.S.MakeMe(r.GetRooms))
}
