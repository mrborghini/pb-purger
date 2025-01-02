package pb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PB struct {
	Url               string
	Token             string
	Username          string
	Password          string
	AccountCollection string
}

func NewPB(url string, username string, password string, accountCollection string) *PB {
	return &PB{
		Url:               url,
		Token:             "",
		Username:          username,
		Password:          password,
		AccountCollection: accountCollection,
	}
}

func (pb *PB) Login() error {
	apiUrl := fmt.Sprintf("%s/api/collections/%s/auth-with-password", pb.Url, pb.AccountCollection)
	data, err := json.Marshal(LoginBody{
		Identity: pb.Username,
		Password: pb.Password,
	})
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(data)
	req, err := http.NewRequest("POST", apiUrl, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var logincreds LoginResponse
	json.Unmarshal(bodyBytes, &logincreds)

	pb.Token = logincreds.Token
	return nil
}

func (pb *PB) RetrieveLastUpdated(collection string) ListSearch {
	apiUrl := fmt.Sprintf("%s/api/collections/%s/records?perPage=10000", pb.Url, collection)

	req, err := http.NewRequest("GET", apiUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pb.Token))
	if err != nil {
		return ListSearch{}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ListSearch{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ListSearch{}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ListSearch{}
	}

	var response ListSearch
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return ListSearch{}
	}

	return response
}

func (pb *PB) Delete(id string, collection string) bool {
	apiUrl := fmt.Sprintf("%s/api/collections/%s/records/%s", pb.Url, collection, id)
	req, err := http.NewRequest("DELETE", apiUrl, nil)
	if err != nil {
		return false
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pb.Token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent
}
