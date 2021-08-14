package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/harlesbayu/bookstore_oauth-api/src/domain/users"
	"github.com/harlesbayu/bookstore_oauth-api/src/utils/errors"
)

var baseUrl = "http://localhost:3000"

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type userRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &userRepository{}
}

func (r *userRepository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	requestBody, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	resp, err := http.Post(fmt.Sprintf("%s/users/login", baseUrl), "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, errors.NewInternalServerError("error request when trying to login user")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewNotFoundError("user not found")
	}

	var user users.User
	json.NewDecoder(resp.Body).Decode(&user)

	return &user, nil
}
