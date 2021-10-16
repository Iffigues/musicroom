package user

import(
	"net/http"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(c *gin.Context) {
     err := TokenValid(c.Request)
     if err != nil {
        c.JSON(http.StatusUnauthorized, err.Error())
        c.Abort()
        return
     }
     c.Next()
 }
