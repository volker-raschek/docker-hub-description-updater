package hub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/volker-raschek/docker-hub-description-updater/pkg/types"
)

var (
	dockerHubAPI = "https://hub.docker.com/v2"
)

func GetRepository(namespace string, name string, token *types.Token) (*types.Repository, error) {

	if len(namespace) <= 0 {
		return nil, errorNoNamespaceDefined
	}

	if len(name) <= 0 {
		return nil, errorNoRepositoryDefined
	}

	client := new(http.Client)

	url, err := url.Parse(fmt.Sprintf("%v/repositories/%v/%v", dockerHubAPI, namespace, name))
	if err != nil {
		return nil, fmt.Errorf("Can not prase URL: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Can not create request to get repository: %v", err)
	}

	if token != nil {
		req.Header.Add("Authorization", fmt.Sprintf("JWT %v", token.Token))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("An error has occured: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid HTTP-Statuscode: Get %v but expect 200", resp.StatusCode)
	}

	repository := new(types.Repository)
	jsonDecoder := json.NewDecoder(resp.Body)
	if err := jsonDecoder.Decode(repository); err != nil {
		return nil, fmt.Errorf("Can not encode JSON from Repository struct: %v", err)
	}

	return repository, nil
}

func GetToken(loginCredentials *types.LoginCredentials) (*types.Token, error) {

	if len(loginCredentials.User) <= 0 {
		return nil, errorNoUserDefined
	}

	if len(loginCredentials.Password) <= 0 {
		return nil, errorNoPasswordDefined
	}

	client := new(http.Client)

	loginBuffer := new(bytes.Buffer)
	jsonEncoder := json.NewEncoder(loginBuffer)
	if err := jsonEncoder.Encode(loginCredentials); err != nil {
		return nil, fmt.Errorf("Can not encode JSON from LoginCredential struct: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/users/login/", dockerHubAPI), loginBuffer)
	if err != nil {
		return nil, fmt.Errorf("Can not create request to get token from %v: %v", dockerHubAPI, err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("An error has occured after sending the http request to get a JWT token from %v: %v", dockerHubAPI, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid HTTP-Statuscode while getting the JWT Token: Get %v but expect 200", resp.StatusCode)
	}

	token := new(types.Token)
	jsonDecoder := json.NewDecoder(resp.Body)
	if err := jsonDecoder.Decode(token); err != nil {
		return nil, fmt.Errorf("Can not decode token: %v", err)
	}

	return token, nil
}

func PatchRepository(repository *types.Repository, token *types.Token) (*types.Repository, error) {

	if len(repository.Namespcace) <= 0 {
		return nil, errorNoNamespaceDefined
	}

	if len(repository.Name) <= 0 {
		return nil, errorNoRepositoryDefined
	}

	// repositoryBuffer := new(bytes.Buffer)
	// jsonEncoder := json.NewEncoder(repositoryBuffer)
	// if err := jsonEncoder.Encode(repository); err != nil {
	// 	return nil, fmt.Errorf("Can not encode JSON from Repository struct: %v", err)
	// }

	client := new(http.Client)

	patchURL, err := url.Parse(fmt.Sprintf("%v/repositories/%v/%v", dockerHubAPI, repository.Namespcace, repository.Name))
	if err != nil {
		return nil, fmt.Errorf("Can not prase URL: %v", err)
	}

	data := url.Values{}
	data.Set("full_description", repository.FullDescription)

	req, err := http.NewRequest(http.MethodPatch, patchURL.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Can not create http request to update file: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("JWT %v", token.Token))
	//req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("An error has occured: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Invalid HTTP-Statuscode: Get %v but expect 200: %v", resp.StatusCode, string(bodyBytes))
	}

	patchedRepository := new(types.Repository)

	jsonDecoder := json.NewDecoder(resp.Body)
	if err := jsonDecoder.Decode(patchedRepository); err != nil {
		return nil, fmt.Errorf("Can not encode JSON from Repository struct: %v", err)
	}

	return patchedRepository, nil
}
