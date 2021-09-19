package logger

import (
    "os"
    "log"
)


type logs struct {
    file *os.File
    Warning *log.Logger
    Info    *log.Logger
    Error *log.Logger
    Fatals *log.Logger
}

func NewLog(aa string) (a logs){
	var err error
	a.file, err = os.OpenFile(aa, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	a.Info = log.New(a.file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	a.Warning = log.New(a.file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	a.Error = log.New(a.file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	a.Fatals = log.New(a.file, "Fatal: ", log.Ldate|log.Ltime|log.Lshortfile)
	return
}

func (aa logs)Fatal(a string) {
	aa.Fatals.Println(a)
	os.Exit(1)
}
