## Smart Contract per il canale tra Org2 e Org3

Questo è lo Smart Contract per gestire lo scambio di Legname tra Org2 (fabbrica di mobili) e Org3 (Fornitore di legna), rappresentando la legna come un Token
fungibile incluso nel wallet.

Lo Smart Contract è stato scritto interamente in Go e sono state previste le principali operazioni di gestione dei wallet. Per rimanere nel constesto dei
ruoli, è stato impedito a Org2 di richiamare `Mint` grazie al quale può procurarsi la *legna* senza passare per il fornitore.

Anche qui non è stato previsto il sistema di scambio di denaro, come anche il sistema di consumazione della legna per creare gli NFT da utilizzare nello Smart Contract 1-2 (Probabilmente sarà fatto a livello applicativo se il prof è d'accordo)
