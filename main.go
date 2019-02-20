package main

import (
	"flag"
	"os"

	"github.com/Sirupsen/logrus"
)

// Log - Global Logger variable to use with Logrus instance
var Log *logrus.Logger

func main() {
	Log = LoadLogger()

	genesisCommand := flag.NewFlagSet("genesis", flag.ExitOnError)
    nodeCommand := flag.NewFlagSet("node", flag.ExitOnError)

	genesisCreate := genesisCommand.Bool("create", false, "to create a genesis block file")
	genesisCreateAddresses := genesisCommand.String("addresses", "./addresses.txt", "path to geth addresses file")

	if len(os.Args) < 2 {
	    Log.Warn("a command is required to continue")
	    os.Exit(1)
	}

	switch os.Args[1] {
    case "genesis":
        genesisCommand.Parse(os.Args[2:])

		if (*genesisCreate == true) {
			Log.Info("Initializing POAGod...")

			w := wizard{network: GetEnv("ACCOUNT_ID", "55f")}
			w.makeGenesis(*genesisCreateAddresses)
			w.saveGenesis(w.conf)
			os.Exit(1)
		}

		genesisCommand.PrintDefaults()
    case "node":
        nodeCommand.Parse(os.Args[2:])
		nodeCommand.PrintDefaults()
    default:
        flag.PrintDefaults()
        os.Exit(1)
    }
}
