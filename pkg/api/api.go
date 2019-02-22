package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pakerfeldt/lovi/pkg/core"
)

func Init(router *mux.Router) {
	router.HandleFunc("/alert/trigger/{policy}", getAlert).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getAlert(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["message"]) > 0 {
		params := mux.Vars(r)
		core.HandleAlert(params["policy"], r.URL.Query()["message"][0])
	} else {
		log.Println("Error triggering alert, no message.")
	}
}
