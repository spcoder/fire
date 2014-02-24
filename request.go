package fire

import (
	"net/http"
	"strconv"
)

type Request struct {
	HttpRequest *http.Request
}

func (r Request) QueryInt(key string, def int) int {
	value := r.HttpRequest.URL.Query().Get(key)
	if value == "" {
		return def
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return def
	}

	return i
}

func (r Request) QueryString(key string, def string) string {
	value := r.HttpRequest.URL.Query().Get(key)
	if value == "" {
		return def
	}
	return value
}
