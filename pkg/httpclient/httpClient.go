package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/cleoGitHub/golem/pkg/merror"
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

func (client HttpClient) Do(request *http.Request) ([]byte, int, error) {
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.Token))

	c := http.Client{}
	var resp *http.Response
	var err error
	var unauthorizedOnce bool = false
retryLoop:
	for i := 0; i < client.Config.NbRetry; i++ {
		resp, err = c.Do(request)
		if err != nil {
			// Use retry mecanism before returning an error
			continue
		}
		switch resp.StatusCode {
		case 200:
			break retryLoop
		case http.StatusForbidden, http.StatusUnauthorized:
			if unauthorizedOnce {
				return nil, resp.StatusCode, nil
			}
			unauthorizedOnce = true
			// call was unauthorized, try to authenticate before retrying
			if err = client.Authenticate(); err != nil {
				return nil, 0, merror.Stack(err)
			}
			i--
		default:
			return nil, resp.StatusCode, ErrUnexpectedStatus
		}
	}
	if err != nil {
		return nil, 0, merror.Stack(err)
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, merror.Stack(err)
	}

	return response, 0, nil
}
