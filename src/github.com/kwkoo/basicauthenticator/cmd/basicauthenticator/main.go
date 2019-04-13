package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kwkoo/basicauthenticator"
	"github.com/kwkoo/configparser"
)

var realm string
var db basicauthenticator.UserDB

func unauthorized(w http.ResponseWriter, message string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", realm))

	e := struct {
		Error string `json:"error"`
	}{message}

	var s string
	b, err := json.Marshal(e)
	if err == nil {
		s = string(b)
	} else {
		s = message
	}
	http.Error(w, s, http.StatusUnauthorized)
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Print("request for URI: ", path)

	w.Header().Set("Content-Type", "application/json")

	authheader := r.Header.Get("Authorization")
	if len(authheader) == 0 {
		log.Print("no authorization header")
		unauthorized(w, "no authorization header")
		return
	}

	if !strings.HasPrefix(authheader, "Basic ") {
		log.Printf("authorization header does not contain Basic keyword: %s", authheader)
		unauthorized(w, "authorization header does not contain Basic keyword")
		return
	}

	authheader = authheader[len("Basic "):]
	data, err := base64.StdEncoding.DecodeString(authheader)
	if err != nil {
		log.Printf("could not decode Base64 string %s: %v", authheader, err)
		unauthorized(w, "could not decode Base64 string")
		return
	}

	s := string(data)
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		log.Printf("authentication string must be of the form userid:password - %s", s)
		unauthorized(w, "authentication string must be of the form userid:password")
		return
	}

	u := db.Authenticate(parts[0], parts[1])
	if u == nil {
		log.Printf("%s / %s could not be authenticated", parts[0], parts[1])
		unauthorized(w, "user could not be authenticated")
		return
	}

	b, err := json.Marshal(u)
	if err != nil {
		log.Printf("could not marshall user struct: %v", err)
		fmt.Fprintln(w, "{\"error\":\"could not marshall user struct\"}")
		return
	}

	fmt.Fprintf(w, string(b))
}

func main() {
	config := struct {
		Realm  string `default:"default" usage:"authentication realm"`
		Port   int    `default:"8080" usage:"HTTP listener port"`
		Userdb string `usage:"user database file - this is a text file with an entry for a user on each line, with each line containing tab-separated values of user ID, password, name, email address" mandatory:"true"`
	}{}

	if err := configparser.Parse(&config); err != nil {
		log.Fatal(err)
	}

	initializeUserDB(config.Userdb)
	log.Printf("user database contains %d entries", db.Size())

	realm = config.Realm

	log.Printf("listening on port %v", config.Port)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}

func initializeUserDB(dbfilename string) {
	f, err := os.Open(dbfilename)
	if err != nil {
		log.Fatalf("error while trying to open user database %s: %v", dbfilename, err)
	}
	defer f.Close()

	db = basicauthenticator.NewUserDB(f)
}
