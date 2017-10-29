package api

import (
	"context"
	"net/http"

	"github.com/juju/errors"

	"github.com/localhots/cmdui/backend/api/auth"
	"github.com/localhots/cmdui/backend/api/github"
	"github.com/localhots/cmdui/backend/db"
)

const (
	callbackURL = "/api/auth/callback"
)

func init() {
	router.GET("/api/auth/login", openEndpoint(authLoginHandler))
	router.POST("/api/auth/logout", protectedEndpoint(authLogoutHandler))
	router.GET("/api/auth/session", protectedEndpoint(authSessionHandler))
	router.GET(callbackURL, openEndpoint(authCallbackHandler))
}

// authSessionHandler returns currently authenticated user details.
// GET /auth/session
func authSessionHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	renderJSON(w, auth.User(ctx))
}

// authLogoutHandler clears authentication cookies.
// GET /auth/session
func authLogoutHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	auth.ClearCookie(ctx, w)
	http.Redirect(w, r, rootPath, http.StatusTemporaryRedirect)
}

// authLoginHandler redirects user to a GitHub login page.
func authLoginHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	github.RedirectToLogin(w, r)
}

// authCallbackHandler accepts GitHub auth callback.
func authCallbackHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ghu, err := authWithGithubCode(ctx, r.FormValue("code"))
	if err != nil {
		renderError(w, err, http.StatusInternalServerError, "GitHub login failed")
		return
	}

	u, err := findOrCreateUser(ghu)
	if err != nil {
		renderError(w, err, http.StatusInternalServerError, "Failed to find a user using GitHub profile")
		return
	}

	sess := db.NewSession(u.ID)
	if err := sess.Create(); err != nil {
		renderError(w, err, http.StatusInternalServerError, "Failed to create a session")
		return
	}

	ctx = auth.ContextWithSession(ctx, sess)
	auth.AuthorizeResponse(ctx, w)
	auth.CacheSession(sess)

	http.Redirect(w, r, rootPath, http.StatusTemporaryRedirect)
}

func authWithGithubCode(ctx context.Context, code string) (github.User, error) {
	if code == "" {
		return github.User{}, errors.New("Missing authentication code")
	}

	accessToken, err := github.ExchangeCode(ctx, code)
	if err != nil {
		return github.User{}, errors.Annotate(err, "Failed to exchange code for access token")
	}

	ghu, err := github.AuthDetails(accessToken)
	if err != nil {
		return ghu, errors.Annotate(err, "Failed to fetch authenticated GitHub user details")
	}

	return ghu, nil
}

func findOrCreateUser(ghu github.User) (db.User, error) {
	u, err := db.FindUserByGithubID(ghu.ID)
	if err != nil {
		return db.User{}, errors.Annotate(err, "Failed to find GitHub user")
	}
	if u != nil {
		importGithubProfile(u, ghu)
		if err := u.Update(); err != nil {
			return *u, errors.Annotate(err, "Failed to update a user")
		}
	} else {
		eu := db.NewUser()
		u = &eu
		importGithubProfile(u, ghu)
		if err := u.Create(); err != nil {
			return *u, errors.Annotate(err, "Failed to create a user")
		}
	}

	return *u, nil
}

func importGithubProfile(u *db.User, ghu github.User) {
	u.GithubID = ghu.ID
	u.GithubLogin = ghu.Login
	u.Name = ghu.Name
	u.Picture = ghu.Picture
}
