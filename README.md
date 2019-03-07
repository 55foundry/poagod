# Introduction

A utility for running the steps to deploying an Ethereum PoA network. 

[![CircleCI](https://circleci.com/gh/55foundry/poagod/tree/master.svg?style=svg)](https://circleci.com/gh/55foundry/poagod/tree/master)


# Getting started
#### Installing dependencies
```
go get -v -t -d ./...
```

#### Testing the code
```
go test -v ./...
```

## > Genesis
First application of PoAGod is the ability to deploy a new genesis. The corresponding command to invoke, is `genesis`

You are required to have a file that contains ETH addresses in order to build your Genesis that assigns prefunds to accounts. The format of your `addresses.txt` should be of the following:

```
Address: {5a13b2f165004d1866eca8503e080352f63925a0}
Address: {029b8f68a60bc794e5c4156a8df02140919db833}
Address: {79fd6ae5fc99dfd80be9c7302a0058c6543b66f5}
Address: {2ea4c67445923800ac02fa67f15b5b48e6f2f47c}
Address: {03e332823530e85e7f0e69c4173288781a1f6989}
Address: {2361c52d851622167306a8728fe3fccab4b29097}
Address: {634e1f77e57a8d366179927c52c670ceba146b28}
Address: {8e0772e3fbe31491c6bfb37b9cfe9b1b00217800}
Address: {c39f84525577e3d90876c685d52e40c2d4490f6a}
Address: {bb47999fa01bed8a26abeb039df15a105c8cfa49}
```

You can have `n` number of addresses in the file, just do not have none. You can generate this file based on this command:

```
for ((n=0;n<10;n++)); do geth account new --password {{YOUR PASSFILE}} >> addresses.txt; done
```

Once generated, invoke the CLI to generate a valid Ethereum and Harmony Genesis file:

```
poagod genesis -create -addresses ./addresses.txt
```

## > Node
This is a work in progress. The idea behind this is:

1) Start a bootnode.

2) Start `n` sealer nodes.

3) Peer each node to one another.