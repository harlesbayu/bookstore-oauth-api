package access_token

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/harlesbayu/bookstore-utils-go/rest_errors"
)

const (
	expirationTime            = 24
	grantTypePassword         = "password"
	grantTypeClientCredential = "client_credential"
)

type AccessToken struct {
	AccessToken string `json:"accessToken"`
	UserId      int64  `json:"userId"`
	ClientId    int64  `json:"clientId"`
	Expires     int64  `json:"expires"`
}

type AccessTokenRequest struct {
	GrantType    string `json:"grantType"`
	Scope        string `json:"scope"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func (data *AccessTokenRequest) Validate() *rest_errors.RestErr {
	switch data.GrantType {
	case grantTypePassword:
		break
	case grantTypeClientCredential:
		break
	default:
		return rest_errors.NewBadRequestError("invalid grant_type parameter")
	}

	return nil
}

func (data *AccessToken) Validate() *rest_errors.RestErr {
	data.AccessToken = strings.TrimSpace(data.AccessToken)

	if data.AccessToken == "" {
		return rest_errors.NewBadRequestError("invalid access token")
	}

	if data.UserId <= 0 {
		return rest_errors.NewBadRequestError("invalid user id")
	}

	if data.ClientId <= 0 {
		return rest_errors.NewBadRequestError("invalid client id")
	}

	if data.Expires <= 0 {
		return rest_errors.NewBadRequestError("invalid expiration time")
	}

	return nil
}

func GetNewAccessToken(userId int64) *AccessToken {
	expires := time.Now().Add(time.Minute * 15).Unix()
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = expires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("secret-access-token"))
	if err != nil {
		return nil
	}
	return &AccessToken{
		AccessToken: token,
		UserId:      userId,
		Expires:     expires,
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}
