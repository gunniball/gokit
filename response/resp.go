package response

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"github.com/pquerna/ffjson/ffjson"
	"net/http"
	"sync"
)

type respError struct {
	Error string `json:"Error"`
}

var respErrorPool = sync.Pool{New: func() interface{} {
	return new(respError)
}}

func ErrorInternalServer(err string, rw web.ResponseWriter) {
	rw.Header().Set("Content-type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusInternalServerError)
	WriteError(err, rw)
}

func ErrorBadRequest(err string, rw web.ResponseWriter) {
	rw.Header().Set("Content-type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusBadRequest)
	WriteError(err, rw)
}

func ErrorNotAuthorized(err string, rw web.ResponseWriter) {
	rw.Header().Set("Content-type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusUnauthorized)
	WriteError(err, rw)
}

func ErrorForbidden(err string, rw web.ResponseWriter) {
	rw.Header().Set("Content-type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusForbidden)
	WriteError(err, rw)
}

func Error(err string, rw web.ResponseWriter) {
	rw.Header().Set("Content-type", "application/json; charset=utf-8")
	WriteError(err, rw)
}

func WriteError(err string, rw web.ResponseWriter) {
	data := respErrorPool.Get().(*respError)
	data.Error = err
	jsn, _ := ffjson.Marshal(data)
	respErrorPool.Put(data)
	rw.Write(jsn)
	ffjson.Pool(jsn)
}

func JsonIntent(data interface{}, rw web.ResponseWriter) {
	jsn, _ := json.MarshalIndent(&data, " ", " ")
	rw.Header().Set("Content-type", "application/json; charset=utf-8")
	rw.Write(jsn)
}

func Json(data interface{}, rw web.ResponseWriter) {
	jsn, _ := ffjson.Marshal(&data)
	rw.Header().Set("Content-type", "application/json; charset=utf-8")
	rw.Write(jsn)
	ffjson.Pool(jsn)
}

func Html(data string, rw web.ResponseWriter) {
	rw.Header().Set("Content-type", "text/html; charset=utf-8")
	fmt.Fprintf(rw, "%s", data)
}

func String(data string, rw web.ResponseWriter) {
	Plain([]byte(data), rw)
}

func Plain(data []byte, rw web.ResponseWriter) {
	rw.Header().Set("Content-type", "text/plain; charset=utf-8")
	fmt.Fprintf(rw, "%s", data)
}
