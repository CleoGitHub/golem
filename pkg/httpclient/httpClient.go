package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cleoGitHub/golem/pkg/stringtool"
)

type HttpClient struct {
	SecurityClient SecurityClient
	Config         HttpClientConfig
	Token          string
	RefreshToken   string
}

func (client HttpClient) NewRequest(action string, endpoint string, body []byte, headers map[string]string) (*http.Request, error) {

	request, err := http.NewRequest(action, stringtool.RemoveDuplicate(client.Config.Host+":"+client.Config.Port+"/"+endpoint, '/'), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	return request, nil
}

func (client *HttpClient) Authenticate() error {
	if client.Token == "" || client.RefreshToken == "" {
		token, refreshToken, err := client.SecurityClient.Authenticate(client.Config.SA, client.Config.Password)
		if err != nil {
			return err
		}
		client.Token = token
		client.RefreshToken = refreshToken
	} else {
		token, refreshToken, err := client.SecurityClient.RefreshToken(client.RefreshToken)
		if err != nil {
			return err
		}
		client.Token = token
		client.RefreshToken = refreshToken
	}
	return nil
}

func (client HttpClient) Do(request *http.Request) ([]byte, error) {
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.Token))

	c := http.Client{}
	var resp *http.Response
	var err error
	for i := 0; i < client.Config.NbRetry; i++ {
		resp, err = c.Do(request)
		if err == nil && resp.StatusCode == 200 {
			break
		}
		if err != nil {
			fmt.Printf("%+v Error: %s\n", time.Now().Format(time.TimeOnly), err)
		}
		if resp != nil {
			defer resp.Body.Close()

			response, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			fmt.Printf("statusCode: %d, message: %s\n", resp.StatusCode, response)
			if resp.StatusCode == http.StatusNotFound {
				fmt.Printf("Url not foud: %v\n", request.URL)
			}
		}
		fmt.Printf("\n")
		if i < client.Config.NbRetry-1 {
			time.Sleep(time.Second * time.Duration(client.Config.IntervalRetry))
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d\nurl: %s\nresponse: %s", resp.StatusCode, request.URL.RawPath, string(response))
	}

	return response, nil
}

func (client HttpClient) DoPost(endpoint string) ([]byte, error) {
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return client.Do(request)
}
