package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	/*
	"strconv"
	"strings"
	*/
//	"github.com/hyperledger/fabric/common/util"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
//Insert struct field in alphabetic order => to achieve determinism accross languages
// golang keeps the order when marshal to json but doesn't order automatically
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

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "asset1", Color: "blue", Type: "table", Price: 300},
		{ID: "asset2", Color: "red", Type: "chair", Price: 400},
		{ID: "asset3", Color: "green", Type: "table", Price: 500},
		{ID: "asset4", Color: "yellow", Type: "chair", Price: 600},
		{ID: "asset5", Color: "black", Type: "table", Price: 700},
		{ID: "asset6", Color: "white", Type: "chair", Price: 800},
	}
	var walletOrg1, walletOrg2 Wallet
	walletOrg1.Owner = "Org1MSP"
	
	walletOrg1.NFT = make(map[string]Asset)
	walletOrg1.NFT[assets[0].ID] = assets[0]
	walletOrg1.NFT[assets[1].ID] = assets[1]
	walletOrg1.NFT[assets[2].ID] = assets[2]
	
	walletOrg2.Owner = "Org2MSP"
	walletOrg2.NFT = make(map[string]Asset)
	walletOrg2.NFT[assets[3].ID] = assets[3]
	walletOrg2.NFT[assets[4].ID] = assets[4]
	walletOrg2.NFT[assets[5].ID] = assets[5]
	
	walletJSON, err := json.Marshal(walletOrg1)
	
	if err != nil {
		return err
	}
	
	err = ctx.GetStub().PutState("Org1MSP", walletJSON)
	if err != nil {
		return fmt.Errorf("Errore nell'aggiunta di Org1 allo stato globale. %v", err)
	}
	
	walletJSON, err = json.Marshal(walletOrg2)
	
	if err != nil {
		return err
	}
	
	err = ctx.GetStub().PutState("Org2MSP", walletJSON)
	if err != nil {
		return fmt.Errorf("Errore nell'aggiunta di Org2 allo stato globale. %v", err)
	}
	
	return nil
}

// CreateAsset puo' essere richiamato solo da Org2 e crea un token, aggiungengolo al wallet di Org2
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, Type string, price int) error {
	
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("Errore nella lettura dell'identita'")
	}
	
	if clientOrgID != "Org2MSP" {
		return fmt.Errorf("L'unico a poter chiamare questa funzione e' Org2MSP")
	}
/*	
	invokeArgs := util.ToChaincodeArgs("ReadWood")

	result := ctx.GetStub().InvokeChaincode("basic-23", invokeArgs, "chanorg2org3")
	
	var str = result.String()

    	var split = strings.Split(str, "status:200 payload:\"")
    	if len(split) < 2 {
    		return fmt.Errorf("Errore lancio query: " + result.String())
    	}
    	var split2 = strings.TrimRight(split[1], "\" ")
    	intRes, err := (strconv.Atoi(split2))
    	if err != nil {
    		return fmt.Errorf("Errore di conversione: %v", err)
    	}
	if intRes < 100 {
		return fmt.Errorf("quantita' di legna insufficiente: " + strconv.Itoa(intRes))
	}
*/	
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("L'asset con ID: %s esiste gia'", id)
	}

	asset := Asset{
		ID:             id,
		Color:          color,
		Type:           Type,
		Price:		price,
	}
	
	walletJSON, err := ctx.GetStub().GetState("Org2MSP")
	
	if err != nil {
		return fmt.Errorf("Errore nella lettura dello stato globale: %v", err)
	}
	
	var wallet Wallet
	err = json.Unmarshal(walletJSON, &wallet)
	wallet.NFT[asset.ID] = asset
	
	walletJSON, err = json.Marshal(wallet)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("Org2MSP", walletJSON)
	
	if err != nil {
		return fmt.Errorf("Errore nell'aggiunta di Org2 allo stato globale. %v", err)
	}
	
	payload, err := json.Marshal(asset)
	err = ctx.GetStub().SetEvent("CreateAsset", payload)
	if err != nil {
		return fmt.Errorf("Errore nel settaggio dell'evento. %v", err)
	}
	
	return nil
}


// ReadAsset ritorna l'asset con l'ID specificato se contenuto all'interno del proprio wallet
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (Asset, error) {
	
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	var asset Asset
	if err != nil {
		return asset, fmt.Errorf("failed to read from world state: %v", err)
	}
	
	
	walletJSON, err := ctx.GetStub().GetState(clientOrgID)
	
	if walletJSON == nil {
		return asset, fmt.Errorf("the wallet %s does not exist", id)
	}

	var wallet Wallet
	
	err = json.Unmarshal(walletJSON, &wallet)
	if err != nil {
		return asset, err
	}
	
	_, presente := wallet.NFT[id]
	
	if presente {
		return wallet.NFT[id], nil
	}

	return asset, fmt.Errorf("ID non trovato")
}


// UpdateAsset aggiorna le informazioni di un token contenuto nel propio wallet
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, Type string, price int) error {
	
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	
	walletJSON, err := ctx.GetStub().GetState(clientOrgID)
	
	var wallet Wallet
	err = json.Unmarshal(walletJSON, &wallet)
	
	if err != nil {
		return fmt.Errorf("Errore nella lettura del wallet")
	}
	
	_ , presente := wallet.NFT[id]
	
	if ! presente {
		return fmt.Errorf("ID non trovato")
	}

	// Creo nuovo token
	asset := Asset{
		ID:             id,
		Color:          color,
		Type:           Type,
		Price:		price,
	}
	
	wallet.NFT[id] = asset
	walletJSON, err = json.Marshal(wallet)
	
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(clientOrgID, walletJSON)
}

