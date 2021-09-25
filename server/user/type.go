package user

import (
	"github.com/iffigues/musicroom/server"
	"github.com/gin-gonic/gin"
)

type User struct {
	email	string `json:"email"`
	password string `json:"pwd"`
	types	bool
	mailVerif bool
}

func NewUser(s *server.Server) (u *User) {
	u = new(User)
	return
}

func (u *User)UserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	}
}

func (u *User) WWW(s *server.Server) {
	s.NewR("/user/login", "user", "GET", u.UserHandler, 1,  nil)
}
