package room

import (
    "testing"
	"github.com/iffigues/musicroom/servertest"
)

func TestAddRoom(t *testing.T) {
	s := servertest.Serves();

	r := &Room{Name:"Un", CreatorId:1}

	c := NewRoom(s)
	err := c.AddRoom(r)
	if err != nil {
		t.Fatalf(err.Error())
	}

	s.AddHH(c)
}
//Golang http, ouvrir client
func TestGetRoom(t *testing.T) {
	s := servertest.Serves();

	r := &Room{Id:1}

	c := NewRoom(s)
	err := c.GetRoom(r)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if r.Name != "Un" {
		t.Fatalf(err.Error())
	}

	s.AddHH(c)
}
//DelRoom  DeleteRoom
func TestDeleteRoom(t *testing.T) {
	s := servertest.Serves();


	c := NewRoom(s)
	err := c.DeleteRoom(1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	s.AddHH(c)
	go servertest.LanceServe(s)
}
