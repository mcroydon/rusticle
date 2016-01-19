package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// TODO: Get path from cli/env
	steamPath, err := FindSteam()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found Steam: %s", steamPath)
	rusticle, err := NewRustHandler(steamPath)
	defer rusticle.Close()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/data", rusticle.dataHandler)
	http.HandleFunc("/img", rusticle.imageHandler)

	log.Println("Server running on :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func getSqlitePath(steamPath string) (string, error) {
	p := filepath.Join(steamPath, "SteamApps", "common", "Rust", "cache", "Storage.db")
	_, err := os.Stat(p)
	return p, err
}

func NewRustHandler(steamPath string) (*rustHandler, error) {
	dbPath, err := getSqlitePath(steamPath)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &rustHandler{
		db: db,
	}, nil
}

type rustHandler struct {
	db *sql.DB
}

func (handler *rustHandler) Close() error {
	return handler.db.Close()
}

func (handler *rustHandler) imageHandler(w http.ResponseWriter, r *http.Request) {
	entity := r.URL.Query().Get("entity")
	crc := r.URL.Query().Get("crc")

	if len(entity) == 0 && len(crc) == 0 {
		http.NotFound(w, r)
		return
	}

	// TODO: Probably should do this once when the server boots and exit if there's an error.
	var data []byte
	query := "SELECT data FROM data where entity = ? and crc = ?"
	err := handler.db.QueryRow(query, entity, crc).Scan(&data)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
	case err != nil:
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (handler *rustHandler) dataHandler(w http.ResponseWriter, r *http.Request) {
	last := r.URL.Query().Get("last")
	args := []interface{}{}
	query := "SELECT crc, entity, num, lastaccess FROM data"
	if len(last) != 0 {
		query += " where lastaccess > ?"
		args = append(args, last)
	}
	query += " order by lastaccess desc"

	rows, err := handler.db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
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
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		d := Data{Entity: entity, Crc: crc, Num: num, LastAccess: lastaccess}
		dataResults = append(dataResults, d)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataResults)
}
