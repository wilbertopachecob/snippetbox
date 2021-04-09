package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"wilbertopachecob/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func getEnvVar(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func showHelp() {
	fmt.Print(`
	Usage Snippets CLI
	Options:
		-show	        Shows [snippets]			(usage="snippets")
		get	[options]	Get a user or snippet		(usage=-model "user" -id 1)
			-model      The model name				(usage="user|snippet")
			-id         The id of the model			(usage=1)
	`)
}

func setFlag(flag *flag.FlagSet) {
	flag.Usage = func() {
		showHelp()
	}
}

func main() {
	var strF string
	flag.StringVar(&strF, "show", "", "")

	getCMD := flag.NewFlagSet("get", flag.ExitOnError)
	model := getCMD.String("model", "", "string representing the model")
	id := getCMD.Int("id", 0, "an integer representing the model id")

	setFlag(flag.CommandLine)
	flag.Parse()

	dns := fmt.Sprintf("%s:%s@/%s?parseTime=true", getEnvVar("DB_USERNAME"), getEnvVar("DB_PASSWORD"), getEnvVar("DB_DATABASE"))
	db, err := openDB(dns)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if strF != "" {
		switch strF {
		case "snippets":
			sm := &mysql.SnippetModel{DB: db}
			ss, err := sm.Latest()
			if err != nil {
				log.Fatal(err)
			}
			for index, s := range ss {
				snippet, err := json.MarshalIndent(s, "", "   ")
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Snippet %d: %+v \n", index, string(snippet))
			}
		case "users":
			fmt.Print("Users selected")
		default:
			showHelp()
		}
	}
	switch os.Args[1] {
	case "get":
		getCMD.Parse(os.Args[2:])
		if *model != "" && *id > 0 {
			err = getModel(db, *model, *id)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			showHelp()
		}
	}
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

func getModel(db *sql.DB, model string, id int) error {
	switch model {
	case "user":
		m := &mysql.UserModel{DB: db}
		s, err := m.Get(id)
		if err != nil {
			return err
		}
		mo, err := json.MarshalIndent(s, "", "   ")
		if err != nil {
			return err
		}
		fmt.Printf("%+v \n", string(mo))
	case "snippet":
		m := &mysql.SnippetModel{DB: db}
		s, err := m.Get(id)
		if err != nil {
			return err
		}
		mo, err := json.MarshalIndent(s, "", "   ")
		if err != nil {
			return err
		}
		fmt.Printf("%+v \n", string(mo))
	}
	return nil
}
