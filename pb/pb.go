package pb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginBody struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type PBEntry struct {
	Id      string `json:"id"`
	Updated string `json:"updated"`
}

type ListSearch struct {
	Items []PBEntry `json:"items"`
}

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
	fmt.Println(apiUrl)

	req, err := http.NewRequest("GET", apiUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pb.Token))
	if err != nil {
		fmt.Printf("Failed to create request: %s", err)
		return ListSearch{}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to retrieve last updated: %s", err)
		return ListSearch{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to retrieve last updated: %s", resp.Status)
		return ListSearch{}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %s", err)
		return ListSearch{}
	}

	var response ListSearch
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal response: %s", err)
		return ListSearch{}
	}

	return response
}
