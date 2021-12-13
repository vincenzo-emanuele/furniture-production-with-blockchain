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
//	"encoding/json"
	"strings"
//	"strconv"
	"bufio"
//	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)



func main() {
	log.Println("============ AVVIO app Org3MSP ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet3")
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
		"org3.example.com",
		"connection-org3.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()
	
	network, err := gw.GetNetwork("chanorg2org3")

	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("basic-23")
	
	
	
	reader := bufio.NewReader(os.Stdin)
	

	log.Println("--> Controllo se il ledger è stato inizializzato")
	_, err = contract.EvaluateTransaction("ReadWood")
	if err != nil {
		log.Println("--> Inizializzo ledger")
		_, err := contract.SubmitTransaction("InitLedger")
		if err != nil {
			log.Fatalf("Failed to Submit transaction: %v", err)
		}
	} else {
		log.Println("--> Ledger già inizializzato")
	}
	
	
	//Inizio
	for {
		fmt.Print("Cosa vuoi fare?\n1. InitLedger\n2. GetAllWallets\n3. ReadWood\n4. RemoveWood\n5. Mint\n6. TransferWood\n7. Exit\n: ")
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
				log.Println("Wallets: " + string(result))
				fmt.Print("\n")
				
			case "3":
				result, err := contract.EvaluateTransaction("ReadWood")
				if err != nil {
					log.Println("Failed to evaluate transaction ReadWood: %v", err)
					break
				}
				log.Println("Wood: " + string(result))
				fmt.Print("\n")
				
			case "4":
				fmt.Print("Inserisci amount: ")
				amount, _ := reader.ReadString('\n')
				amount = strings.Replace(amount, "\n", "", -1)
				result, err := contract.SubmitTransaction("RemoveWood", amount)
				if err != nil {
					log.Println("Failed to evaluate transaction ReadWood: %v", err)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result) + "\n")
			case "5":
				fmt.Print("Inserisci amount: ")
				amount, _ := reader.ReadString('\n')
				amount = strings.Replace(amount, "\n", "", -1)
				result, err := contract.SubmitTransaction("Mint", amount)
				if err != nil {
					log.Println("Failed to evaluate transaction Mint: %v", err)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result) + "\n")
			case "6":
				fmt.Print("Inserisci amount: ")
				amount, _ := reader.ReadString('\n')
				amount = strings.Replace(amount, "\n", "", -1)
				fmt.Print("Inserisci destinatario (Org2MSP o Org3MSP): ")
				dest, _ := reader.ReadString('\n')
				dest = strings.Replace(dest, "\n", "", -1)
				
				result, err := contract.SubmitTransaction("TransferWood", amount, dest)
				if err != nil {
					log.Println("Failed to evaluate transaction Mint: %v", err)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result) + "\n")
			
				
		}//switch
		
		if op == "7" {
			break
		}
	}//for
	
	log.Println("============ CHIUSURA APP Org3MSP ============")
	
	log.Println("Cancello la cartella appena creata")
	err = os.RemoveAll("./wallet3")
	if err != nil {
		log.Fatalf("ERRORE: %v", err)
	}
	


}


func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org3.example.com",
		"users",
		"User1@org3.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "User1@org3.example.com-cert.pem")
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

	identity := gateway.NewX509Identity("Org3MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}



