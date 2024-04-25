package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestValidateUserInput(t *testing.T) {
	tests := []struct {
		name    string
		want    serverInstance
		wantErr bool
		osArgs  []string
	}{
		{name: "Default params", want: serverInstance{"sql.test.com", false}, wantErr: false, osArgs: []string{"cmd", "sql.test.com"}},
		{name: "No parameters", want: serverInstance{}, wantErr: true, osArgs: []string{"cmd"}},
		{name: "Encryption enabled", want: serverInstance{"sql.test.com", true}, wantErr: false, osArgs: []string{"cmd", "--encrypt", "sql.test.com"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// saving original os args for ref
			origOsArgs := os.Args

			// function to run after or tests are done
			defer func() {
				os.Args = origOsArgs
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			}()

			os.Args = tc.osArgs // setting os args for test
			got, err := validateUserInput()
			if (err != nil) != tc.wantErr {
				t.Errorf("getConnectionData() error = %v, wantErr = %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("getConnectionData() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCreateConnectionString(t *testing.T) {
	tests := []struct {
		name string
		want string
		args serverInstance
	}{
		{name: "Default instance", want: "sqlserver://sql.test.com?app+name=sqltester&encrypt=disable&keepalive=35", args: serverInstance{instanceName: "sql.test.com"}},
		{name: "Named instance", want: "sqlserver://sql.test.com/namedinstance?app+name=sqltester&encrypt=disable&keepalive=35", args: serverInstance{instanceName: `sql.test.com\namedinstance`}},
		{name: "With encryption", want: "sqlserver://sql.test.com?app+name=sqltester&encrypt=mandatory&keepalive=35", args: serverInstance{instanceName: "sql.test.com", encrypt: true}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := createConnectionString(tc.args)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("createConnectionString() got = %v, want %v", got, tc.want)
			}
		})
	}
}
