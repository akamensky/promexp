package main

import (
	"fmt"
	"github.com/akamensky/prommetric"
	"log"
	"net/http"
	"time"
)

func main() {
	prommetric.DefaultExpiration = time.Second * 10
	prommetric.DefaultInterval = time.Second * 1
	e := prommetric.NewExporter()
	e.SetGauge("one_off_should_disappear", 1, "a metric that should disappear soon", map[string]string{"lbl": "val"})
	go func() {
		for true {
			e.SetGauge("time_since_epoch", float64(time.Now().Unix()), "seconds since Unix epoch", nil)
			time.Sleep(1 * time.Second)
		}
	}()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, e.Render())
		if err != nil {
			panic(err)
		}
	})

	listenAddr := ":8080"
	log.Println("listening on", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Panic(err)
	}
}
