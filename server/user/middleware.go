package user

import(
	"fmt"
	"github.com/gin-gonic/gin"
)

func (u *UserUtils)DummyMiddleware() gin.HandlerFunc {
  // Do some initialization logic here
  // Foo()
  return func(c *gin.Context) {
	 fmt.Println("zzzz")
	  c.Next()
  }
}

