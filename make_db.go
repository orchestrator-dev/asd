package main
import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)
func main() {
	db, _ := sql.Open("sqlite3", "test.db")
	db.Exec("CREATE TABLE users(id INTEGER PRIMARY KEY, name TEXT);")
	db.Exec("INSERT INTO users (name) VALUES ('Alice'), ('Bob');")
	db.Close()
}
