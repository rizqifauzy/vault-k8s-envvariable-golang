package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	message := os.Getenv("MESSAGE")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	if message == "" {
		message = "Hello from Golang!"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n\n", message)

		// Menampilkan daftar pengguna MySQL di web
		db, err := sql.Open("mysql", dbURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to connect to database: %v", err), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to ping database: %v", err), http.StatusInternalServerError)
			return
		}

		rows, err := db.Query("SELECT User, Host FROM mysql.user")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to execute query: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var user, host string
		fmt.Fprintf(w, "List of MySQL Users :\n\n")
		for rows.Next() {
			if err := rows.Scan(&user, &host); err != nil {
				http.Error(w, fmt.Sprintf("Failed to scan row: %v", err), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "User: %s, Host: %s\n", user, host)
		}
	})

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}