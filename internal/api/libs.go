package api

import (
	"encoding/json"
	"net/http"
)

func writeOkResponse(w http.ResponseWriter, okCode int, body []byte) {
	if okCode > 299 || okCode < 200 {
		okCode = 200
	}
	w.Header().Add("Content-Type", JSON_CONTENT_TYPE_HEADER_VALUE)
	w.WriteHeader(okCode)
	w.Write(body)
}

func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	type errorResp struct {
		ErrorMsg string `json:"error"`
	}
	w.WriteHeader(errorCode)
	errorRespJson := errorResp{ErrorMsg: errorMsg}
	errorRespJsonData, err := json.Marshal(errorRespJson)
	if err == nil {
		w.Write(errorRespJsonData)
		return
	}
	w.Write([]byte("something went wrong"))
}
