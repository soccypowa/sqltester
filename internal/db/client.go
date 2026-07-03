package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

type Client struct {
	db *sql.DB
}

func New(connstr string) (*Client, error) {
	u, err := url.Parse(connstr)
	if err != nil {
		return nil, fmt.Errorf("invalid connection string: %v", err)
	}
	q := u.Query()
	if q.Get("app name") == "" {
		q.Set("app name", "sqltester")
	}
	u.RawQuery = q.Encode()
	connstr = u.String()

	db, err := sql.Open("sqlserver", connstr)
	if err != nil {
		return nil, fmt.Errorf("could not open connection: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("db ping failed: %v", err)
	}

	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(1)
	db.SetConnMaxIdleTime(60 * time.Second)

	return &Client{
		db: db,
	}, nil
}

func (c *Client) Close() {
	c.db.Close()
}

type ConnectionInfo struct {
	SPID             int
	ServerName       string
	ServiceName      string
	Database         string
	AuthScheme       string
	EncryptionOption bool
}

func (c *Client) GetServerInfo(ctx context.Context) (ConnectionInfo, error) {
	const query = `
SELECT TOP(1)
	s.session_id,
	@@SERVERNAME 						AS server_name,
	@@SERVICENAME 						AS service_name,
	ISNULL(DB_NAME(s.database_id), '') 	AS database_name,
	ISNULL(c.auth_scheme, '') 			AS auth_scheme,
	ISNULL(c.encrypt_option, '') 		AS encryption_option
FROM sys.dm_exec_sessions s
	join sys.dm_exec_connections c
	ON s.session_id = c.session_id
WHERE s.session_id = @@SPID
	`
	var r ConnectionInfo
	err := c.db.QueryRowContext(ctx, query).Scan(&r.SPID, &r.ServerName, &r.ServiceName, &r.Database, &r.AuthScheme, &r.EncryptionOption)

	return r, err
}

func (c ConnectionInfo) String() string {
	return fmt.Sprintf(
		"received this response from the remote host:\n"+
			" • spid:\t\t%d\n"+
			" • server_name:\t\t%s\n"+
			" • service_name:\t%s\n"+
			" • database_name:\t%s\n"+
			" • auth_scheme:\t\t%s\n"+
			" • encryption_option:\t%t\n",
		c.SPID, c.ServerName, c.ServiceName, c.Database, c.AuthScheme, c.EncryptionOption,
	)
}
