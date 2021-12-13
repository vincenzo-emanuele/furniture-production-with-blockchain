## Smart Contract per canale tra Org1 e Org2

Questo Smart Contract permette alle due organizzazioni (Org1 è un negozio e Org2 è una fabbrica) di gestire un proprio wallet costituito da NFT rappresentanti mobili.

È stato realizzato interamente in go e L'NFT è rappresentato da una struct con questa struttura:

```
type Asset struct {
  ID             string
  Price          int    
  Color          string
  Type           string
}
```
Sono state previste le principali operazioni di gestione dei Token e, per rimanere nel constesto dove Org2 è la fabbrica e Org1 è il negozio, abbiamo tolto la
possibilità ad Org1 di richiamare `createAsset` e crearsi da solo i *mobili* senza passare per la fabbrica.

In questo Smart Contract non è previsto il sistema di scambio di valuta, cosa che verrà implementata tramite un altro Smart Contract che risiederà su un canale condiviso da tutte e tre le Org (non necessario)

Inoltre, è stata prevista una gestione degli eventi per permettere ad Org1 di ricevere nuovi asset (ammesso che Org2 abbia disponibilità) ogni volta 
che viene effettuata una vendita e per permettere ad Org2 di consumare le risorse necessarie per la creazione di nuovi asset
