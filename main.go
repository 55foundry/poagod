package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
)

func (w *wizard) makeGenesis() {
	fmt.Println("Building new genesis block")

	genesis := &core.Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		GasLimit:   4700000,
		Difficulty: big.NewInt(1), //big.NewInt(524288),
		Alloc:      make(core.GenesisAlloc),
		Config: &params.ChainConfig{
			HomesteadBlock:      big.NewInt(1),
			EIP150Block:         big.NewInt(2),
			EIP155Block:         big.NewInt(3),
			EIP158Block:         big.NewInt(3),
			ByzantiumBlock:      big.NewInt(4),
			ConstantinopleBlock: big.NewInt(5),
			Clique: &params.CliqueConfig{
				Period: 5,
				Epoch:  30000,
			},
		},
	}

	fileP := "./addresses.txt"
	if _, err := os.Stat(fileP); os.IsNotExist(err) {
	  panic(err)
	}

	file, err := os.Open(fileP)
	if err != nil {
	    panic(err)
	}
	defer file.Close()

	var signers []common.Address

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		match := regexp.MustCompile("\\{(.*?)\\}").FindStringSubmatch(line)[1]
		address := *w.readAddress(match)

		signers = append(signers, address)

		genesis.Alloc[address] = core.GenesisAccount{
			Balance: new(big.Int).Lsh(big.NewInt(1), 256-7),
		}
	}

	if err := scanner.Err(); err != nil {
	    panic(err)
	}

	genesis.ExtraData = make([]byte, 32+len(signers)*common.AddressLength+65)
	for i, signer := range signers {
		copy(genesis.ExtraData[32+i*common.AddressLength:], signer[:])
	}

	for i := int64(0); i < 256; i++ {
		genesis.Alloc[common.BigToAddress(big.NewInt(i))] = core.GenesisAccount{Balance: big.NewInt(1)}
	}

	genesis.Config.ChainID = new(big.Int).SetUint64(uint64(555))

	fmt.Println("Configured new genesis block")

	w.conf.Genesis = genesis
	w.conf.flush()
}

func saveGenesis(folder, network, client string, spec interface{}) {
	var format string

	if (client == "") {
		format = fmt.Sprintf("%s.json", network)
	} else {
		format = fmt.Sprintf("%s-%s.json", network, client)
	}

	path := filepath.Join(folder, format)

	out, _ := json.Marshal(spec)
	if err := ioutil.WriteFile(path, out, 0644); err != nil {
		fmt.Println("Failed to save genesis file", "client", client, "err", err)
		return
	}

	fmt.Println("Saved genesis chain spec", "client", client, "path", path)
}

func main() {
	fmt.Println("Initializing BUtility...")

	w := wizard{network: GetEnv("ACCOUNT_ID", "55f")}
	w.makeGenesis()

	folder, _ := os.Getwd()

	saveGenesis(folder, w.network, "", w.conf)
}
