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
	log.Println("============ AVVIO app Org2MSP ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet2")
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
		"org2.example.com",
		"connection-org2.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("chanorg1org2")
	
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
	
	network2, err := gw.GetNetwork("chanorg2org3")

	if err != nil {
		log.Fatalf("Failed to get network2: %v", err)
	}

	contract := network.GetContract("basic-12")

	contract2 := network2.GetContract("basic-23")


	regCreate, notifierCreate, err := contract.RegisterEvent("CreateAsset")
	if err != nil {
		fmt.Printf("Failed to register contract event: %s", err)
		return
	}
	
	regRemove, notifierRemove, err := contract2.RegisterEvent("RemoveWood")
	if err != nil {
		fmt.Printf("Failed to register contract event: %s", err)
		return
	}


	defer contract.Unregister(regCreate)
	defer contract.Unregister(regRemove)

	go EventoCreate(notifierCreate, contract, contract2)
	go EventoRemove(notifierRemove, contract, contract2)
	

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
	
	log.Println("--> Controllo se il ledger 2 è stato inizializzato")
	_, err = contract2.EvaluateTransaction("ReadWood")
	if err != nil {
		log.Println("--> Inizializzo ledger 2")
		_, err := contract2.SubmitTransaction("InitLedger")
		if err != nil {
			log.Fatalf("Failed to Submit transaction: %v", err)
		}
	} else {
		log.Println("--> Ledger 2 già inizializzato")
	}
	
	
	//Inizio
	for {
		fmt.Print("Cosa vuoi fare?\n1. InitLedger\n2. GetAllWallets\n3. ReadAsset\n4. UpdateAsset\n5. DeleteAsset\n6. AssetExists\n7. TransferAsset\n8. CreateAsset\n9. Initledger2\n10. GetAllWallets2\n11. ReadWood\n12. MockUseWood\n13. Exit\n: ")
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
				
			case "8":
				result, err = contract2.EvaluateTransaction("ReadWood")
				if err != nil {
					log.Println("Failed to evaluate transaction ReadWood: %v", err)
					break
				}
				stringRes := string(result)
				intRes, err := strconv.Atoi(stringRes)
				if err != nil {
					log.Println("Errore di conversione: %v", err)
					break
				}
				
				if intRes >= 100 {
					log.Println("--> Submit Transaction: CreateAsset")
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
					result, err = contract.SubmitTransaction("CreateAsset", id, color, Type, price)
					if err != nil {
						log.Println("Failed to Submit transaction: %v", err)
						break
					}
					log.Println("ESEGUITO")
					fmt.Print(string(result) + "\n")
				} else {
					log.Println("Non Hai abbastanza legna")
				}
			case "9":
				log.Println("--> Submit Transaction: InitLedger2")
				result, err := contract2.SubmitTransaction("InitLedger")
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result) + "\n")
			case "10":
				log.Println("--> Evaluate Transaction: GetAllWallets2")
				result, err := contract2.EvaluateTransaction("GetAllWallets")
				if err != nil {
					log.Println("Failed to evaluate transaction: %v", err)
					break
				}
				log.Println(string(result))
				fmt.Print("\n")
			case "11":
				result, err := contract2.EvaluateTransaction("ReadWood")
				if err != nil {
					log.Println("Failed to evaluate transaction ReadWood: %v", err)
					break
				}
				log.Println("Wood: " + string(result))
				fmt.Print("\n")
				
			case "12":
				fmt.Print("Inserisci amount: ")
				amount, _ := reader.ReadString('\n')
				amount = strings.Replace(amount, "\n", "", -1)
				result, err := contract2.SubmitTransaction("RemoveWood", amount)
				if err != nil {
					log.Println("Failed to evaluate transaction ReadWood: %v", err)
					break
				}
				log.Println("ESEGUITO")
				fmt.Print(string(result) + "\n")
		}//switch
		
		if op == "13" {
			break
		}
	}//for
	
	log.Println("============ CHIUSURA APP Org2MSP ============")
	
	log.Println("Cancello la cartella appena creata")
	err = os.RemoveAll("./wallet2")
	if err != nil {
		log.Fatalf("ERRORE: %v", err)
	}
}

func EventoCreate(notifier <-chan *fab.CCEvent, contract1 *gateway.Contract, contract2 *gateway.Contract) (error){
		
	var ccEvent *fab.CCEvent
	for {
		select {
		case ccEvent = <-notifier:
			_ = ccEvent.Payload
			_, err := contract2.SubmitTransaction("RemoveWood", "100")
			if err != nil {
				log.Fatalf("Failed to evaluate transaction RemoveWood: %v", err)
			}
			
			log.Println("Eseguito il pagamento di 100 wood")
			
		}//select

	}//for
	return nil
}

func EventoRemove(notifier <-chan *fab.CCEvent, contract1 *gateway.Contract, contract2 *gateway.Contract) (error){
	var ccEvent *fab.CCEvent
	for {
		select {
		case ccEvent = <-notifier:
			res := ccEvent.Payload
			resInt, err := strconv.Atoi(string(res))
			if err != nil {
					log.Println("Errore di conversione: " + string(res))
			}
			if resInt < 300 {
				err = doRefill()
				if err != nil {
					log.Println("Errore nel refill: %v", err)
				}
			} else {
				log.Println("Hai ancora abbastanza legna")
			}
			
		}//select

	}//for
	return nil
}

func doRefill() (error) {
	
	wallet, err := gateway.NewFileSystemWallet("wallet3")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet2(wallet)
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
	
	network2, err := gw.GetNetwork("chanorg2org3")

	if err != nil {
		log.Fatalf("Failed to get network2: %v", err)
	}

	contract2 := network2.GetContract("basic-23")
	
	result, err := contract2.EvaluateTransaction("ReadWood")
	if err != nil {
		log.Println("Failed to evaluate transaction ReadWood: %v", err)
		return err
	}
	strRes := string(result)
	intRes, err := strconv.Atoi(strRes)
	if err != nil {
		log.Println("Errore di conversione: %v", err)
		return err
	}
	if intRes > 0 {
		if intRes > 300 {
			strRes = "300"
		}
		result, err = contract2.SubmitTransaction("TransferWood", strRes, "Org2MSP")
		if err != nil {
			log.Println("Failed to evaluate transaction TransferWood: %v", err)
			return err
		}
		
		log.Println("Refill di Wood effettuato")
	} else {
		log.Println("Org3MSP non ha abbastanza legna")
	}
	
	return nil
}

func populateWallet2(wallet *gateway.Wallet) error {
	
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



func populateWallet(wallet *gateway.Wallet) error {
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



