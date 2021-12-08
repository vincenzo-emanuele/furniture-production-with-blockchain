# progetto-sdd
1. Dalla cartella `fabric-samples/test-network` eseguire `./network.sh up` per creare Org1 e Org2
2. SENZA CAMBIARE CARTELLA eseguire il seguente comando: `./scripts/createChannelOrg12.sh` e il comando `./scripts/createChannelOrg23.sh`. In caso di permesso negato fornire i permessi di esecuzione con il comando `chmod +x ./scripts/createChannelOrg12.sh`
3. Eseguire il comando `./network.sh deployCC -c chanorg1org2 -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go` e il comando `./network.sh deployCC23 -c chanorg2org3 -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go` per deployare i ChainCode sul canale chanorg1org2 e sul canale chanorg2org3
4. Nella cartella `fabric-samples/config` eseguire il comando `export FABRIC_CFG_PATH=$PWD`
5. Nella cartella `fabric-samples/test-network` eseguire `export $(./setOrgEnv.sh <orgName>)` per effettuare le chiamate ai ChainCode come `<orgName>`
6. Eseguire il comando `./ccinteract.sh invoke12 chanorg1org2 basic '{"function":"InitLedger","Args":[]}'` per eseguire la funzione di inizializzazione dello ChainCode
7. Eseguire il comando `./ccinteract.sh query chanorg1org2 basic '{"Args":["GetAllWallets"]}'` per lanciare la query ed ottenere tutti i Wallet
8. Eseguire il comando `peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C chanorg1org2 -n basic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"CreateAsset","Args": ["id", "lmao", "asd", "123"]}'` per creare un nuovo asset sul canale chanorg12. 
10. Eseguire il comando `peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C chanorg2org3 -n basic --peerAddresses localhost:11051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"CreateAsset","Args": ["id", "lmao", "asd", "123"]}'` per creare un nuovo asset sul canale chanorg23


Nel file `./scripts/createChannelOrg23.sh` è presente anche il richiamo alle istruzioni per la creazione e l'aggiunta al canale e alla rete dell'Org3 (lo script che viene invocato è `fabric-samples/test-network/addOrg3/addOrg3.sh`. Per settare le variabili d'ambiente ed operare come una determinata organizzazione, lanciare il comando `export $(./setOrgEnv3.sh <orgName>)` dove `<orgName>` indica il nome dell'organizzazione. Per verificare la corretta configurazione dell'infrastruttura è possibile lanciare il comando precedente in combinazione con il comando `peer channel list` per listare tutti i canali a cui si è unita una determinata organizzazione.
Il setup corretto è:
- Org1 collegato a chanorg1org2;
- Org2 collegato a chanorg1org2 e chanorg2org3;
- Org3 collegato a chanorg2org3; 

UPDATE:
Al punto 2 è possibile eseguire il comando `./network.sh createChannel12 -c chanorg1org2` per creare il canale `chanorg1org2` e il comando `./network.sh createChannel23 -c chanorg2org3` per creare il canale `chanorg2org3`.

