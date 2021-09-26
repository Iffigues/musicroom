package user

import (
	"github.com/iffigues/musicroom/server"
	"github.com/gin-gonic/gin"
)

type User struct {
	Email	string `json:"email"`
	Password string `json:"pwd"`
	Types	bool
	MailVerif bool
}

func NewUser(s *server.Server) (u *User) {
	u = new(User)
	return
}

func (u *User)UserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var login User
		c.BindJSON(&login)
		login.Types = true
		c.JSON(200, gin.H{"status": "OK"})
	}
}

func (u *User) WWW(s *server.Server) {
	s.NewR("/user/login", "user", "POST", u.UserHandler, 1,  nil)
}
