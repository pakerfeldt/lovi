package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pakerfeldt/lovi/pkg/core"
	"github.com/pakerfeldt/lovi/pkg/models"
)

func Init(router *mux.Router, settings models.Settings) {
	router.HandleFunc("/event/trigger/{policy}", getEvent).Methods("GET")
	log.Printf("Listening on 0.0.0.0:%s\n", strconv.Itoa(settings.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(settings.Port), router))
}

func getEvent(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["message"]) > 0 {
		params := mux.Vars(r)
		core.HandleEvent(params["policy"], r.URL.Query()["message"][0])
	} else {
		log.Println("Error triggering event, no message.")
	}
}
