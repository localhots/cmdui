package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/localhots/cmdui/backend/api/assets"
	"github.com/localhots/cmdui/backend/api/auth"
	"github.com/localhots/cmdui/backend/config"
	"github.com/localhots/cmdui/backend/db"
	"github.com/localhots/cmdui/backend/log"
)

// Start starts a web server that runs both backend API and serves assets that
// support the UI.
func Start() error {
	assHand := assets.Handler()
	router.NotFound = func(w http.ResponseWriter, r *http.Request) {
		assHand.ServeHTTP(w, r)
	}

	cfg := config.Get().Server
	log.Logger().Infof("Starting command UI server at %s:%d", cfg.Host, cfg.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
}

//
// Endpoints
//

type handle func(ctx context.Context, w http.ResponseWriter, r *http.Request)

const rootPath = "/"

var router = httprouter.New()

func openEndpoint(h handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := contextWithParams(r.Context(), params)
		h(ctx, w, r)
	}
}

func protectedEndpoint(h handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx, err := auth.AuthenticateRequest(w, r)
		if err != nil {
			renderUnauthorized(w, err)
			return
		}

		ctx = contextWithParams(ctx, params)
		h(ctx, w, r)
	}
}

//
// Rendering
//

// renderJSON is a convinience function that encodes any value as JSON and
// writes it to response with appropriate headers included.
func renderJSON(w http.ResponseWriter, v interface{}) {
	body, err := json.Marshal(v)
	if err != nil {
		renderError(w, err, http.StatusInternalServerError, "Failed to encode response into JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func renderError(w http.ResponseWriter, err error, status int, msg string) {
	log.WithFields(log.F{
		"status": status,
		"error":  err,
	}).Warnf("Request failed: %s", msg)
	http.Error(w, msg, status)
}

func renderUnauthorized(w http.ResponseWriter, err error) {
	renderError(w, err, http.StatusUnauthorized, "Unauthorized")
}

//
// Params and context
//

type ctxKey string

const (
	ctxParamsKey ctxKey = "params"
)

func contextWithParams(ctx context.Context, params httprouter.Params) context.Context {
	return context.WithValue(ctx, ctxParamsKey, params)
}

func paramsFromContext(ctx context.Context) (params httprouter.Params, ok bool) {
	v := ctx.Value(ctxParamsKey)
	if v == nil {
		return nil, false
	}

	return v.(httprouter.Params), true
}

func param(ctx context.Context, name string) string {
	if params, ok := paramsFromContext(ctx); ok {
		return params.ByName(name)
	}
	return ""
}

func requestedPage(r *http.Request) db.Page {
	offset := r.FormValue("offset")
	limit := r.FormValue("limit")
	if offset == "" && limit == "" {
		return db.Page{}
	}

	str2uint := func(s string) uint {
		u, _ := strconv.ParseUint(s, 10, 64)
		return uint(u)
	}
	return db.Page{
		Offset: str2uint(offset),
		Limit:  str2uint(limit),
	}
}

//
// Utils
//

// unbufferedWriter is an implementation of http.ResponseWriter that flushes the
// buffer after every write.
type unbufferedWriter struct {
	http.ResponseWriter
}

func (w unbufferedWriter) Write(p []byte) (int, error) {
	n, err := w.ResponseWriter.Write(p)
	if f, ok := w.ResponseWriter.(http.Flusher); ok && err == nil {
		f.Flush()
	}
	return n, err
}
