package room

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"
	"fmt"
)

type Vote struct {
	Song int `json:"song"`
	Vote int `json:"vote"`
}

func (r *RoomUtils) AddVote(c *gin.Context) {
	e, ee :=ExtractTokenMetadata(c.Request)
	if ee != nil {
		return
	}
	var g bool
	var v Vote

	c.BindJSON(&v)

	db, err := r.S.Data.Bdd.Connect()
	if err != nil {
		return
	}
	defer db.Close()
	eo := `SELECT IF(COUNT(*),'true','false') FROM vote WHERE song_id = ? AND user_id = (SELECT id FROM user WHERE uuid = ?)`
	errs := db.QueryRow(eo, v.Song, e.UserId).Scan(&g)
	//errs := db.QueryRow(eo, v.Song, "doublekey").Scan(&g) 

	if errs != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"status": "bad"})
	}
	if !g {
		t := `INSERT INTO vote (song_id, user_id, vote) VALUES((SELECT id FROM song WHERE id =?), (SELECT id FROM user WHERE uuid = ?), ?)`
		//if _, err := db.Exec(t, v.Song, "doublekey", v.Vote); err  != nil {
		if _, err := db.Exec(t, v.Song, e.UserId, v.Vote); err  != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"status": "bad"})
		} else {
			c.JSON(http.StatusOK, "Successfully add vote")
		}
	} else {
		eo = `UPDATE room SET current_position = ?, playlist = ? WHERE id = ?`

		t := `UPDATE vote SET vote = ? WHERE song_id = ? AND user_id = (SELECT id FROM user WHERE uuid = ?)`
		//if _, err := db.Exec(t, v.Vote, v.Song, "doublekey"); err  != nil {
		if _, err := db.Exec(t, v.Vote, v.Song, e.UserId); err  != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"status": "bad"})
		} else {
			c.JSON(http.StatusOK, "Successfully update vote")
		}
	}

}

func (r *RoomUtils) GetRanking(c *gin.Context) {
	e, ee :=ExtractTokenMetadata(c.Request)
	if ee != nil {
		return
	}
	var t int
	song, err := strconv.Atoi(c.Param("song"))
	if err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	eo := `SELECT vote From vote WHERE song_id = ?`
	db, errs := r.S.Data.Bdd.Connect()
	if errs != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	defer db.Close()
	rows, errt := db.Query(eo, song)
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
		t = t + tt
	}
	eo = `UPDATE song SET ranking = ?  WHERE id = ?`
	db.Exec(eo, t, song)
	//fmt.Println(t)
	if err := rows.Err(); err != nil {
	   fmt.Println(err)
	   c.JSON(400, gin.H{"status": "bad"})
	   return
	}
	c.IndentedJSON(http.StatusOK, t)

	
}
