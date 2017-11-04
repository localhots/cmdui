package github

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/juju/errors"

	"github.com/localhots/cmdui/backend/config"
)

type User struct {
	ID      uint   `json:"id"`
	Login   string `json:"login"`
	Name    string `json:"name"`
	Picture string `json:"avatar_url"`
}

const (
	authorizeURL   = "https://github.com/login/oauth/authorize"
	accessTokenURL = "https://github.com/login/oauth/access_token"
	userDetailsURL = "https://api.github.com/user"
)

func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	urlStr := authorizeURL + "?" + url.Values{
		"client_id": {config.Get().Github.ClientID},
		"scope":     {"read:user"},
	}.Encode()
	http.Redirect(w, r, urlStr, http.StatusTemporaryRedirect)
}

func ExchangeCode(ctx context.Context, code string) (accessToken string, err error) {
	cfg := config.Get()
	reqBody := bytes.NewBufferString(url.Values{
		"client_id":     {cfg.Github.ClientID},
		"client_secret": {cfg.Github.ClientSecret},
		"code":          {code},
	}.Encode())
	req, err := http.NewRequest(http.MethodPost, accessTokenURL, reqBody)
	if err != nil {
		return "", errors.Annotate(err, "Failed to create a code exchange request")
	}

	// Passing client request context
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Annotate(err, "Failed to perform code exchange request")
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Annotate(err, "Failed to read access token response")
	}

	uri, err := url.ParseQuery(string(respBody))
	if err != nil {
		return "", errors.Annotate(err, "Failed to parse access token response")
	}

	return uri.Get("access_token"), nil
}

func AuthDetails(accessToken string) (User, error) {
	var u User
	reqURL := userDetailsURL + "?" + url.Values{
		"access_token": {accessToken},
	}.Encode()
	resp, err := http.Get(reqURL)
	if err != nil {
		return u, errors.Annotate(err, "Failed to fetch authenticated user details")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return u, errors.Annotate(err, "Failed to read authenticated user details")
	}
	if err := json.Unmarshal(body, &u); err != nil {
		return u, errors.Annotate(err, "Failed to parse authenticated user details")
	}

	return u, nil
}
