package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/harlesbayu/bookstore-utils-go/rest_errors"
	"net/http"

	"github.com/harlesbayu/bookstore_oauth-api/src/domain/users"
)

var baseUrl = "http://localhost:3000"

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *rest_errors.RestErr)
}

type userRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &userRepository{}
}

func (r *userRepository) LoginUser(email string, password string) (*users.User, *rest_errors.RestErr) {
	requestBody, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	if err != nil {
		return nil, rest_errors.NewInternalServerError("error marshall body", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/users/login", baseUrl), "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, rest_errors.NewInternalServerError("error request when trying to login user", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, rest_errors.NewNotFoundError("user not found")
	}

	var user users.User
	json.NewDecoder(resp.Body).Decode(&user)

	return &user, nil
}
