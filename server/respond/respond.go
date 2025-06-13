package respond

import (
	"encoding/json"
	"net/http"
)

type ResponseStruct struct {
	Err bool   `json:"err"`
	Msg string `json:"msg"`
}

func RespondErrMsg(w http.ResponseWriter, msg string) {
	res := ResponseStruct{Msg: msg, Err: true}
	json.NewEncoder(w).Encode(res)
}

func RespondSuccess(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(ResponseStruct{Err: false, Msg: "success"})
}

var SUCCESS = ResponseStruct{
	Err: false, Msg: "success",
}

