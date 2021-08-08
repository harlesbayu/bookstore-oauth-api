package access_token

import (
	"strings"
	"time"

	"github.com/harlesbayu/bookstore_oauth-api/src/utils/errors"
)

const (
	expirationTime = 24
)

type AccessToken struct {
	AccessToken string `json:"accessToken"`
	UserId      int64  `json:"userId"`
	ClientId    int64  `json:"clientId"`
	Expires     int64  `json:"expires"`
}

func (data *AccessToken) Vaidate() *errors.RestErr {
	data.AccessToken = strings.TrimSpace(data.AccessToken)

	if data.AccessToken == "" {
		return errors.NewBadRequestError("invalid access token")
	}

	if data.UserId <= 0 {
		return errors.NewBadRequestError("invalid user id")
	}

	if data.ClientId <= 0 {
		return errors.NewBadRequestError("invalid client id")
	}

	if data.Expires <= 0 {
		return errors.NewBadRequestError("invalid expiration time")
	}

	return nil
}

func GetNewAccessToken() AccessToken {
	return AccessToken{
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}
