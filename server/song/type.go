package song

import (
	"github.com/iffigues/musicroom/server"
	"github.com/gin-gonic/gin"
)

type SongUtils struct {
	S *server.Server
}

type Song struct {
	Name	string `json:"name"`
	Author	string `json:"author"`
	Ranking int `json:"ranking"`
	IsPlayed	bool   `json:"isplayed"`
}

func NewSong(s *server.Server) (u *SongUtils) {
	u = new(SongUtils)
	u.S = s
	u.InitSong()
	return
}

func (u *SongUtils)SongHandler(c *gin.Context) {
	var id Song
	c.BindJSON(&id)
	if err := u.AddSong(&id); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	c.JSON(200, gin.H{"status": "OK"})
}

func (u *SongUtils)GetSongs(c *gin.Context) {
	var id Song
	c.BindJSON(&id)
	if err := u.GetSong(&id); err != nil {
		c.JSON(400, gin.H{"status": "bad"})
		return
	}
	c.JSON(200, gin.H{"status": "OK"})
}

func (u *SongUtils) WWW(s *server.Server) {
	s.NewR("/song/add", "song", "POST", 1, u.S.MakeMe(u.SongHandler))
	s.NewR("/song/get", "songs", "POST", 1, u.S.MakeMe(u.GetSongs))
}
