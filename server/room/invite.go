package room

import (
	"github.com/gin-gonic/gin"
	"github.com/iffigues/musicroom/user"


	"net/http"
	"strconv"
	"fmt"
)

type Invite struct {
	Room int `json:"room" binding:"required"`
	Invite int `json:"invite" binding:"required"`
}

func (r *RoomUtils) AddInvite(c *gin.Context) {
	e, ee := user.ExtractTokenMetadata(c.Request)
	if ee != nil {
		return
	}
	var g bool
	var i Invite

	c.BindJSON(&i)
	if i.invite == e.UserId {
		c.JSON(400, gin.H{"status": "doublon"})
		return
	}

	db, err := r.S.Data.Bdd.Connect()
	if err != nil {
		return
	}
	defer db.Close()
	eo := `SELECT IF(COUNT(*),'true','false') FROM invite WHERE room_id = ? AND user_id = ?`
	errs := db.QueryRow(eo, i.Room, i.Invite).Scan(&g)
	if errs != nil {
		fmt.Println(err)
		//err.Errors()
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	if !g {
		t := `INSERT INTO invite (room_id, user_id) VALUES((SELECT id FROM room WHERE id =?),(SELECT id FROM user WHERE id = ?))`
		if _, err := db.Exec(t, i.Room, i.Invite); err  != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"status": "bad"})
			return
		} else {
			fmt.Println("guy is invited")
			return
		}
	} else {
		c.JSON(400, gin.H{"status": "doublon"})
		return
	}
}

func (r *RoomUtils) ShowInvite(c *gin.Context) {
	var t []int
	room, err := strconv.Atoi(c.Param("room"))
	if err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	e, err := user.ExtractTokenMetadata(c.Request)
	if err != nil {
		return
	}
	eo := `SELECT user_id From invite WHERE room_id = ?`
	db, errs := r.S.Data.Bdd.Connect()
	if errs != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	defer db.Close()
	rows, errt := db.Query(eo, room)
	if errt != nil {
		fmt.Println(errt)
		c.JSON(400, gin.H{"status": "badrow"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tt int
		if err := rows.Scan(&tt); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"status": "badrow"})
			return
		}
		t = append(t, tt)
	}
	//fmt.Println(t)
	if err := rows.Err(); err != nil {
	   fmt.Println(err)
	   c.JSON(400, gin.H{"status": "bad"})
	   return
	}
	c.IndentedJSON(http.StatusOK, t)

}

func (r *RoomUtils) DelInvite(c *gin.Context) {
	e, err := user.ExtractTokenMetadata(c.Request)
	if err != nil {
		return
	}
	var i Invite
	c.BindJSON(&i)

	db, errd := r.S.Data.Bdd.Connect()
	if errd != nil {
		return
	}
	defer db.Close()
	if !isOwner(e.UserId, i.Room) {
	//if !isOwner(r, "testkey", i.Room) {
		c.JSON(400, gin.H{"status": "unauthorized"})
		return
	}

	eo := `DELETE FROM invite WHERE room_id = ? AND  user_id = ?`
	db.Exec(eo, i.Room, i.Invite)
	c.JSON(200, gin.H{"status": "OK"})
}

func isOwner(r *RoomUtils, uid string, rid int) (g bool) {
	db, errd := r.S.Data.Bdd.Connect()
	if errd != nil {
		return
	}
	defer db.Close()
	eo := `SELECT IF(COUNT(*),'true','false') FROM room WHERE id = ? AND creator_id = (SELECT id FROM user WHERE uuid = ?)  LIMIT 1`
	errs := db.QueryRow(eo, rid, uid).Scan(&g)
	if errs != nil {
		return false
	}
	return g
}
