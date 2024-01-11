package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	_ "github.com/microsoft/go-mssqldb"
)

type serverInstance struct {
	instanceName string
	encrypt      bool
}

func main() {
	serverinstance, err := validateUserInput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conString := createConnectionString(serverinstance)
	if err := testConnection(conString); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("succesfully connected to: %s\n", serverinstance.instanceName)
}

func validateUserInput() (serverInstance, error) {
	// We need a servername to continue
	if len(os.Args) < 2 {
		return serverInstance{}, errors.New("a server/instance name is required")
	}
	flag.Usage = appUsage
	// We need a flag to know if we are going for encryption or not
	// These flag(s) will show up with --help
	encrypt := flag.Bool("encrypt", false, "Force connection encryption")

	flag.Parse() //Parses the flags from the terminal

	instanceName := flag.Arg(0)

	// We return our instance information
	return serverInstance{instanceName, *encrypt}, nil
}

func appUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [Options] instance\n", filepath.Base(os.Args[0]))
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Simple tool to test SQL server connectivity from the command line\n")
	fmt.Fprintf(os.Stderr, "Windows auth is used and encryption can be forced\n")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
}

func createConnectionString(serverInfo serverInstance) string {
	query := url.Values{}
	query.Add("app name", "sqltester")
	// We use "disable" and "mandatory" below as that is what MSFT uses (We don't use "optional" as it has some weird side effects)
	if !serverInfo.encrypt {
		query.Add("encrypt", "disable")
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

func testConnection(connectionString string) error {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return fmt.Errorf("failed to process connection string: %v", err)
	}
	defer db.Close()
	// We need to ping the db, open just checks that the connection string is valid
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	return nil
}
