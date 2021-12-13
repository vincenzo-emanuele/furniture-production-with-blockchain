/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"encoding/json"
	"strings"
	"strconv"
	"bufio"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type Asset struct {
	ID             string `json:"ID"`
	Price		int    `json:"Price"`
	Color          string `json:"Color"`
	Type           string    `json:"Type"`
}

//Realizzo un wallet di Token NFT
type Wallet struct{
	Owner string	`json:"Owner"`
	NFT	map[string]Asset	`json:"NFT"`
}

func main() {
	//Questa roba è tutta magia
	log.Println("============ AVVIO app Org1MSP ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()
//Ottengo il riferimento al canale che mi interessa
	network, err := gw.GetNetwork("chanorg1org2")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
//Ottengo il riferimento al chaincode che mi interessa
	contract := network.GetContract("basic-12")

	reg, notifier, err := contract.RegisterEvent("DeleteAsset")
	if err != nil {
		fmt.Printf("Failed to register contract event: %s", err)
		return
	}
	//fmt.Println("Reg: ", reg)
	//fmt.Println("Notifier: ", notifier)
	defer contract.Unregister(reg)


	go Evento(notifier)
	
	reader := bufio.NewReader(os.Stdin)
	
	log.Println("--> Controllo se il ledger è stato inizializzato")
	result , err := contract.EvaluateTransaction("GetAllWallets")
	var wallets []Wallet
	err = json.Unmarshal(result, &wallets)
	if err != nil {
		log.Println("--> Inizializzo il ledger")
		_ , err = contract.SubmitTransaction("InitLedger")
		if err != nil {
			log.Fatalf("Failed to Submit transaction: %v", err)
			return 
		}
	} else {
		log.Println("--> Ledger già inizializzato")
	}
	
	
	//Inizio
	for {
		fmt.Print("Cosa vuoi fare?\n1. InitLedger\n2. GetAllWallets\n3. ReadAsset\n4. UpdateAsset\n5. MockSell\n6. AssetExists\n7. TransferAsset\n8. Exit\n: ")
		op, _ := reader.ReadString('\n')
		op = strings.Replace(op, "\n", "", -1)
		switch op {
			case "1":
				log.Println("--> Submit Transaction: InitLedger")
				result, err := contract.SubmitTransaction("InitLedger")
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result) + "\n")
				
			case "2":
				log.Println("--> Evaluate Transaction: GetAllWallets")
				result, err := contract.EvaluateTransaction("GetAllWallets")
				if err != nil {
					log.Println("Failed to evaluate transaction: %v", err)
					break
				}
				log.Println(string(result))
				fmt.Print("\n")
				
			case "3":
				log.Println("--> Evaluate Transaction: ReadAsset")
				fmt.Print("Inserisci ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				result, err := contract.EvaluateTransaction("ReadAsset", id)
				if err != nil {
					log.Println("Failed to evaluate transaction: %v\n", err)
					break
				}
				log.Println(string(result))
				fmt.Print("\n")
				
			case "4":
				log.Println("--> Evaluate Transaction: UpdateAsset")
				fmt.Print("Inserisci ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				fmt.Print("Inserisci colore: ")
				color, _ := reader.ReadString('\n')
				color = strings.Replace(color, "\n", "", -1)
				fmt.Print("Inserisci tipo: ")
				Type, _ := reader.ReadString('\n')
				Type = strings.Replace(Type, "\n", "", -1)
				fmt.Print("Inserisci prezzo: ")
				price, _ := reader.ReadString('\n')
				price = strings.Replace(price, "\n", "", -1)
				result, err := contract.SubmitTransaction("UpdateAsset", id, color, Type, price)
				if err != nil {
					log.Println("Failed to evaluate transaction: %v\n", err)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result) + "\n")
				
			case "5":
				fmt.Print("Inserisci ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				result, err := contract.SubmitTransaction("DeleteAsset", id)
				if err != nil {
					log.Println("Failed to evaluate transaction: %v", err)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result) + "\n")
				
			case "6":
				log.Println("--> Evaluate Transaction: AssetExists")
				fmt.Print("Inserisci ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				result, err := contract.EvaluateTransaction("AssetExists", id)
				if err != nil {
					log.Println("Failed to evaluate transaction: %v\n", err)
					break
				}
				if string(result) == "true" {
					fmt.Println("L'asset con ID " + id + " esiste")
				} else if string(result) == "false" {
					fmt.Println("L'asset con ID " + id + " NON esiste")
				} else {
					fmt.Println("Errore")
				}
				fmt.Print("\n")
				
			case "7":
				log.Println("--> Submit Transaction: TransferAsset")
				fmt.Print("Inserisci ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.Replace(id, "\n", "", -1)
				fmt.Print("Inserisci destinatario (Org1MSP o Org2MSP): ")
				dest, _ := reader.ReadString('\n')
				dest = strings.Replace(dest, "\n", "", -1)
				_, err := contract.SubmitTransaction("TransferAsset", id, dest)
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result) + "\n")
		}//switch
		
		if op == "8" {
			break
		}
	}//for
	
	
	log.Println("============ CHIUSURA APP Org1MSP ============")
	
	log.Println("Cancello la cartella appena creata")
	err = os.RemoveAll("./wallet")
	if err != nil {
		log.Fatalf("ERRORE: %v", err)
	}
	
	
}


func Evento(notifier <-chan *fab.CCEvent) (error){
		
		var ccEvent *fab.CCEvent
		for {
			select {
			case ccEvent = <-notifier:
				res, err := strconv.Atoi(string(ccEvent.Payload))
				if err != nil {
					log.Println("Errore di conversione: " + string(ccEvent.Payload))
				}
				
				if res > 2 {
					log.Println("Hai ancora più di 2 mobili")
				} else {
					log.Println("Oh no, refill")
					err = doRefill()
					if err != nil {
						log.Println("Failed to Submit transaction: %v", err)
					}
				}
				
			}//select
		
		}//for
		return nil
	}

func doRefill() (error){

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
		return err
	}

	wallet, err := gateway.NewFileSystemWallet("wallet2")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
		return err
	}

	if !wallet.Exists("appUser") {
		err = populateWallet2(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
			return err
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org2.example.com",
		"connection-org2.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
		return err
	}
	defer gw.Close()
//Ottengo il riferimento al canale che mi interessa
	network, err := gw.GetNetwork("chanorg1org2")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
		return err
	}
//Ottengo il riferimento al chaincode che mi interessa
	contract := network.GetContract("basic-12")

	
	result, err := contract.EvaluateTransaction("GetAllWallets")
	var w []Wallet
	err = json.Unmarshal(result, &w)
	if err != nil {
		log.Println("Failed to evaluate transaction: %v", err)
	}
	var index int
	
	if w[0].Owner == "Org2MSP" {
		index = 0
	} else {
		index = 1
	}
	
	if l := len(w[index].NFT); l != 0 {
	
		for i, asset := range w[index].NFT {
			
			_, err := contract.SubmitTransaction("TransferAsset", asset.ID, "Org1MSP")
			if err != nil {
				log.Println("Failed to Submit transaction: %v", err)
				break
			}
			y, _ := strconv.Atoi(i)
			if  y >= 2 {
				break
			}
		}//for
		
		log.Println("Refill effettuato con successo")
		fmt.Print("\n")
    	} else {
    		log.Println("Org2MSP non ha più asset nel wallet...")
    	}
	
	
	
	return nil
}


func populateWallet2(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org2.example.com",
		"users",
		"User1@org2.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "User1@org2.example.com-cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org2MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}


func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "User1@org1.example.com-cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}

