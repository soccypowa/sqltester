package main

import (
	"errors"
	"flag"
	"os"
)

type serverInstance struct {
	instanceName string
	encrypt      bool
}

func main() {

}

func getConnectionData() (serverInstance, error) {
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
