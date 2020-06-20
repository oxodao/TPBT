package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"tpbt/bterrors"
	"tpbt/dal"
	"tpbt/helix"
	"tpbt/models"
	"tpbt/oauth"
	"tpbt/services"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

func RunServer(prv *services.Provider) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	r := mux.NewRouter()

	r.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Access-Control-Allow-Origin", "*")
		url := oauth.Configuration.AuthCodeURL("test")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	r.HandleFunc("/auth/callback/{code}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		code := mux.Vars(r)["code"]

		if len(code) == 0 {
			fmt.Println("No code given!")
			return
		}

		tk, err := oauth.Configuration.Exchange(context.Background(), code)
		if err != nil {
			fmt.Println(err)
			return
		}

		user := &models.BTStreamer{
			AccessToken: tk.AccessToken,
			RefreshToken: tk.RefreshToken,
			MutexWS: &sync.Mutex{},
		}
		err = helix.FetchUser(prv, user, true)
		if bterrors.WriteError(w, err) {
			return
		}

		err = dal.FindOrCreateUser(prv, user)
		if err != nil {
			fmt.Println("Can't get user!")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		toSend, err := json.Marshal(user)

		w.WriteHeader(http.StatusOK)
		w.Write(toSend)
	})

	r.HandleFunc("/auth/verify/{app_token}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		token := mux.Vars(r)["app_token"]

		user, err := dal.FetchUserFromToken(prv, token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)

		toSend, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(toSend)
	})

	r.HandleFunc("/connect/{app_token}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		token := mux.Vars(r)["app_token"]

		user, err := dal.FetchUserFromToken(prv, token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = helix.FetchUser(prv, user, true)
		if bterrors.WriteError(w, err) {
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//defer c.Close()

		user.Websocket = c

		_ = SendCommand(user, "SET_USER", user)

		scores, err := dal.FetchScoreboard(prv, user)
		if err != nil {
			fmt.Println("Something went wrong fetching scores")
		} else {
			_ = SendCommand(user, "SET_LEADERBOARD", scores)
		}

		ListenWS(prv, user)
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%v", prv.Config.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
