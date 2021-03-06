package server

import (
	"regexp"
	"github.com/gin-contrib/cors"
	"net/http"
	"strings"
	"github.com/iffigues/musicroom/config"
	"github.com/iffigues/musicroom/pk"
	"github.com/iffigues/musicroom/ws"
	"github.com/gin-gonic/gin"
	"github.com/iffigues/musicroom/api"
	"os"
	"io"
	"sync"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/iffigues/musicroom/worker"
	"github.com/gin-contrib/sessions/cookie"
)

/*
	interface for allow plugin to add route.
	the plugin structure have to have this function instance or you will get error in compilation
	example:
	type Man struct {
	}
	func (m *Man) WWW(s *Server) {
			s.NewR( path , route name , <"GET","POST", ...>, integer role (anonyme, user, admin) , u.S.MakeMe(m.MyHandlerFunc))
	}
*/
type HH interface {
	WWW(*Server)
}


type Mumu struct {
	Mu sync.RWMutex
	Wg sync.WaitGroup
}

/*
	Data structure contain all configuration and tool for the server and plugin
	Store: session
	Bdd: sql 
	Conf: configuration 
	api: tool for make api request
*/
type Data struct {
	Store cookie.Store
	Bdd  *pk.Pk
	Conf  *config.Conf
	Api map[string]*api.Config
	Lock map[int]*Mumu
	Worker *worker.Worker
	Ws  *ws.Channel
}

/*
	server structure is the server object.
	contains Router who is the template http motor (router, session, request)
	Data is the data use by server and plugin
	Handle is the route for the server
	Give is the array interface for allow plugin add Handle
*/
type Server struct {
	Router *gin.Engine
	Data   *Data
	Handle map[string]*Handle
	Give   []HH
}

/*
	Handl structure is the route use by server
*/

type Handle struct {
	Role	int
	Route	string
	Method	string
	Handle	[]gin.HandlerFunc
}

func (s *Server) AddPk (a *pk.Pk) {
	s.Data.Bdd = a
}

func (s *Server) AddHH(p ...HH) {
	for _, val := range p {
		s.Give = append(s.Give, val)
	}
}

func (s *Server)MakeMe(a ...gin.HandlerFunc)(b []gin.HandlerFunc) {
	return a
}

func (s *Server) StartHH() {
	for _, val := range s.Give {
		val.WWW(s)
	}
}

func GinConfig(i *config.Conf) {
	DebugMode := "debug"
	ReleaseMode := "release"
	TestMode := "test"
	var tt string
	var ttt string
	t := i.GetValue("gin", "mode")
	if t == nil {
		ttt = DebugMode
	} else {
		ttt = t.(string)
	}
	if DebugMode == ttt {
		gin.SetMode(gin.DebugMode)
		tt = "./log/gin-debug.log"
	} else if ReleaseMode == ttt {
		gin.SetMode(gin.ReleaseMode)
		tt = "./log/gin-release.log"
	} else if TestMode == ttt {
		gin.SetMode(gin.TestMode)
		tt = "./log/gin-test.log"
	} else {
		panic("gin mode unknown: " + tt + " (available mode: debug release test)")
	}
	gin.DisableConsoleColor();
	f, _ := os.Create(tt)
	gin.DefaultWriter = io.MultiWriter(f)
}

func (s *Server) FourTwo() {
	ap:= &api.Config{
		Host:  "https://api.insee.fr/",
		Oauth: api.Oauth{},
		Headers: map[string]string{
			"grant_type": "client_credentials",
		},
	}
	ap.Oauth.ClientID = "86023b24c48480f95e5b24b5a0d90939815fe16781adea9eb04ab34d3537b026"
	ap.Oauth.ClientSecret = "2411b140d3f8fb889e27dd89236f8008bb03779e1ba2333693c38a38e9bcb33c"
	ap.Oauth.TokenURL = "https://api.intra.42.fr/oauth/token"
	ap.Oauth.AuthURL = "https://api.intra.42.fr/oauth/authorize"
	ap.Oauth.AuthParam = map[string]string{
		"grant_type": "client_credentials",
	}
	ap.Oauth.RedirectURL = "http://gopiko.fr:9000/user/token"
	s.Data.Api["42"] = ap
}

func NewServer(i *config.Conf) (a *Server) {
	GinConfig(i)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
    }))
	a = &Server{
		Data: &Data{
			Store: cookie.NewStore([]byte("secret")),
			Bdd: pk.NewPk(*i),
			Conf:  i,
			Api: make(map[string]*api.Config),
			Lock: make(map[int]*Mumu),
			Worker: worker.NewWorkers(10),
			Ws: ws.NewChannel(),
		},
		Router: router,
		Handle: make(map[string]*Handle),
	}
	router.Use(sessions.Sessions("mysession", a.Data.Store))
	a.FourTwo()
	return
}

func (r *Server) AddLock(i int) (err error) {
	if _, ok := r.Data.Lock[i]; ok {
		return fmt.Errorf("already exists")
	}
	var mu  sync.RWMutex
	var wg sync.WaitGroup
	r.Data.Lock[i] = &Mumu{
		Mu: mu,
		Wg: wg,
	}
	return nil
}

func (r *Server) DelLock(i int) (err error) {
	if _, ok := r.Data.Lock[i]; !ok {
		return fmt.Errorf("don't exists")
	}
	delete(r.Data.Lock, i)
	return nil
}

func (r *Server) NewR(route, key string, method string,i int, handler []gin.HandlerFunc) {
	route = strings.ToLower(route)
	r.Handle[key] = &Handle{Method: method, Route: route, Handle: handler, Role: i}
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
	g.Mm()
	for _, h := range g.Handle {
		if h.Method == "GET" {
			g.Router.GET(h.Route, h.Handle...)
		}
		if h.Method == "POST" {
			g.Router.POST(h.Route, h.Handle...)
		}
	}
	return &http.Server{
		Addr:    g.Data.Conf.GetValue("http","socket").(string),
		Handler: g.Router,
	}
}

func (g *Server) Mm() {
	g.Data.Ws.NewChan("zzz")
}
