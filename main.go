package main

import (
	"os"
	"net/http"
)

var avatar_dir = "avatars"

func status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"response\":200,\"status\":0}"))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var avatarid = "0"
	if _, err := os.Stat(avatar_dir + "/" + r.URL.Path[1:] + ".png"); !os.IsNotExist(err) {
		avatarid = r.URL.Path[1:]
	}
	http.ServeFile(w, r, avatar_dir + "/" + avatarid + ".png")
}

func main() {
	if _, err := os.Stat(avatar_dir); os.IsNotExist(err) {
		os.Mkdir(avatar_dir, 0777)
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/status", status)
	http.ListenAndServe(":4999", nil)
}
