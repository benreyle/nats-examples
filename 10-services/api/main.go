package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"nats-examples/10-services/api/store"
	"nats-examples/10-services/helpers"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
)

var conn stan.Conn

func main() {
	conn = helpers.ConnectToSTAN("api")

	r := mux.NewRouter()
	r.Handle("/position", CreatePosition()).Methods("POST")
	r.Handle("/position", GetPosition()).Methods("GET")
	r.Handle("/position", DeletePosition()).Methods("DELETE")
	r.Handle("/positions", ListPositions()).Methods("GET")

	http.Handle("/", r)

	log.Println("Listening :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func CreatePosition() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.ResponseJSON(rw, "invalid position", http.StatusBadRequest)
			return
		}

		var position *helpers.Position

		err = json.Unmarshal(b, &position)
		if err != nil {
			helpers.ResponseJSON(rw, err, http.StatusBadRequest)
			return
		}

		position.ID = uuid.New()

		store.Save(position)

		body, err := json.Marshal(position.ID)
		if err != nil {
			helpers.ResponseJSON(rw, err, http.StatusInternalServerError)
			return
		}

		err = conn.Publish("create-position", body)
		if err != nil {
			helpers.ResponseJSON(rw, err, http.StatusInternalServerError)
			return
		}

		log.Println("Publishing position created event", position.ID.String())

		helpers.ResponseJSON(rw, position.ID, http.StatusCreated)
	})
}

func GetPosition() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		strID := r.URL.Query().Get("id")
		if strID == "" {
			helpers.ResponseJSON(rw, "invalid id", http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(strID)
		if err != nil {
			helpers.ResponseJSON(rw, err, http.StatusBadRequest)
			return
		}

		position, err := store.Get(id)
		if err != nil {
			helpers.ResponseJSON(rw, err, http.StatusNotFound)
			return
		}

		helpers.ResponseJSON(rw, position, http.StatusOK)
	})
}

func ListPositions() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		positions := store.List()
		helpers.ResponseJSON(rw, positions, http.StatusOK)
	})
}

func DeletePosition() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		strID := r.URL.Query().Get("id")
		if strID == "" {
			helpers.ResponseJSON(rw, "invalid id", http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(strID)
		if err != nil {
			helpers.ResponseJSON(rw, err, http.StatusBadRequest)
			return
		}

		err = store.Delete(id)
		if err != nil {
			helpers.ResponseJSON(rw, err, http.StatusNotFound)
			return
		}

		body, err := json.Marshal(id)
		if err != nil {
			helpers.ResponseJSON(rw, err, http.StatusInternalServerError)
			return
		}

		err = conn.Publish("delete-position", body)
		if err != nil {
			helpers.ResponseJSON(rw, err, http.StatusInternalServerError)
			return
		}

		log.Println("Publishing position deleted event", id.String())

		helpers.ResponseJSON(rw, nil, http.StatusNoContent)
	})
}
