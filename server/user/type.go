package user

import (
	"github.com/iffigues/musicroom/server"
	"github.com/iffigues/musicroom/util"
	"github.com/iffigues/musicroom/regex"

	"github.com/satori/go.uuid"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"

	"net/http"
	"os"
	"log"
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
		Addr: dsn,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	u.Client = client
	return
}

func (us *UserUtils)verify(users *User) (vrai bool) {
	if regex.Verifie.IsEmail(users.Email) {
		if regex.Verifie.ValidPassword(users.Password) {
			return true
		}
	}
	return false
}

func (u *UserUtils)UserHandler(c *gin.Context) {
	var login User
	c.BindJSON(&login)

	if !u.verify(&login) {
		c.JSON(403, gin.H{"status":"bading"})
		return
	}
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
	if delErr != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

func (u *UserUtils)UserVerif(c *gin.Context) {
	username := c.Param("token")
	if util.IsValidUUID(username) == true {
		if err := u.GetUseriVerif(username); err == nil {
			c.JSON(200,gin.H{"match":"true"})
			return
		}
	}
	c.JSON(400,gin.H{"match":"false"})
}

func (u *UserUtils) Get42(c *gin.Context) {
	cc, err := u.S.Data.Api["42"].NewClient()
	if err != nil {
		log.Fatal(err)
	}
	cc.Types = 2;
	println(cc.GetURL())
}

func (u *UserUtils) WWW(s *server.Server) {
	s.NewR("/user/signup", "user", "POST", 1, u.S.MakeMe(u.UserHandler))
	s.NewR("/user/signin", "users", "POST", 1, u.S.MakeMe(u.GetUsers))
	s.NewR("/user/signout","deluser", "GET", 1, u.S.MakeMe(TokenAuthMiddleware ,u.DelUser))
	s.NewR("user/verif/:token", "verifuser", "GET", 1, u.S.MakeMe(u.UserVerif))
	s.NewR("/user/42", "get42", "GET", 1, u.S.MakeMe(u.Get42))
}
