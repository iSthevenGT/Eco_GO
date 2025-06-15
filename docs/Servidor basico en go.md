
# Servidor basico en go

```go
ackage main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello!")
	})
	r.Run(":8080")
}

```