// DeleteAsset cancella un token dal proprio wallet
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	
	walletJSON, err := ctx.GetStub().GetState(clientOrgID)
	
	if err != nil {
		return fmt.Errorf("Errore nella lettura dello stato globale: %v", err)
	}
	
	var wallet Wallet
	err = json.Unmarshal(walletJSON, &wallet)
	
	_ , presente := wallet.NFT[id]
	
	if ! presente {
		return fmt.Errorf("L'ID non esiste");
	}
	
	delete(wallet.NFT, id)
	
	walletJSON, err = json.Marshal(wallet)
	
	if err != nil {
		fmt.Errorf("Errore nel marshalling")
	}
	
	payload, err := json.Marshal(len(wallet.NFT))
	err = ctx.GetStub().SetEvent("DeleteAsset", payload)
	if err != nil {
		return fmt.Errorf("Errore nel settaggio dell'evento. %v", err)
	}

	return ctx.GetStub().PutState(clientOrgID, walletJSON) 
}

// AssetExists controlla l'esistenza dell'ID specificato all'interno di tutti i wallet
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	
	org := "Org2MSP"
	walletJSON, err := ctx.GetStub().GetState(org)
	
	if err != nil {
		return false, fmt.Errorf("Errore nella lettura dello stato globale: %v", err)
	}
	
	var wallet Wallet
	err = json.Unmarshal(walletJSON, &wallet)
	
	_ , presente := wallet.NFT[id]
	
	if presente {
		return true, nil;
	}
	
	org = "Org1MSP"
	walletJSON, err = ctx.GetStub().GetState(org)
	
	if err != nil {
		return false, fmt.Errorf("Errore nella lettura dello stato globale: %v", err)
	}
	
	err = json.Unmarshal(walletJSON, &wallet)
	
	if err != nil {
		return false, fmt.Errorf("Errore nell'unmarshalling: %v", err)
	}
	
	_ , presente = wallet.NFT[id]
	
	return presente, nil
}





// TransferAsset trasferisce un token da un wallet ad un altro
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	
	
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	
	if err != nil {
		return err
	}
	
	if (newOwner == "Org1MSP" || newOwner == "Org2MSP") && newOwner != clientOrgID {
	 
		var wallet1, wallet2 Wallet
		
		walletJSON1, err := ctx.GetStub().GetState("Org1MSP")
		walletJSON2, err := ctx.GetStub().GetState("Org2MSP")
		
		err = json.Unmarshal(walletJSON1, &wallet1)
		if err != nil {
			return fmt.Errorf("Errore nell'unmarshalling: %v", err)
		}
		
		err = json.Unmarshal(walletJSON2, &wallet2)
		if err != nil {
			return fmt.Errorf("Errore nell'unmarshalling: %v", err)
		}
		
		if clientOrgID == "Org1MSP" {
			_, presente := wallet1.NFT[id]
			
			if ! presente {
				return fmt.Errorf("ID non trovato")
			}
			
			wallet2.NFT[id] = wallet1.NFT[id]
			delete(wallet1.NFT, id)
			
			
			walletJSON1, err = json.Marshal(wallet1)
			if err != nil {
				fmt.Errorf("Errore nel marshalling")
			}
			
			walletJSON2, err = json.Marshal(wallet2)
			if err != nil {
				fmt.Errorf("Errore nel marshalling")
			}
			
			err = ctx.GetStub().PutState("Org1MSP", walletJSON1)
			if err != nil {
				fmt.Errorf("Errore nella scrittura di wallet1 nello stato globale")
			}
			
			err = ctx.GetStub().PutState("Org2MSP", walletJSON2)
			if err != nil {
				fmt.Errorf("Errore nella scrittura di wallet2 nello stato globale")
			}
			
			return nil
			
		} else if clientOrgID == "Org2MSP" {
			
			_, presente := wallet2.NFT[id]
			
			if ! presente {
				return fmt.Errorf("ID non trovato")
			}
			
			wallet1.NFT[id] = wallet2.NFT[id]
			delete(wallet2.NFT, id)
			
			
			walletJSON1, err = json.Marshal(wallet1)
			if err != nil {
				fmt.Errorf("Errore nel marshalling")
			}
			
			walletJSON2, err = json.Marshal(wallet2)
			if err != nil {
				fmt.Errorf("Errore nel marshalling")
			}
			
			err = ctx.GetStub().PutState("Org1MSP", walletJSON1)
			if err != nil {
				fmt.Errorf("Errore nella scrittura di wallet1 nello stato globale")
			}
			
			err = ctx.GetStub().PutState("Org2MSP", walletJSON2)
			if err != nil {
				fmt.Errorf("Errore nella scrittura di wallet2 nello stato globale")
			}
			
			return nil
		} else {
			return fmt.Errorf("Attenzione, controlla chi riceve cosa")
		}	
	} else {
		return fmt.Errorf("Attenzione, controlla chi riceve cosa")
	}
}







// GetAllAssets ritorna tutti i wallet nello stato globale
func (s *SmartContract) GetAllWallets(ctx contractapi.TransactionContextInterface) ([]*Wallet, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var wallets []*Wallet
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var wallet Wallet
		err = json.Unmarshal(queryResponse.Value, &wallet)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, &wallet)
	}

	return wallets, nil
}
