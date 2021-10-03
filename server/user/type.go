package user

import (
	"github.com/iffigues/musicroom/server"

	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

type UserUtils struct {
	S *server.Server
}

type User struct {
	Email	string `json:"email"`
	Uid	string `json:"uuid"`
	Password string `json:"pwd"`
	Types	bool
	Buy	bool   `json:"buy"`
	MailVerif bool
	TokenEmail uuid.UUID
}

func NewUser(s *server.Server) (u *UserUtils) {
	u = new(UserUtils)
	u.S = s
	u.InitUser()
	return
}

func (u *UserUtils)UserHandler(c *gin.Context) {
	var login User
	c.BindJSON(&login)
	if err := u.AddUser(&login); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	c.JSON(200, gin.H{"status": "OK"})
}

func (u *UserUtils)GetUsers(c *gin.Context) {
	var login User
	c.BindJSON(&login)
	if err := u.GetUser(&login); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	c.JSON(200, gin.H{"status": "OK"})
}

func (u *UserUtils) WWW(s *server.Server) {
	s.NewR("/user/signin", "user", "POST", 1, []gin.HandlerFunc{u.UserHandler})
	s.NewR("/user/signup", "users", "POST", 1, []gin.HandlerFunc{u.DummyMiddleware(), u.GetUsers})
}
