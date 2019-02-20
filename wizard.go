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
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

type config struct {
	path      string   // File containing the configuration values
	bootnodes []string // Bootnodes to always connect to by all nodes
	ethstats  string   // Ethstats settings to cache for node deploys

	Genesis *core.Genesis     `json:"genesis,omitempty"` // Genesis block to cache for node deploys
	Servers map[string][]byte `json:"servers,omitempty"`
}

func (c config) flush() {
	os.MkdirAll(filepath.Dir(c.path), 0755)

	out, _ := json.MarshalIndent(c, "", "  ")
	if err := ioutil.WriteFile(c.path, out, 0644); err != nil {
		Log.Warn("Failed to save puppeth configs", "file", c.path, "err", err)
	}
}

type wizard struct {
	network string // Network name to manage
	conf    config // Configurations from previous runs

	services map[string][]string   // Ethereum services known to be running on servers

	in   *bufio.Reader // Wrapper around stdin to allow reading user input
	lock sync.Mutex    // Lock to protect configs during concurrent service discovery
}

func (w *wizard) makeGenesis(addresses string) {
	Log.Info("Building new genesis block")

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

	if _, err := os.Stat(addresses); os.IsNotExist(err) {
	  panic(err)
	}

	file, _ := os.Open(addresses)
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

	Log.Info("Configured new genesis block")

	w.conf.Genesis = genesis
	w.conf.flush()
}

func (w *wizard) readAddress(a string) *common.Address {
	for {
		if a = strings.TrimSpace(a); a == "" {
			return nil
		}

		if len(a) != 40 {
			Log.Error("Invalid address length, please retry")
			continue
		}

		bigaddr, _ := new(big.Int).SetString(a, 16)
		address := common.BigToAddress(bigaddr)

		return &address
	}
}

func (w *wizard) readString(a string) string {
	for {
		if a = strings.TrimSpace(a); a != "" {
			return a
		}
	}
}

func (w *wizard) saveGenesis(spec interface{}) {
	client := ""
	format := fmt.Sprintf("%s.json", w.network)
	folder, _ := os.Getwd()

	path := filepath.Join(folder, format)

	out, _ := json.Marshal(spec)
	if err := ioutil.WriteFile(path, out, 0644); err != nil {
		Log.Warn("Failed to save genesis file", "client", client, "err", err)
		return
	}

	Log.Info("Saved genesis chain spec", "client", client, "path", path)
}
