package web

import (
	"net/http"

	"github.com/gorilla/mux"

	"encoding/json"

	log "github.com/morriswinkler/cloudglog"
)

type Result struct {
	Success bool   `json:"sucess"`
	Message string `json:"message"`
}

type JsonResult struct {
	Result `json:"result"`
}

func ServiceStartRecipe(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	id := vars["id"]

	log.V(3).Infoln("[API] startRecipe: ", id)

	jsonMesage := &JsonResult{
		Result{
			true,
			"",
		},
	}

	jsonString, err := json.Marshal(&jsonMesage)
	if err != nil {
		log.Errorln("[API] error: ", err)
	}

	setHeaders(res)

	res.Write(jsonString)
}

func setHeaders(res http.ResponseWriter) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type,Accept")
}
