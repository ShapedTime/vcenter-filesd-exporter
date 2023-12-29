package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"vcenter-prom-filesd-go/model"
	"vcenter-prom-filesd-go/promhelper"
)

type PromController struct {
	promModel *model.PromModel
}

func NewPromController(promModel *model.PromModel) *PromController {
	return &PromController{promModel: promModel}
}

func (c PromController) PromHandler(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	promhelper.RequestReceived(path)

	vms, err := c.promModel.GetVMs(path)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte("cannot get vms"))
		promhelper.Error()
		log.Println(err)
		return
	}

	resp, err := json.Marshal(vms)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte("cannot marshal json"))
		promhelper.Error()
		log.Println(err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(resp)
	if err != nil {
		promhelper.Error()
		log.Println(err)
		return
	}
}
