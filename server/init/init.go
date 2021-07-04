package inits

import(
	"log"
	"github.com/go-ini/ini"
)

type Init struct {
	file *ini.File
}

func NewInit(path string) (i *Init, err error) {
	i = new(Init)
	i.file, err  = ini.Load(path);
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (i *Init) GetKey(section , key string)(inu string){
	ar,err:=  i.file.Section(section).GetKey(key)
	if err != nil {
		log.Panic(err)
	}
	return ar.String()
}
