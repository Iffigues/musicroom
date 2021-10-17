package postmail

import (
	"github.com/iffigues/musicroom/config"

	"html/template"
	"net/smtp"
	"bytes"
)

type Email struct {
	From string
	Pwd string
	To []string
	User string
	Mime  string
	SmtpHost string
	SmtpPort string
	Message string
	Auth smtp.Auth
}

func (e *Email) Users(a string) {
	e.User = a
}

func NewEmail(Conf  *config.Conf) (*Email) {
	if Conf != nil {
		return &Email{
			From: Conf.GetValue("email","from").(string),
			Pwd: Conf.GetValue("email","pwd").(string),
			User: Conf.GetValue("email","user").(string),
			Mime: `MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n`,
			SmtpHost: Conf.GetValue("email","smtphost").(string),
			SmtpPort: Conf.GetValue("email","smtpport").(string),
		}
	}
	return &Email{
	}
}

func (e *Email) Froms(from string) {
	e.From = from
}

func (e *Email) Pwds (pwd string) {
	e.Pwd = pwd
}

func (e *Email) Tos(to []string) {
	e.To = to
}

func (e *Email) AddTos (to string) {
	e.To = append(e.To, to)
}

func (e *Email) SmtpHosts (h string) {
	e.SmtpHost = h
}

func (e *Email) SmtpPorts (h string) {
	e.SmtpPort = h
}

func (e *Email) Messages(a string) {
	e.Message = a
}

func (e *Email) Mimes(b string) {
	e.Mime = b
}

func (e *Email) Auths() {
	e.Auth = smtp.PlainAuth("", e.User, e.Pwd, e.SmtpHost)
}

func (e *Email) Html(path string, Data interface{}) (err error){
	t, errs := template.ParseFiles(path)
	if errs != nil {
		return errs
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, Data); err != nil {
		return err
	}
	e.Message = buf.String()
	return nil
}

func (e *Email) Send() error {
	println(e.Mime)
	message := []byte("Subject: " + e.User + "!\n" + "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" + e.Message)
	return smtp.SendMail(e.SmtpHost+":"+e.SmtpPort, e.Auth, e.From, e.To, message)
}
