package service 

import (
	"net/http"
	"fmt"
	"github.com/go-playground/log"
)

func logAccess(req *http.Request, respWriter http.ResponseWriter) (*http.Request, bool, error) {
	defer log.WithTrace().Info(fmt.Sprintf("[%s] %s", req.Method, req.URL))
	return req, true, nil
}