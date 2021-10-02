package dto

import "core/domain"

type GetFileResponse struct {
	Path    string `json:"path,omitempty"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func ToFileResponceDto(path, message string, status bool) *GetFileResponse {
	var response GetFileResponse
	if path == "" {
		response = GetFileResponse{
			Message: message,
			Status:  status,
		}
	} else {
		response = GetFileResponse{
			Message: message,
			Status:  status,
			Path:    path,
		}
	}
	return &response
}

type FileUploadStatusResponse struct {
	BirthRegImg      bool   `json:"birth_reg_certificate"`
	NICImg           bool   `json:"nic"`
	PostalIdImg      bool   `json:"postal_id"`
	DrivingLicensImg bool   `json:"driving_license"`
	PassportImg      bool   `json:"passport"`
	Status           bool   `json:"status"`
	Message          string `json:"message"`
}

func ToFileStatusDto(us *domain.FileUploadStatus, msg string, status bool) *FileUploadStatusResponse {
	var response FileUploadStatusResponse
	if us == nil {
		response = FileUploadStatusResponse{
			Message: msg,
			Status:  status,
		}
	} else {
		response = FileUploadStatusResponse{
			BirthRegImg:      us.BirthRegImg,
			NICImg:           us.NICImg,
			PostalIdImg:      us.PostalIdImg,
			DrivingLicensImg: us.DrivingLicensImg,
			PassportImg:      us.PassportImg,
			Message: msg,
			Status: status,
		}
	}
	return &response
}
