package access_token

import (
	"strings"

	"github.com/harlesbayu/bookstore_oauth-api/src/domain/users"
	"github.com/harlesbayu/bookstore_oauth-api/src/utils/errors"
)

type DbRepository interface {
	GetById(string) (*AccessToken, *errors.RestErr)
	Create(*AccessToken) *errors.RestErr
	UpdateExpirationTime(AccessToken) *errors.RestErr
}

type RestRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}
type Service interface {
	GetById(string) (*AccessToken, *errors.RestErr)
	Create(AccessTokenRequest) (*AccessToken, *errors.RestErr)
	UpdateExpirationTime(AccessToken) *errors.RestErr
}

type service struct {
	dbRepository   DbRepository
	restRepository RestRepository
}

func NewService(restRepo RestRepository, dbRepo DbRepository) Service {
	return &service{
		restRepository: restRepo,
		dbRepository:   dbRepo,
	}
}

func (s *service) GetById(accessTokenId string) (*AccessToken, *errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)

	if len(accessTokenId) == 0 {
		return nil, errors.NewBadRequestError("invalid access token id")
	}

	accessToken, err := s.dbRepository.GetById(accessTokenId)

	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (s *service) Create(request AccessTokenRequest) (*AccessToken, *errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	user, err := s.restRepository.LoginUser(request.Username, request.Password)

	if err != nil {
		return nil, errors.NewBadRequestError("invalid access token id")
	}

	accessToken := GetNewAccessToken(user.Id)

	if accessToken == nil {
		return nil, errors.NewBadRequestError("failed generate access token")
	}

	if err := s.dbRepository.Create(accessToken); err != nil {
		return nil, errors.NewBadRequestError("failed create access token")
	}

	return accessToken, nil
}

func (s *service) UpdateExpirationTime(data AccessToken) *errors.RestErr {
	if err := data.Validate(); err != nil {
		return err
	}

	return s.dbRepository.UpdateExpirationTime(data)
}
