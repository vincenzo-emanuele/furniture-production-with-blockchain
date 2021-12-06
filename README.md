# progetto-sdd
1. Dalla cartella `fabric-samples/test-network` eseguire `network.sh up` per creare Org1 e Org2
2. SENZA CAMBIARE CARTELLA eseguire il seguente comando: `./scripts/createChannelOrg12.sh` e il comando `./scripts/createChannelOrg23.sh`. In caso di permesso negato fornire i permessi di esecuzione con il comando `chmod +x ./scripts/createChannelOrg12.sh`

Nel file `./scripts/createChannelOrg23.sh` è presente anche il richiamo alle istruzioni per la creazione e l'aggiunta al canale e alla rete dell'Org3 (lo script che viene invocato è `fabric-samples/test-network/addOrg3/addOrg3.sh`. Per settare le variabili d'ambiente ed operare come una determinata organizzazione, lanciare il comando `export $(./setOrgEnv3.sh <orgName>)` dove `<orgName>` indica il nome dell'organizzazione. Per verificare la corretta configurazione dell'infrastruttura è possibile lanciare il comando precedente in combinazione con il comando `peer channel list` per listare tutti i canali a cui si è unita una determinata organizzazione.
Il setup corretto è:
- Org1 collegato a chanorg1org2;
- Org2 collegato a chanorg1org2 e chanorg2org3;
- Org3 collegato a chanorg2org3; 
