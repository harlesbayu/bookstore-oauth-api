package access_token

import (
	"github.com/harlesbayu/bookstore-utils-go/rest_errors"
	"strings"

	"github.com/harlesbayu/bookstore_oauth-api/src/domain/access_token"
	"github.com/harlesbayu/bookstore_oauth-api/src/repository/db"
	"github.com/harlesbayu/bookstore_oauth-api/src/repository/rest"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
}

type service struct {
	restUsersRepo rest.RestUsersRepository
	dbRepo        db.DbRepository
}

func NewService(usersRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepo:        dbRepo,
	}
}

func (s *service) GetById(accessToken string) (*access_token.AccessToken, rest_errors.RestErr) {
	accessToken = strings.TrimSpace(accessToken)

	if len(accessToken) == 0 {
		return nil, rest_errors.NewBadRequestError("invalid access token id")
	}

	resp, err := s.dbRepo.GetById(accessToken)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	user, err := s.restUsersRepo.LoginUser(request.Username, request.Password)

	if err != nil {
		return nil, rest_errors.NewBadRequestError("invalid access token id")
	}

	accessToken := access_token.GetNewAccessToken(user.Id)

	if accessToken == nil {
		return nil, rest_errors.NewBadRequestError("failed generate access token")
	}

	if err := s.dbRepo.Create(accessToken); err != nil {
		return nil, rest_errors.NewBadRequestError("failed create access token")
	}

	return accessToken, nil
}

func (s *service) UpdateExpirationTime(data access_token.AccessToken) rest_errors.RestErr {
	if err := data.Validate(); err != nil {
		return err
	}

	return s.dbRepo.UpdateExpirationTime(data)
}
