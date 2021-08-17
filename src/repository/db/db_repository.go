package db

import (
	"github.com/gocql/gocql"
	"github.com/harlesbayu/bookstore-utils-go/rest_errors"
	"github.com/harlesbayu/bookstore_oauth-api/src/clients/cassandra"
	"github.com/harlesbayu/bookstore_oauth-api/src/domain/access_token"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	queryCreateAccessToken = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES(?,?,?,?);"
	queryUpdateExpires     = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *rest_errors.RestErr)
	Create(*access_token.AccessToken) *rest_errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *rest_errors.RestErr
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

type dbRepository struct{}

func (r *dbRepository) GetById(accessToken string) (*access_token.AccessToken, *rest_errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, accessToken).Scan(
		&result.AccessToken,
		&result.UserId,
		&result.ClientId,
		&result.Expires,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, rest_errors.NewNotFoundError("no access token found with given id")
		}
		return nil, rest_errors.NewInternalServerError("database error", err)
	}

	return &result, nil
}

func (r *dbRepository) Create(data *access_token.AccessToken) *rest_errors.RestErr {
	if err := cassandra.GetSession().Query(queryCreateAccessToken,
		data.AccessToken,
		data.UserId,
		data.ClientId,
		data.Expires).Exec(); err != nil {
		return rest_errors.NewInternalServerError("database error", err)
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(data access_token.AccessToken) *rest_errors.RestErr {
	if err := cassandra.GetSession().Query(queryUpdateExpires,
		data.Expires,
		data.AccessToken,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerError("database error", err)
	}
	return nil
}
