package model

type CreateUsersResponse struct {
	RequestedCount int    `json:"requestedCount"`
	Gender         string `json:"requestedGender"`
}
