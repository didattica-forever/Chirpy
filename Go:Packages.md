I package in Go sono il modo in cui il linguaggio organizza e gestisce il codice. Funzionano come moduli o librerie che ti permettono di riutilizzare il codice, strutturare i progetti e controllare la visibilità di funzioni e variabili.

-----

### Concetti Fondamentali

  * **Organizzazione del codice**: Un package è una cartella contenente uno o più file `.go`. Tutti i file all'interno di quella cartella devono dichiarare lo stesso nome di package all'inizio del file.
  * **Visibilità**: In Go, la visibilità (pubblico o privato) non è gestita con parole chiave come `public` o `private`. Dipende invece dalla **lettera iniziale** del nome.
      * Un nome (di funzione, variabile, tipo, ecc.) che inizia con una **lettera maiuscola** è **pubblico** ed è visibile (cioè, può essere esportato) da altri package che lo importano.
      * Un nome che inizia con una **lettera minuscola** è **privato** ed è visibile solo all'interno del suo package.
  * **Importazione**: Per usare le funzionalità di un package, devi importarlo usando la parola chiave `import`. Il nome di un package importato corrisponde al suo percorso di importazione. Ad esempio, `import "fmt"` importa il package standard per la formattazione di input/output.

-----

### Anatomia di un Package

Ogni file `.go` inizia con una dichiarazione `package`.

```go
// package main è il punto di ingresso per i programmi eseguibili
package main

// Importazione di package
import (
    "fmt"
    "log"
)

// Funzione pubblica (inizia con lettera maiuscola)
func Saluta() {
    fmt.Println("Ciao!")
}

// Funzione privata (inizia con lettera minuscola)
func salutaInterno() {
    log.Println("Ciao interno!")
}

func main() {
    Saluta()
    salutaInterno()
}
```

### Struttura e Nomenclatura

  * **`package main`**: Questo nome è speciale. Indica che il package è un'applicazione eseguibile. La funzione `main()` al suo interno è il punto di ingresso del programma.
  * **Nome del package**: Il nome di un package, per convenzione, è lo stesso della cartella che lo contiene. Per esempio, un package in una cartella chiamata `utils` dovrebbe chiamarsi `package utils`.
  * **Percorso di importazione**: Per importare un package, usi il suo percorso di importazione. Se hai un package in `github.com/tuo-utente/tuo-progetto/utils`, l'importazione sarà `import "github.com/tuo-utente/tuo-progetto/utils"`.

### Il Ruolo del `go.mod`

Il file `go.mod` gestisce le dipendenze del tuo progetto, specificando il modulo principale e le versioni dei package esterni. Quando importi un package, Go cerca prima nel tuo modulo locale e poi nelle dipendenze elencate in `go.mod`. Questo sistema assicura che il tuo codice sia portabile e riproducibile su qualsiasi macchina.


Sì, la documentazione ufficiale di Go è la fonte più affidabile. C'è un'ottima guida sul sito ufficiale di Go.

-----

### Cos'è `go.mod`?

Il file `go.mod` (Go Module) è il cuore del sistema di gestione delle dipendenze in Go a partire dalla versione 1.11. 📦 Svolge il ruolo di un manifest che definisce il modulo corrente, traccia le sue dipendenze e gestisce le versioni dei pacchetti esterni. Ha essenzialmente lo stesso scopo del `package.json` in Node.js o del `requirements.txt` in Python.

Il vantaggio principale è che **non hai più bisogno di mettere il tuo codice nel `$GOPATH`** e puoi avere più moduli nel tuo computer.

-----

### Come usare `go.mod`

#### 1\. Iniziare un nuovo modulo

Per creare un nuovo progetto Go, il primo passo è inizializzare un nuovo modulo. Apri il terminale nella directory del tuo progetto e digita:

```bash
go mod init <nome-del-modulo>
```

Sostituisci `<nome-del-modulo>` con il nome che vuoi dare al tuo progetto, solitamente il percorso del repository (es. `github.com/tuo-utente/il-mio-progetto`). Questo comando creerà un nuovo file `go.mod` nella directory corrente.

#### 2\. Aggiungere una dipendenza

Quando importi e usi un pacchetto di terze parti nel tuo codice, il comando `go` lo rileverà automaticamente e lo aggiungerà a `go.mod` la prossima volta che eseguirai un comando come `go build`, `go run` o `go test`.

In alternativa, puoi aggiungere una dipendenza manualmente con:

```bash
go get <percorso-pacchetto>
```

Ad esempio, `go get github.com/gin-gonic/gin` aggiungerà il pacchetto Gin al tuo `go.mod` e scaricherà il codice nel tuo cache locale.

#### 3\. Comandi utili

  * `go mod tidy`: Pulisce il tuo `go.mod`, rimuovendo le dipendenze non più usate e aggiungendo quelle che mancano. È una buona pratica eseguirlo regolarmente.
  * `go list -m all`: Mostra tutte le dipendenze del tuo modulo, incluse quelle indirette.
  * `go mod graph`: Visualizza il grafo delle dipendenze, utile per capire le relazioni tra i pacchetti.

#### 4\. File `go.sum`

Oltre a `go.mod`, verrà creato anche un file `go.sum`. Questo file contiene i checksum crittografici (hash) di ogni dipendenza.  È fondamentale per garantire che i pacchetti scaricati non siano stati manomessi, rafforzando la sicurezza del tuo progetto. **Non dovresti mai modificarlo manualmente.**