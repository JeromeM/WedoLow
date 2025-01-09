package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRandomUsers(t *testing.T) {
	tests := []struct {
		name       string
		count      int
		gender     string
		mockServer *httptest.Server
		wantErr    bool
	}{
		{
			name:   "valid response",
			count:  1,
			gender: "any",
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, `{"results":[{"name":{"first":"John","last":"Doe"},"gender":"male","email":"john.doe@example.com","phone":"123-456-7890"}]}`)
			})),
			wantErr: false,
		},
		{
			name:   "invalid response",
			count:  1,
			gender: "any",
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.mockServer.Close()
			client := NewRandomUserClient(tt.mockServer.URL)
			_, err := client.GetRandomUsers(tt.count, tt.gender)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRandomUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
