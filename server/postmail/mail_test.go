package postmail

import (
    "testing"
)

func TestHelloEmpty(t *testing.T) {
	e := NewEmail(nil)
	e.Froms("denoyelle.boris@gmail.com")
	e.Users("denoyelle.boris@gmail.com")
	e.Pwds("Datura1426")
	e.SmtpHosts("smtp.gmail.com")
	e.AddTos("denoyelle.boris@gmail.com")
	e.SmtpPorts("587")
	ee := "From: denoyelle.boris@gmail.com\r\n" +
        "To: roger.denoyelle.boris\r\n" +
        "Subject: Test mail\r\n\r\n" +
        "Email body\r\n"
	e.Messages(ee)
	e.Auths()
	if err := e.Send(); err != nil {
		t.Fatal(err)
	}
}
