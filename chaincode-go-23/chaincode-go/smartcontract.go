package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}


//Realizzo un wallet di Token NFT
type Wallet struct{
	Owner string	`json:"Owner"`
	Wood uint32	`json:"Wood"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	
	var wallet2, wallet3 Wallet
	wallet2.Owner = "Org2MSP"
	wallet2.Wood = 100
	
	wallet3.Owner = "Org3MSP"
	wallet3.Wood = 100
	
	walletJSON, err := json.Marshal(wallet2)
	if err != nil {
		return err
	}
	
	err = ctx.GetStub().PutState("Org2MSP", walletJSON)
	if err != nil {
		return fmt.Errorf("Errore nell'aggiunta di Org1 allo stato globale. %v", err)
	}
	
	walletJSON, err = json.Marshal(wallet3)
	
	if err != nil {
		return err
	}
	
	err = ctx.GetStub().PutState("Org3MSP", walletJSON)
	if err != nil {
		return fmt.Errorf("Errore nell'aggiunta di Org2 allo stato globale. %v", err)
	}
	
	return nil
}

// CreateAsset puo' essere richiamato solo da Org2 e crea un token, aggiungengolo al wallet di Org2
func (s *SmartContract) Mint(ctx contractapi.TransactionContextInterface, amount uint32) error {
	
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("Errore nella lettura dell'identita'")
	}
	
	if clientOrgID != "Org3MSP" {
		return fmt.Errorf("L'unico a poter chiamare questa funzione e' Org3MSP")
	}
	
	walletJSON, err := ctx.GetStub().GetState("Org3MSP")
	
	if err != nil {
		return fmt.Errorf("Errore nella lettura dello stato globale: %v", err)
	}
	
	var wallet Wallet
	err = json.Unmarshal(walletJSON, &wallet)
	wallet.Wood += amount
	
	walletJSON, err = json.Marshal(wallet)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState("Org3MSP", walletJSON)
	
	if err != nil {
		return fmt.Errorf("Errore nell'aggiunta di Org3 allo stato globale. %v", err)
	}
	
	return nil
}


// ReadAsset ritorna l'asset con l'ID specificato se contenuto all'interno del proprio wallet
func (s *SmartContract) ReadWood(ctx contractapi.TransactionContextInterface) (uint32, error) {
	
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}
	
	
	walletJSON, err := ctx.GetStub().GetState(clientOrgID)
	if walletJSON == nil {
		return 0, fmt.Errorf("the wallet %s does not exist", clientOrgID)
	}

	var wallet Wallet
	
	err = json.Unmarshal(walletJSON, &wallet)
	if err != nil {
		return 0, err
	}

	return wallet.Wood, nil
}




// DeleteAsset cancella un token dal proprio wallet
func (s *SmartContract) RemoveWood(ctx contractapi.TransactionContextInterface, amount uint32 ) error {
	
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
	
	if amount > wallet.Wood {
		return fmt.Errorf("Non disponi di quella quantita'")
	}
	
	wallet.Wood -= amount
	
	walletJSON, err = json.Marshal(wallet)
	
	if err != nil {
		return fmt.Errorf("Errore nel marshalling")
	}

	return ctx.GetStub().PutState(clientOrgID, walletJSON) 
}



// TransferAsset trasferisce un token da un wallet ad un altro
func (s *SmartContract) TransferWood(ctx contractapi.TransactionContextInterface, amount uint32, newOwner string) error {
	
	
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	
	if err != nil {
		return err
	}
	
	if (newOwner == "Org2MSP" || newOwner == "Org3MSP") && newOwner != clientOrgID {
	 
		var wallet2, wallet3 Wallet
		
		walletJSON2, err := ctx.GetStub().GetState("Org2MSP")
		walletJSON3, err := ctx.GetStub().GetState("Org3MSP")
		
		err = json.Unmarshal(walletJSON2, &wallet2)
		if err != nil {
			return fmt.Errorf("Errore nell'unmarshalling: %v", err)
		}
		
		err = json.Unmarshal(walletJSON3, &wallet3)
		if err != nil {
			return fmt.Errorf("Errore nell'unmarshalling: %v", err)
		}
		
		if clientOrgID == "Org2MSP" {
						
			if amount > wallet2.Wood {
				return fmt.Errorf("Non disponi di questa quantita'")
			}
			
			wallet2.Wood -= amount
			wallet3.Wood += amount
			
			walletJSON2, err = json.Marshal(wallet2)
			if err != nil {
				fmt.Errorf("Errore nel marshalling")
			}
			
			walletJSON3, err = json.Marshal(wallet3)
			if err != nil {
				fmt.Errorf("Errore nel marshalling")
			}
			
			err = ctx.GetStub().PutState("Org2MSP", walletJSON2)
			if err != nil {
				fmt.Errorf("Errore nella scrittura di wallet2 nello stato globale")
			}
			
			err = ctx.GetStub().PutState("Org3MSP", walletJSON3)
			if err != nil {
				fmt.Errorf("Errore nella scrittura di wallet3 nello stato globale")
			}
			
			return nil
			
		} else if clientOrgID == "Org3MSP" {
			
			if amount > wallet3.Wood {
				return fmt.Errorf("Non disponi di questa quantita'") 
			}
			
			wallet3.Wood -= amount
			wallet2.Wood += amount
			
			walletJSON2, err = json.Marshal(wallet2)
			if err != nil {
				return fmt.Errorf("Errore nel marshalling")
			}
			
			walletJSON3, err = json.Marshal(wallet3)
			if err != nil {
				return fmt.Errorf("Errore nel marshalling")
			}
			
			err = ctx.GetStub().PutState("Org2MSP", walletJSON2)
			if err != nil {
				return fmt.Errorf("Errore nella scrittura di wallet2 nello stato globale")
			}
			
			err = ctx.GetStub().PutState("Org3MSP", walletJSON3)
			if err != nil {
				return fmt.Errorf("Errore nella scrittura di wallet3 nello stato globale")
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
