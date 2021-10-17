package song

import (
    "testing"
	"github.com/iffigues/musicroom/servertest"
)

func TestAddSong(t *testing.T) {
	s := servertest.Serves();

	b := &Song{Name:"Un", Author:"JC", Ranking:10, IsPlayed:false}

	c := NewSong(s)
	err := c.AddSong(b)
	if err != nil {
		t.Fatalf(err.Error())
	}

	s.AddHH(c)
	servertest.LanceServe(s)
}
