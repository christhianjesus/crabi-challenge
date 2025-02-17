package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/christhianjesus/crabi-challenge/internal/domain"
)

type pldRepository struct {
	client HTTPClient
	url    string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type pldRequest struct {
	*domain.User
}

func (p *pldRequest) MarshalJSON() ([]byte, error) {
	if p.User == nil {
		return nil, errors.New("Nil user")
	}

	return json.Marshal(map[string]string{
		"first_name": p.FirstName,
		"last_name":  p.LastName,
		"email":      p.Email,
	})
}

type pldResponse struct {
	IsInBlacklist bool `json:"is_in_blacklist"`
}

func NewPLDRepository(client HTTPClient, url string) domain.PLDRepository {
	return &pldRepository{client, url}
}

func (ms *pldRepository) IsValidUser(ctx context.Context, user *domain.User) (bool, error) {
	data, err := json.Marshal(&pldRequest{user})
	if err != nil {
		return false, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ms.url+"/check-blacklist", bytes.NewReader(data))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	if resp, err = ms.client.Do(req); err != nil {
		return false, err
	}

	defer resp.Body.Close()

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return false, err
	}

	var response pldResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return false, err
	}

	return !response.IsInBlacklist, nil
}
