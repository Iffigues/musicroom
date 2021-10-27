package user

import (
	"github.com/gin-gonic/gin"
	"fmt"
)

type Adder struct {
	F string `json:"friend"`
}

func (u *UserUtils) AddFriend(c *gin.Context) {
	e, ee :=ExtractTokenMetadata(c.Request)
	fmt.Println(e.UserId, ee)
	var g bool
	var f Adder
	c.BindJSON(&f)
	db, err := u.S.Data.Bdd.Connect()
	if err != nil {
		return 
	}
	eo := `SELECT IF(COUNT(*),'true','false') FROM friends WHERE user_id = (SELECT id FROM user WHERE uuid = ?) AND friend_id = (SELECT id FROM user WHERE uuid = ?)`
	errs := db.QueryRow(eo, e.UserId, f.F).Scan(&g)
	if errs != nil {
		fmt.Println(err)
	}
	if !g {
		t := `INSERT INTO friends (user_id, friend_id) VALUES((SELECT id FROM user WHERE uuid =?),(SELECT id FROM user WHERE uuid = ?))`
		if _, err := db.Exec(t, e.UserId, f.F); err  != nil {
			fmt.Println(err)
		} else {
			fmt.Println("oui")
		}
	}
}

