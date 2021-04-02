package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"wilbertopachecob/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	infolog  *log.Logger
	errorlog *log.Logger
	snippets *mysql.SnippetModel
}

func getEnvVar(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Importantly, we use the flag.Parse() function to parse the command-line
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any er
	// encountered during parsing the application will be terminated.
	flag.Parse()

	//using a file to store info logs
	f, err := os.OpenFile("./tpm/info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dns := fmt.Sprintf("%s:%s@/%s?parseTime=true", getEnvVar("DB_USERNAME"), getEnvVar("DB_PASSWORD"), getEnvVar("DB_DATABASE"))
	db, err := openDB(dns)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		infolog:  infoLog,
		errorlog: errorLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields
	// that the server uses the same network address and routes as before, and
	// the ErrorLog field so that the server now uses the custom errorLog logge
	// the event of any problems.
	svr := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	//log.Printf("Starting server on port %s", getEnvVar("PORT"))
	infoLog.Printf("Starting server on port %s", *addr)
	err = svr.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
