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
	if ee != nil {
		return
	}
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
		t := `INSERT INTO friends (user_id, friend_id) VALUES((SELECT id FROM user WHERE uuid =?),(SELECT id FROM user WHERE uuid = ?))`
		if _, err := db.Exec(t, e.UserId, f.F); err  != nil {
			fmt.Println(err)
		} else {
			fmt.Println("oui")
		}
	}
}

func (a *UserUtils) GetAllFriend(c *gin.Context) {
	fmt.Print("eazeaz")
	e, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return
	}
	db, errs := a.S.Data.Bdd.Connect()
	if errs != nil {
		return
	}
	eo := `SELECT a.id FROM user as q
		JOIN friend as a ON a.user_id = q.id OR a.friend_id = q.id
		WHERE q.uuid = ? 
	`
	rows, errt := db.Query(eo, e.UserId)
	if errt != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var t int
		rows.Scan(&t)
		fmt.Println("t=", t)
	}
}

func (u *UserUtils) ShowFriend(c *gin.Context) {
	var t []int
	e, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return
	}
	eo := `SELECT user.id From user WHERE id = (SELECT user_id FROM friends WHERE friend_id = (SELECT id FROM user WHERE uuid = ?))`
	db, errs := u.S.Data.Bdd.Connect()
	if errs != nil {
		return
	}
	rows, errt := db.Query(eo, e.UserId)
	if errt != nil {
		fmt.Println(errt)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tt int
		if err := rows.Scan(&tt); err != nil {
			fmt.Println(err)
		}
		t = append(t, tt)
	}
	fmt.Println(t)
	if err := rows.Err(); err != nil {
	   fmt.Println(err)
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
