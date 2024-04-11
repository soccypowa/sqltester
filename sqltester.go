package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
)

type serverInstance struct {
	instanceName string
	encrypt      bool
}

func main() {
	serverinstance, err := validateUserInput()
	checkErr(err)

	conString := createConnectionString(serverinstance)

	db, err := openConnection(conString)
	checkErr(err)

	err = dbPinger(db)
	checkErr(err)

	db.Close()
}

func validateUserInput() (serverInstance, error) {
	// We need a servername to continue
	if len(os.Args) < 2 {
		return serverInstance{}, errors.New("a server/instance name is required (use fqdn)")
	}
	flag.Usage = appUsage
	// We need a flag to know if we are going for encryption or not
	// These flag(s) will show up with --help
	encrypt := flag.Bool("encrypt", false, "sets encrypt=mandatory (default: encrypt=disable)")

	flag.Parse() //Parses the flags from the terminal

	instanceName := flag.Arg(0)

	// We return our instance information
	return serverInstance{instanceName, *encrypt}, nil
}

func appUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [Options] instance\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "Example: %s myserver.mydomain.com\n", filepath.Base(os.Args[0]))
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Simple tool to test SQL server connectivity from the command line\n")
	fmt.Fprintf(os.Stderr, "Windows auth is used and encryption can be forced\n")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
}

func createConnectionString(serverInfo serverInstance) string {
	query := url.Values{}
	//TODO: Hardcode ttl
	query.Add("app name", "sqltester")
	// We use "disable" and "mandatory" below as that is what MSFT uses (We don't use "optional" as it has some weird side effects)
	if !serverInfo.encrypt {
		query.Add("encrypt", "disable")
	} else {
		query.Add("encrypt", "mandatory")
	}

	// To handle named instances vi have to split the string on the backslash
	var host string
	var path string
	if !strings.Contains(serverInfo.instanceName, `\`) {
		host = strings.Split(serverInfo.instanceName, `\`)[0]
	} else {
		host = strings.Split(serverInfo.instanceName, `\`)[0]
		path = strings.Split(serverInfo.instanceName, `\`)[1]
	}

	// Construct the url, empty variables is omitted
	url := url.URL{
		Scheme:   "sqlserver",
		Host:     host,
		Path:     path,
		RawQuery: query.Encode(),
	}
	return url.String()
}

func openConnection(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return &sql.DB{}, fmt.Errorf("failed to process connection string: %v", err)
	}
	return db, nil
}

func dbPinger(db *sql.DB) error {
	// We need to ping the db, open just checks that the connection string is valid
	fmt.Println("Ping ->")
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	fmt.Printf("\t<- Pong\n")
	return nil
}

func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, err.Error())
	os.Exit(1)
}

func checkErr(err error) {
	if err != nil {
		exitGracefully(err)
	}
}
