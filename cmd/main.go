package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/soccypowa/sqltester/internal/db"
)

func main() {
	connStr := flag.String("conn", "", "SQL Server connection string\n"+
		"  sqlserver://user:password@host[:port][/instance]\n"+
		"  sqlserver://host?trusted_connection=true  (Windows auth)\n"+
		"  For further information checkout: https://github.com/microsoft/go-mssqldb")
	flag.Parse()
	if *connStr == "" {
		fmt.Fprintln(os.Stderr, "Error: provide --conn")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, `  sqlmon --conn "sqlserver://sa:password@localhost"`)
		fmt.Fprintln(os.Stderr, `  sqlmon --conn "sqlserver://myserver:1433" --interval 10s`)
		fmt.Fprintln(os.Stderr, `  sqlmon --conn "sqlserver://myserver?trusted_connection=true"`)
		fmt.Fprintln(os.Stderr, "  For further information checkout: https://github.com/microsoft/go-mssqldb")
		os.Exit(1)
	}

	name := extractHost(*connStr)

	fmt.Fprintf(os.Stderr, "connecting to %s...\n\n", name)
	client, err := db.New(*connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nCould not connect: %v\n\n", err)
		fmt.Fprintln(os.Stderr, "Things to check:")
		fmt.Fprintln(os.Stderr, "  • SQL Server is running and reachable")
		fmt.Fprintln(os.Stderr, "  • Username / password are correct")
		fmt.Fprintln(os.Stderr, "  • Login has VIEW SERVER STATE permission:")
		fmt.Fprintln(os.Stderr, "      GRANT VIEW SERVER STATE TO [yourlogin];")
		os.Exit(1)
	}
	defer client.Close()

	response, err := client.GetServerInfo(context.Background())
	if err != nil {
		fmt.Printf("received a error from the remote host: %v", err)
		os.Exit(1)
	}
	fmt.Printf(
		"Received this response from the remote host:\n"+
			" • spid:\t\t%d\n"+
			" • server_name:\t\t%s\n"+
			" • service_name:\t%s\n"+
			" • database_name:\t%s\n"+
			" • auth_scheme:\t\t%s\n"+
			" • encryption_option:\t%t\n",
		response.SPID, response.ServerName, response.ServiceName, response.Database, response.AuthScheme, response.EncryptionOption,
	)
}

func extractHost(connStr string) string {
	s := connStr

	// Strip Scheme
	if i := strings.Index(s, "://"); i >= 0 {
		s = s[i+3:]
	}

	// Strip user:pass@ - find the last @ before any / or ?
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '@' {
			s = s[i+1:]
			break
		}
		if s[i] == '/' || s[i] == '?' {
			break
		}
	}
	// Strip query string
	if i := strings.IndexByte(s, '?'); i >= 0 {
		s = s[:i]
	}

	if s == "" {
		return "SQL Server"
	}
	return s
}
