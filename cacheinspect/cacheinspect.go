package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/user"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/img", imageHandler)
	log.Println("Server running on :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func getSqlitePath() string {
	var basepath string
	usr, err := user.Current()
	checkErr(err)
	switch runtime.GOOS {
	case "darwin":
		basepath = filepath.Join(usr.HomeDir, "Library", "Application Support", "Steam")
	case "linux":
		basepath = filepath.Join(usr.HomeDir, ".local", "share", "Steam")
	case "windows":
		if runtime.GOARCH == "amd64" {
			basepath = filepath.Join("C:", "Program Files (x86)", "Steam")
		} else if runtime.GOARCH == "386" {
			basepath = filepath.Join("C:", "Program Files", "Steam")
		} else {
			panic(fmt.Sprintf("Unable to handle architecture %v", runtime.GOARCH))
		}
	}
	return filepath.Join(basepath, "SteamApps", "common", "Rust", "cache", "Storage.db")
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	entity := r.URL.Query().Get("entity")
	crc := r.URL.Query().Get("crc")

	if len(entity) == 0 && len(crc) == 0 {
		http.NotFound(w, r)
		return
	}

	query := "SELECT data FROM data where entity = ? and crc = ?"
	db, err := sql.Open("sqlite3", getSqlitePath())
	defer db.Close()
	var data []byte
	err = db.QueryRow(query, entity, crc).Scan(&data)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
	case err != nil:
		checkErr(err)
	}
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(data)
	checkErr(err)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	last := r.URL.Query().Get("last")
	query := "SELECT crc, entity, num, lastaccess FROM data"
	if len(last) != 0 {
		query += " where lastaccess > ?"
	}
	query += " order by lastaccess desc"

	db, err := sql.Open("sqlite3", getSqlitePath())
	checkErr(err)
	defer db.Close()

	rows, err := db.Query(query, last)
	checkErr(err)
	defer rows.Close()

	type Data struct {
		Entity     int64
		Crc        int64
		Num        int64
		LastAccess int64
	}

	var dataResults []Data

	for rows.Next() {
		var crc int64
		var entity int64
		var num int64
		var lastaccess int64
		err = rows.Scan(&crc, &entity, &num, &lastaccess)
		checkErr(err)
		d := Data{Entity: entity, Crc: crc, Num: num, LastAccess: lastaccess}
		dataResults = append(dataResults, d)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataResults)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
