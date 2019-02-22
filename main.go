package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
)

// Log - Global Logger variable to use with Logrus instance
var Log *logrus.Logger

func main() {
	Log = LoadLogger()

	genesisCommand := flag.NewFlagSet("genesis", flag.ExitOnError)
	genesisCommand.Usage = func() {
	    fmt.Fprintf(os.Stderr, "\033[32mgenesis\033[0m - Manage, create/import a genesis json file.\n")
	}
    nodeCommand := flag.NewFlagSet("node", flag.ExitOnError)
	nodeCommand.Usage = func() {
		fmt.Fprintf(os.Stderr, "\033[32mnode\033[0m - Deploy ethereum nodes and connect to poa.\n")
	}

	genesisCreate := genesisCommand.Bool("create", false, "to create a genesis block file")
	genesisCreateAddresses := genesisCommand.String("addresses", "./addresses.txt", "path to geth addresses file")

	if len(os.Args) < 2 {
		genesisCommand.Usage()
	    genesisCommand.PrintDefaults()

		nodeCommand.Usage()
		nodeCommand.PrintDefaults()

	    os.Exit(1)
	}

	switch os.Args[1] {
    case "genesis":
        genesisCommand.Parse(os.Args[2:])

		if (*genesisCreate == true) {
			w := wizard{network: GetEnv("ACCOUNT_ID", "55f")}
			w.makeGenesis(*genesisCreateAddresses)
			w.saveGenesis(w.conf)
			os.Exit(1)
		}

		genesisCommand.Usage()
		genesisCommand.PrintDefaults()
    case "node":
        nodeCommand.Parse(os.Args[2:])

		nodeCommand.Usage()
		nodeCommand.PrintDefaults()
    default:
        flag.PrintDefaults()
        os.Exit(1)
    }
}
