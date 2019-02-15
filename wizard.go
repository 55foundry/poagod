package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
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
		log.Warn("Failed to save puppeth configs", "file", c.path, "err", err)
	}
}

type wizard struct {
	network string // Network name to manage
	conf    config // Configurations from previous runs

	services map[string][]string   // Ethereum services known to be running on servers

	in   *bufio.Reader // Wrapper around stdin to allow reading user input
	lock sync.Mutex    // Lock to protect configs during concurrent service discovery
}

func (w *wizard) readAddress(a string) *common.Address {
	for {
		if a = strings.TrimSpace(a); a == "" {
			return nil
		}

		if len(a) != 40 {
			log.Error("Invalid address length, please retry")
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
