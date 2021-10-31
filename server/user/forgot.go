package user

import (
	"github.com/iffigues/musicroom/util"
	"github.com/gin-gonic/gin"
	"github.com/iffigues/musicroom/postmail"
)

type Ml struct {
	Email string	`json:"email"`
	Uuid  string	`json:"pwd1"`
}


func (a *UserUtils) SendMailForgot(u, email string) (err error) {
	e := postmail.NewEmail(a.S.Data.Conf)
	e.AddTos(email)
	e.Html("./mailtemplate/forgot", "http://gopiko.fr:9000/user/forgot/" + u)
	e.Auths()
	return e.Send()
}

func (u *UserUtils) Forgot(c *gin.Context) {
	var f  Ml
	c.BindJSON(&f)
	db, errs := u.S.Data.Bdd.Connect()
	if errs != nil {
		return
	}
	g := `INSERT INTO forgot (user_id, uuid) VALUES ((SELECT id FROM user WHERE email = ?), ?)`
	db.Exec(g, f.Email, util.Uid().String())
}

func (u *UserUtils) Change_Forgot(c *gin.Context) {
	var f  Ml
	c.BindJSON(&f)
	db, errs := u.S.Data.Bdd.Connect()
	if errs != nil {
		return
	}
	var i bool
	g := `SELECT IF(COUNT(*),'true','false') FROM forgot WHERE uuid = ?`
	err := db.QueryRow(g, f.Uuid).Scan(&i)
	if err != nil {
		println(err.Error())
	}
	println(i)
}
