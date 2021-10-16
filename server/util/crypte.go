package util

import(
	"golang.org/x/crypto/bcrypt"
	"github.com/satori/go.uuid"
	op "github.com/google/uuid"
)

func IsValidUUID(u string) bool {
	_, err := op.Parse(u)
	return err == nil
 }

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
