package app

import (
	"github.com/gin-gonic/gin"
	"github.com/harlesbayu/bookstore_oauth-api/src/domain/access_token"
	"github.com/harlesbayu/bookstore_oauth-api/src/http"
	"github.com/harlesbayu/bookstore_oauth-api/src/repository/db"
	"github.com/harlesbayu/bookstore_oauth-api/src/repository/rest"
)

var (
	router = gin.Default()
)

func StartApplication() {
	atService := access_token.NewService(rest.NewRepository(), db.NewRepository())
	atHandler := http.NewHandler(atService)

	router.GET("/oauth/access_token/:accessTokenId", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8080")
}
