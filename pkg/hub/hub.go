package hub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/volker-raschek/docker-hub-description-updater/pkg/types"
)

var (
	dockerHubAPI = "https://hub.docker.com/v2"
)

type Hub struct {
	client      *http.Client
	credentials *types.LoginCredentials
	token       *types.Token
}

// GetRepository returns a repository struct
func (h *Hub) GetRepository(namespace string, name string) (*types.Repository, error) {

	if len(namespace) <= 0 {
		return nil, errorNoNamespaceDefined
	}

	if len(name) <= 0 {
		return nil, errorNoRepositoryDefined
	}

	rawURL := fmt.Sprintf("%v/repositories/%v/%v", dockerHubAPI, namespace, name)
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", errorFailedToParseURL, err)
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", errorFailedToCreateRequest, err)
	}

	if h.token == nil {
		token, err := h.getToken()
		if err != nil {
			return nil, err
		}
		h.token = token
	}
	req.Header.Add("Authorization", fmt.Sprintf("JWT %v", h.token.Token))

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", errorFailedToSendRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v: expect %v, received %v", errorUnexpectedHTTPStatuscode, http.StatusOK, resp.StatusCode)
	}

	repository := new(types.Repository)
	jsonDecoder := json.NewDecoder(resp.Body)
	if err := jsonDecoder.Decode(repository); err != nil {
		return nil, fmt.Errorf("%v: %v", errorFailedToParseJSON, err)
	}

	return repository, nil
}

func (h *Hub) getToken() (*types.Token, error) {
	loginBuffer := new(bytes.Buffer)
	jsonEncoder := json.NewEncoder(loginBuffer)
	if err := jsonEncoder.Encode(h.credentials); err != nil {
		return nil, fmt.Errorf("%v: %v", errorFailedToParseJSON, err)
	}

	rawURL := fmt.Sprintf("%v/users/login/", dockerHubAPI)
	req, err := http.NewRequest(http.MethodPost, rawURL, loginBuffer)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", errorFailedToCreateRequest, err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", errorFailedToCreateRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v: expect %v, received %v", errorUnexpectedHTTPStatuscode, http.StatusOK, resp.StatusCode)
	}

	token := new(types.Token)
	jsonDecoder := json.NewDecoder(resp.Body)
	if err := jsonDecoder.Decode(token); err != nil {
		return nil, fmt.Errorf("%v: %v", errorFailedToParseJSON, err)
	}

	return token, nil
}

// PatchRepository updates the docker hub repository
func (h *Hub) PatchRepository(repository *types.Repository) (*types.Repository, error) {

	if len(repository.Namespcace) <= 0 {
		return nil, errorNoNamespaceDefined
	}

	if len(repository.Name) <= 0 {
		return nil, errorNoRepositoryDefined
	}

	if h.token == nil {
		token, err := h.getToken()
		if err != nil {
			return nil, err
		}
		h.token = token

	}

	rawURL := fmt.Sprintf("%v/repositories/%v/%v", dockerHubAPI, repository.Namespcace, repository.Name)
	patchURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", errorFailedToParseURL, err)
	}

	data := &url.Values{}
	data.Set("full_description", repository.FullDescription)
	patchURL.RawQuery = data.Encode()

	req, err := http.NewRequest(http.MethodPatch, patchURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", errorFailedToCreateRequest, err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("JWT %v", h.token.Token))

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", errorFailedToCreateRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v: expect %v, received %v", errorUnexpectedHTTPStatuscode, http.StatusOK, resp.StatusCode)
	}

	patchedRepository := new(types.Repository)
	jsonDecoder := json.NewDecoder(resp.Body)
	if err := jsonDecoder.Decode(h.token); err != nil {
		return nil, fmt.Errorf("%v: %v", errorFailedToParseJSON, err)
	}

	return patchedRepository, nil
}

func New(credentials *types.LoginCredentials) *Hub {
	return &Hub{
		client: &http.Client{
			Timeout: time.Second * 15,
		},
		credentials: credentials,
	}
}
