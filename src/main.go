package main

import (
	"io"
	"log"
	"net/http"
	"shopstic/prom-labels-injector/config"
	"shopstic/prom-labels-injector/label_injector"
	"shopstic/prom-labels-injector/util"
)

func main() {
	settings := config.LoadFromEnvVars()

	log.Printf("Configuration: %+v", settings)

	handleMetrics := func(w http.ResponseWriter, r *http.Request) {
		metricsHandler(&settings, w, r)
	}

	http.HandleFunc("/metrics", handleMetrics)
	log.Fatal(http.ListenAndServe(":"+util.Uint16ToString(settings.Server.Port), nil))
}

func metricsHandler(settings *config.Config, w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(settings.PrometheusTarget.Address)
	if err != nil {
		log.Printf("Got error [%v] while making request to target [%v]", err, settings.PrometheusTarget.Address)
		w.WriteHeader(500)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Got error while [%v] while reading response body from target [%v]", err, settings.PrometheusTarget.Address)
		w.WriteHeader(500)
		return
	}
	metrics := string(body)
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(label_injector.InjectLabels(metrics, &settings.LabelInjector)))
}
