package user

import (
	"github.com/iffigues/musicroom/server"

	"github.com/satori/go.uuid"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"

	"net/http"
	"fmt"
	"os"
)

type UserUtils struct {
	S *server.Server
	Client *redis.Client
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

type AccessDetails struct {
    AccessUuid string
    UserId   string
}

func NewUser(s *server.Server) (u *UserUtils) {
	u = new(UserUtils)
	u.S = s
	u.InitUser()
	var  client *redis.Client
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	u.Client = client
	return
}

func (u *UserUtils)UserHandler(c *gin.Context) {
	var login User
	c.BindJSON(&login)
	if err := u.AddUser(&login); err != nil {
		fmt.Println(err)
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
	ts, err := u.CreateToken(login, true)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	u.CreateAuth(login.Uid, ts)
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

func (u *UserUtils)DelUser(c *gin.Context) {
	au, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleted, delErr := u.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

func (u *UserUtils) WWW(s *server.Server) {
	s.NewR("/user/signup", "user", "POST", 1, u.S.MakeMe(u.UserHandler))
	s.NewR("/user/signin", "users", "POST", 1, u.S.MakeMe(u.GetUsers))
	s.NewR("/user/signout","deluser", "GET", 1, u.S.MakeMe(TokenAuthMiddleware ,u.DelUser))
}
