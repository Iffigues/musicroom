package server

import (
	"regexp"
	"net/http"
	"strings"
	"github.com/iffigues/musicroom/config"
	"github.com/iffigues/musicroom/pk"
	//"github.com/gorilla/mux"
	"github.com/gin-gonic/gin"
	//"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	//"github.com/gorilla/sessions"
)

type HH interface {
	WWW(*Server)
}

type Data struct {
	Store cookie.Store
	Bdd  *pk.Pk
	Conf  *config.Conf
}

type Server struct {
	Router *gin.Engine
//	Router *gin.RouterGroup
	Data   *Data
	Handle map[string]*Handle
	Give   []HH
}

type Handle struct {
	Role	int
	Route	string
	Method	string
	Handle	func () gin.HandlerFunc
	Mi	[]func(c *gin.Context)
}

func (s *Server) AddHH(p ...HH) {
	for _, val := range p {
		s.Give = append(s.Give, val)
	}
}

func (s *Server) StartHH() {
	for _, val := range s.Give {
		val.WWW(s)
	}
}

func NewServer(i *config.Conf) (a *Server) {
	//router := mux.NewRouter()
	router := gin.Default()
	//router.StrictSlash(true)
	return &Server{
		Data: &Data{
			Store: cookie.NewStore([]byte("secret")),
			Bdd: pk.NewPk(*i),
			Conf:  i,
		},
		Router: router,
		Handle: make(map[string]*Handle),
	}
}

func (r *Server) NewR(route, key string, method string, handler func () gin.HandlerFunc , i int, mi []func(c *gin.Context)) {
	route = strings.ToLower(route)
	r.Handle[key] = &Handle{Method: method, Route: route, Handle: handler, Role: i, Mi: mi}
}


func (s *Server) Middleware(next http.Handler, a *Handle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
	})
}


func UserHandler(c *gin.Context) {
    r, err := regexp.Compile(`[a-zA-Z0-9]`)
    if err != nil {
       panic(err)
       return
    }
    username := c.Param("regex")
    if r.MatchString(username) == true {
        c.JSON(200,gin.H{"match":"true"})
    } else {
        c.JSON(400,gin.H{"match":"false"})
    }
}

func (g *Server) Servers() (srv *http.Server) {
	g.StartHH()
	g.Router.GET("/users/:regex/lol",UserHandler)
	for _, h := range g.Handle {
		print(h)
		mi := g.Router.Handle(h.Method,h.Route, h.Handle())
		for _, v := range h.Mi {
			mi.Use(v)
		}
		//g.Router.Handle(h.Route, g.Middleware(h.Handle, h)).Methods(h.Method...)
	}
	return &http.Server{
		Addr:    g.Data.Conf.GetValue("http","socket").(string),
		Handler: g.Router,
	}
}
