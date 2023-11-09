package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"

	_ "github.com/microsoft/go-mssqldb"
)

type serverInstance struct {
	instanceName string
	encrypt      bool
}

var DB *sql.DB

func main() {
	serverinstance, err := validateUserInput()
	if err != nil {
		fmt.Printf("failed to validate input: %v", err)
	}
	conString := createConnectionString(serverinstance)
	connect(conString)
	defer DB.Close()
}

func validateUserInput() (serverInstance, error) {
	// We need a servername to continue
	if len(os.Args) < 2 {
		return serverInstance{}, errors.New("a server/instance name is required")
	}

	// We need a flag to know if we are going for encryption or not
	// These flag(s) will show up with --help
	encrypt := flag.Bool("encrypt", false, "Use connection encryption")

	flag.Parse() //Parses the flags from the terminal

	instanceName := flag.Arg(0)

	// We return our instance information
	return serverInstance{instanceName, *encrypt}, nil
}

func createConnectionString(serverInfo serverInstance) string {
	query := url.Values{}
	query.Add("app name", "sqltester")
	if !serverInfo.encrypt {
		query.Add("encrypt", "Optional")
	} else {
		query.Add("encrypt", "Mandatory")
	}
	url := url.URL{
		Scheme:   "sqlserver",
		Host:     serverInfo.instanceName,
		RawQuery: query.Encode(),
	}
	return url.String()
}

func connect(connectionString string) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		fmt.Printf("failed to connect to server: %v", err)
	}
	DB = db
	fmt.Print("Connected!")
}
