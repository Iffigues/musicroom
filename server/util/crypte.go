package util

import(
	"golang.org/x/crypto/bcrypt"
	"github.com/satori/go.uuid"
)

func Crypte(pass string)(io []byte,err error){
	return bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
}

func Decrypt(a []byte,b []byte)(vrai bool){
	if err := bcrypt.CompareHashAndPassword(a, b); err != nil {
		return false
	}
	return true
}

func Uid() (e uuid.UUID) {
	return uuid.NewV4()
}
