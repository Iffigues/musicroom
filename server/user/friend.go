package user

import (
	"github.com/iffigues/musicroom/util"
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
	if f.F == e.UserId {
		return
	}
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
		t := `INSERT INTO friends (uuid, user_id, friend_id) VALUES(?, (SELECT id FROM user WHERE uuid =?),(SELECT id FROM user WHERE uuid = ?))`
		if _, err := db.Exec(t, e.UserId, f.F); err  != nil {
			fmt.Println(err)
		} else {
			fmt.Println("oui")
		}
	}
}

func (u *UserUtils) AcceptFriend(c *gin.Context) {
	e, err := ExtractTokenMetadata(c.Request)
	var g bool
	if err != nil {
		return
	}
	var f Adder
	c.BindJSON(&f)
	db, errd := u.S.Data.Bdd.Connect()
	if errd != nil {
		return
	}
	eo := `SELECT IF(COUNT(*),'true','false') FROM friends WHERE user_id = (SELECT id FROM user WHERE uuid = ?)   AND friend_id = (SELECT id FROM user WHERE uuid = ?)  LIMIT 1`
	errs := db.QueryRow(eo, f.F, e.UserId).Scan(&g)
	if errs != nil {
		fmt.Println(errs)
		return
	}
	if !g {
		print("eee")
		return
	}
	eo = `INSERT INTO  friend  (user_id, friend_id)  VALUES  ((SELECT id FROM user WHERE uuid = ?), (SELECT id FROM user WHERE uuid = ?))`
	_, errs = db.Exec(eo, f.F, e.UserId)
	if errs != nil {
		fmt.Println(errs)
		return
	}
	eo = `DELETE FROM friends WHERE user_id = (SELECT id FROM user WHERE uuid = ?) OR user_id = (SELECT id FROM user WHERE uuid = ?) AND friend_id = (SELECT id FROM user WHERE uuid = ?) OR friend_id = (SELECT id FROM user WHERE uuid = ?)`
	_, err = db.Exec(eo, e.UserId, f.F, f.F, e.UserId)
	if err != nil {
		fmt.Println(err)
	}
}


func (u *UserUtils) BanFriend(c *gin.Context) {
	e, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return
	}
	var f Adder
	c.BindJSON(&f)
	db, errd := u.S.Data.Bdd.Connect()
	if errd != nil {
		return
	}
	eo := `UPDATE friend SET ban = true  WHERE user_id = (SELECT id FROM user WHERE uuid = ?)  OR user_id = (SELECT id FROM user WHERE uuid = ?) AND  friend_id = (SELECT id FROM user WHERE uuid = ?) OR firend_id = (SELECT id FROM user WHERE uuid = ?)`
	db.Exec(eo, e.UserId, f.F,  f.F, e.UserId)
}

func (u *UserUtils) DelFriend(c *gin.Context) {
	e, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return
	}
	var f Adder
	c.BindJSON(&f)
	db, errd := u.S.Data.Bdd.Connect()
	if errd != nil {
		return
	}
	eo := `DELETE FROM friend WHERE user_id = (SELECT id FROM user WHERE uuid = ?)  OR user_id = (SELECT id FROM user WHERE uuid = ?) AND  friend_id = (SELECT id FROM user WHERE uuid = ?) OR firend_id = (SELECT id FROM user WHERE uuid = ?)`
	db.Exec(eo, e.UserId, f.F,  f.F, e.UserId)
}

func (u *UserUtils) Deban(c *gin.Context) {
	e, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return
	}
	var f Adder
	c.BindJSON(&f)
	db, errs := u.S.Data.Bdd.Connect()
	if errs != nil {
		return
	}
	eo := `UPDATE friend SET ban = true WHERE user_id = (SELECT id FROM user WHERE uuid = ?) OR user_id = (SELECT id FROM user WHERE uuid = ?) AND friend_id = (SELECT id FROM user WHERE uuid = ?) OR friend_id = (SELECT id FROM user WHERE uuid = ?)`
	db.Exec(eo, e.UserId, f.F, f.F, e.UserId)
}

func (u *UserUtils) RefuseFriend(c *gin.Context) {
	e, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return
	}
	var f Adder
	c.BindJSON(&f)
	db, errd := u.S.Data.Bdd.Connect()
	if errd != nil {
		return
	}
	eo := `DELETE FROM friends WHERE user_id = (SELECT id FROM user WHERE uuid = ?)  OR user_id = (SELECT id FROM user WHERE uuid = ?) AND  friend_id = (SELECT id FROM user WHERE uuid = ?) OR firend_id = (SELECT id FROM user WHERE uuid = ?)`
	db.Exec(eo, e.UserId, f.F,  f.F, e.UserId)
}
