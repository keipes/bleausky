package db

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
)

type BleauDB interface {
	Create()
}

func CreateDB() {
	// connect to sqlite3 database named "BleauDB"
	db, err := sql.Open("sqlite", "BleauDB")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// get db version and log
	var version string
	row := db.QueryRow("SELECT sqlite_version()")
	err = row.Scan(&version)
	if err != nil {
		panic(err)
	}
	db.Exec("CREATE TABLE IF NOT EXISTS bleau (id INTEGER PRIMARY KEY, name TEXT)")
	// table posts (repo string, cid string, text string, seq int64)
	db.Exec("CREATE TABLE IF NOT EXISTS posts (repo TEXT, cid TEXT, text TEXT, seq INTEGER)")
	//db.Exec("INSERT INTO bleau (name) VALUES ('BleauDB2')")
	println("SQLite version: ", version)
	// print all fields in bleau table
	rows, err := db.Query("SELECT * FROM bleau")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		println(id, name)
	}
}

func InsertPost(cid string, did string, seq int64) {
	db, err := sql.Open("sqlite", "BleauDB")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Exec("INSERT INTO posts (repo, cid, seq) VALUES (?, ?, ?)", cid, did, seq)
}

func ListPosts() {
	db, err := sql.Open("sqlite", "BleauDB")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT repo, cid, seq FROM posts")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var repo string
		var cid string
		var _ string
		var seq int64
		err = rows.Scan(&repo, &cid, &seq)
		if err != nil {
			panic(err)
		}
		println(repo, cid, seq)
	}
}

func CountPosts() {
	db, err := sql.Open("sqlite", "BleauDB")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT count(*) FROM posts")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		err = rows.Scan(&count)
		if err != nil {
			panic(err)
		}
		println(count)
	}
}

func DropTables() {
	db, err := sql.Open("sqlite", "BleauDB")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Exec("DROP TABLE posts")
	db.Exec("DROP TABLE bleau")
}
