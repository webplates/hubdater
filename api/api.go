package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://hub.docker.com/v2/"
const userAgent = "hubdater"

type Client struct {
	httpClient    *http.Client
	authenticator Authenticator

	BaseURL *url.URL

	UserAgent string

	common service

	Repositories *RepositoriesService
}

type service struct {
	client *Client
}

type Authenticator interface {
	Authenticate(req *http.Request) error
}

type JwtAuthenticator struct {
	Token string
}

type ApiError struct {
	Caller string
	Code   int
	Detail string `json:"detail"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("Cannot %s [%d]: %s", e.Caller, e.Code, e.Detail)
}

func NewClient(httpClient *http.Client, authenticator Authenticator) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	client := &Client{
		httpClient: httpClient,
		authenticator: authenticator,
		BaseURL:    baseURL,
		UserAgent:  userAgent,
	}

	client.common.client = client

	client.Repositories = (*RepositoriesService)(&client.common)

	return client
}

func (c *Client) NewRequest(method, URL string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, result interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func (c *Client) Login(username string, password string) (string, error) {
	body := map[string]string{
		"username": username,
		"password": password,
	}

	req, err := c.NewRequest("POST", "users/login", body)
	if err != nil {
		return "", err
	}

	b := make(map[string]string)

	resp, err := c.Do(req, &b)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", c.ParseError(fmt.Sprintf("login as %s", username), resp)
	}

	return b["token"], nil
}

func (c *Client) Authenticate(req *http.Request) error {
	if c.authenticator == nil {
		return errors.New("You need to login first to use this service")
	}

	c.authenticator.Authenticate(req)

	return nil
}

func (a *JwtAuthenticator) Authenticate(req *http.Request) error {
	req.Header.Add("Authorization", fmt.Sprintf("JWT %s", a.Token))

	return nil
}

func (c *Client) ParseError(caller string, resp *http.Response) error {
	decoder := json.NewDecoder(resp.Body)

	err := &ApiError{
		Caller: caller,
		Code:   resp.StatusCode,
	}

	decoder.Decode(&err)

	return err
}
