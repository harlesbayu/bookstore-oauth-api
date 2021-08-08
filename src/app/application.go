package app

import (
	"github.com/gin-gonic/gin"
	"github.com/harlesbayu/bookstore_oauth-api/src/clients/cassandra"
	"github.com/harlesbayu/bookstore_oauth-api/src/domain/access_token"
	"github.com/harlesbayu/bookstore_oauth-api/src/http"
	"github.com/harlesbayu/bookstore_oauth-api/src/repository/db"
)

var (
	router = gin.Default()
)

func StartApplication() {
	session, dbErr := cassandra.GetSession()

	if dbErr != nil {
		panic(dbErr)
	}

	defer session.Close()

	atService := access_token.NewService(db.NewRepository())
	atHandler := http.NewHandler(atService)

	router.GET("/oauth/access_token/:accessTokenId", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8080")
}
