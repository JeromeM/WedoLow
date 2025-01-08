package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RandomUserClient struct {
	baseURL string
}

type RandomUserResponse struct {
	Results []struct {
		Name struct {
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Gender string `json:"gender"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
	} `json:"results"`
}

func NewRandomUserClient(baseURL string) *RandomUserClient {
	return &RandomUserClient{baseURL: baseURL}
}

func (c *RandomUserClient) GetRandomUsers(count int, gender string) (*RandomUserResponse, error) {
	url := fmt.Sprintf("%s?results=%d", c.baseURL, count)
	if gender != "any" {
		url = fmt.Sprintf("%s&gender=%s", url, gender)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RandomUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
