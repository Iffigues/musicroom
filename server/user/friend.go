package user

import (
	"github.com/gin-gonic/gin"	
	"fmt"
)

func (u *UserUtils) AddFriend(c *gin.Context) {
	fmt.Println(ExtractTokenMetadata(c.Request))
}
