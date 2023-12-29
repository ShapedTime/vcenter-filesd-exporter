package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"vcenter-prom-filesd-go/controller"
	"vcenter-prom-filesd-go/model"
	_ "vcenter-prom-filesd-go/promhelper"
	"vcenter-prom-filesd-go/vcenter-helper"
)

var (
	promModel      *model.PromModel
	promController *controller.PromController
	vc             *vcenterhelper.VCenterHelper
)

func main() {
	datacenter := os.Getenv("DATACENTER")

	host := os.Getenv("VC_HOST")
	user := os.Getenv("VC_USER")
	password := os.Getenv("VC_PASSWORD")

	port := os.Getenv("PORT")

	tlsEnabled := flag.Bool("tls", false, "enable TLS")
	flag.Parse()

	if port == "" {
		port = "8080"
		if *tlsEnabled {
			port = "443"
		}
	}

	if datacenter == "" || host == "" || user == "" || password == "" {
		panic("DATACENTER or VC_HOST or VC_USER or VC_PASSWORD env var is empty")
	}

	setupMC(host, user, password, datacenter)

	createServer(port, *tlsEnabled)
}

// function that creates a http server and listens on port 8080
func createServer(port string, tlsEnabled bool) {
	http.HandleFunc("/prom", promController.PromHandler)
	http.Handle("/metrics", promhttp.Handler())
	if tlsEnabled {
		err := http.ListenAndServeTLS(fmt.Sprintf(":%s", port), "vcenter-exporter.cert", "vcenter-exporter.key", nil)
		if err != nil {
			panic(err)
		}
	} else {
		err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
		if err != nil {
			panic(err)
		}
	}
}

func setupMC(host, user, password, datacenter string) {
	vc = vcenterhelper.NewVCenterHelper(host, user, password)

	promModel = model.NewPromModel(vc, datacenter)

	promController = controller.NewPromController(promModel)
}
