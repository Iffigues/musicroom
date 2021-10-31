package user

import (
	"github.com/iffigues/musicroom/server"
	"github.com/iffigues/musicroom/util"
	"github.com/iffigues/musicroom/regex"

	"github.com/satori/go.uuid"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"

	"strconv"
	"encoding/json"
	"net/http"
	"os"
	"log"
	"fmt"
	"context"
	"text/template"
)

type UserUtils struct {
	S	*server.Server
	Client	*redis.Client
}

type App struct {
	Uid	string
}

type Fo struct {
	Roi	int `json:"resource_owner_id"`
	Scopes	[]string `json:"scopes"`
	Eis	int `json:"expires_in_seconds"`
	App	App `json:"application"`
	Ca	int64 `json:"created_at"`
}

type User struct {
	Email		string `json:"email"`
	Uid		string `json:"uuid"`
	Password	string `json:"pwd"`
	Name		string `json:"name"`
	Types		bool
	Buy		bool   `json:"buy"`
	MailVerif	bool
	TokenEmail	uuid.UUID
}

type Four struct {
	Email	string `json:"email"`
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
	Password, err := util.Crypte(login.Password)
	if err != nil {
		println(err.Error())
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	login.Password = string(Password)
	if err := u.AddUser(&login); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	c.JSON(200, gin.H{"status": "OK"})
}



func (u *UserUtils)GetUsers(c *gin.Context) {
	var login User
	c.BindJSON(&login)
	t := login.Password
	if err := u.GetUser(&login); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	if !util.Decrypt([]byte(login.Password), []byte(t)) {
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
			t, err := template.ParseFiles("./template/accept.tmpl")
			if err == nil {
				t.Execute(c.Writer, nil)
				return
			}
			//c.HTML(200, "template/accept.tmpl", nil)
			//return
		}
	}
	c.JSON(400,gin.H{"match":"false"})
}

func (u *UserUtils) Get42(c *gin.Context) {
	cc, err := u.S.Data.Api["42"].NewClient()
	if err != nil {
		log.Fatal(err)
	}
	cc.Types = 1;
	println(cc.GetURL())
}

func (u *UserUtils) Tok(c *gin.Context) {
	target := &Fo{}
	gg := &Four{}
	cc, err := u.S.Data.Api["42"].NewClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	cc.Code = c.Query("code")
	cc.Types = 1
	err = cc.Authenticate(ctx)
	if err != nil {
		println(err.Error())
		return
	}
	resp, errs :=  cc.Client.Get("https://api.intra.42.fr/oauth/token/info")
	if errs != nil {
		fmt.Println(err)
	}
	if  err = json.NewDecoder(resp.Body).Decode(&target); err != nil {
	}
	rr, ee := cc.Client.Get("https://api.intra.42.fr/v2/users/" +  strconv.Itoa(target.Roi))
	if ee != nil {
		fmt.Println(err)
	}
	if err = json.NewDecoder(rr.Body).Decode(&gg); err != nil {
		fmt.Println(err)
	}
	fmt.Println(target, gg)
}

func (u *UserUtils) WWW(s *server.Server) {
	s.NewR("/user/signup", "user", "POST", 1, u.S.MakeMe(u.UserHandler))
	s.NewR("/user/signin", "users", "POST", 1, u.S.MakeMe(u.GetUsers))
	s.NewR("/user/signout","deluser", "GET", 1, u.S.MakeMe(TokenAuthMiddleware ,u.DelUser))
	s.NewR("user/verif/:token", "verifuser", "GET", 1, u.S.MakeMe(u.UserVerif))
	s.NewR("/user/42", "get42", "GET", 1, u.S.MakeMe(u.Get42))
	s.NewR("/user/token", "g42", "GET", 1, u.S.MakeMe(u.Tok))
	s.NewR("/user/friend/add", "fa", "POST", 1, u.S.MakeMe(u.AddFriend))
	s.NewR("/user/friend/accept", "fi", "POST", 1, u.S.MakeMe(u.AcceptFriend))
	s.NewR("/user/friend/refuse", "fd", "POST", 1, u.S.MakeMe(u.RefuseFriend))
	s.NewR("/user/friend/get", "ge", "GET", 1, u.S.MakeMe(u.ShowFriend))
	s.NewR("/user/friend/all", "al", "GET", 1 , u.S.MakeMe(u.GetAllFriend))
}
