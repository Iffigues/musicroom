package user

import (
	"github.com/iffigues/musicroom/util"
	"github.com/gin-gonic/gin"
	"github.com/iffigues/musicroom/postmail"
	"github.com/iffigues/musicroom/regex"
	"fmt"
	"text/template"
)

type Ml struct {
	Email	string	`json:"email"`
	Uuid	string	`json:"uid"`
	Pwd	string	`json:"pwd"`
	Pwda	string	`json:"pwd1"`
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
	bb := util.Uid().String()
	g := `INSERT INTO forgot (user_id, uuid) VALUES ((SELECT id FROM user WHERE email = ?), ?)`
	_, err := db.Exec(g, f.Email, bb)
	if err != nil {
		fmt.Println(err)
	}
	u.SendMailForgot(bb, f.Email)
}

func (u *UserUtils) Change_Forgot(c *gin.Context) {
	token := c.Param("token")
	db, errs := u.S.Data.Bdd.Connect()
	if errs != nil {
		return
	}
	var i bool
	g := `SELECT IF(COUNT(*),'true','false') FROM forgot WHERE uuid = ?`
	err := db.QueryRow(g, token).Scan(&i)
	if err != nil {
		println(err.Error())
		return
	}
	t, err := template.ParseFiles("./template/forgot.tmpl")
	if err == nil {
		t.Execute(c.Writer, token)
		return
	}
}

func (u *UserUtils) Change(c *gin.Context) {
	token := c.Param("token")
	pwd1 := c.PostForm("pwd1")
	pwd := c.PostForm("pwd")
	if pwd != pwd1 {
		fmt.Println("not same")
		return
	}
	if !regex.Verifie.ValidPassword(pwd) {
		fmt.Println("oui")
		return
	}
	db, err := u.S.Data.Bdd.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	var i bool
	g := `SELECT IF(COUNT(*),'true','false') FROM forgot WHERE uuid = ?`
	err = db.QueryRow(g, token).Scan(&i)
	if err != nil {
		println(err.Error())
		return
	}
	g = `UPDATE user SET password = ? WHERE id = (SELECT user_id FROM forgot WHERE uuid = ?)`
	ii, tt := util.Crypte(pwd)
	if tt != nil {
		fmt.Println(tt)
		return
	}
	_, err = db.Query(g, ii, token)
	if err != nil {
		fmt.Println(err)
		return
	}
	g = "DELETE FROM forgot  WHERE uuid = ?"
	_, err =  db.Exec(g, token)
	if err != nil {
		fmt.Println(err)
	}
}
