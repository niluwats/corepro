package dto

import (
	"fmt"
	"core/domain"
)

type NewUserRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type NewUserResponse struct {
	Email   string `json:"email,omitempty"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

//	ToUserDto gets a struct of User and form it into struct of type *NewUserResponse and returns it.
func ToUserDto(us *domain.User, msg string, status bool) *NewUserResponse {
	var response NewUserResponse
	if us == nil {
		response = NewUserResponse{
			Message: msg,
			Status:  status,
		}
	} else {
		response = NewUserResponse{
			Email:   us.Email,
			Message: msg,
			Status:  status,
		}
	}
	fmt.Println(response)
	return &response
}
