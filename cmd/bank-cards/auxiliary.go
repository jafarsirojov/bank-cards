package main

import (
	"github.com/jafarsirojov/bank-cards/pkg/core/auth"
	"fmt"
	"github.com/jafarsirojov/mux/pkg/mux/middleware/jwt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

var jwtMiddleware = jwt.JWT(reflect.TypeOf((*auth.Auth)(nil)).Elem(), []byte("top secret"))

func requestIdier(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		request.Header.Set("X-Id", fmt.Sprintf("%d", time.Now().Unix()))
		log.Printf("header set: %s", request.Header.Get("X-Id"))
		handler(writer, request)
	}
}

func logger(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Print("before")
		handler(writer, request)
		log.Print("after")
	}
}

func FromFlagOrEnv(flag string, env string) (value string, ok bool) {
	if flag != "" {
		return flag, true
	}

	return os.LookupEnv(env)
}
