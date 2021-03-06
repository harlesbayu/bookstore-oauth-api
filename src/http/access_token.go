package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harlesbayu/bookstore-utils-go/rest_errors"
	atDomain "github.com/harlesbayu/bookstore_oauth-api/src/domain/access_token"
	"github.com/harlesbayu/bookstore_oauth-api/src/services/access_token"
)

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Create(*gin.Context)
}

type accessTokenHandler struct {
	service access_token.Service
}

func NewAccessTokenHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (h *accessTokenHandler) GetById(c *gin.Context) {
	accessToken, err := h.service.GetById(c.Param("accessToken"))

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, accessToken)
}

func (h *accessTokenHandler) Create(c *gin.Context) {
	var request atDomain.AccessTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	accessToken, err := h.service.Create(request)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, accessToken)
}
