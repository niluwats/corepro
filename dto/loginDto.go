package dto

import "core/domain"

type Response struct {
}

type NewLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewLoginResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	Email       string `json:"email,omitempty"`
	UserId      string `json:"user_id,omitempty"`
	Message     string `json:"message,omitempty"`
	Status      bool   `json:"status"`
}

func ToULoginDto(u *domain.User, msg string, status bool) *NewLoginResponse {
	var response NewLoginResponse
	if u == nil {
		response = NewLoginResponse{
			Message: msg,
			Status:  status,
		}
	} else {
		response = NewLoginResponse{
			Email:   u.Email,
			UserId:  u.Id,
			Message: msg,
			Status:  status,
		}
	}
	return &response
}
