package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/beldur/kraken-go-api-client"
	"github.com/joho/godotenv"
	"github.com/kardianos/osext"
)

func init() {
	directory, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load(filepath.Join(directory, ".env"))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	api := krakenapi.New(os.Getenv("KRAKEN_API_KEY"), os.Getenv("KRAKEN_PRIVATE_KEY"))
	ticker, err := api.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		log.Fatal(err)
	}

	// Reference price
	r, err := strconv.ParseFloat(os.Getenv("BTC_PRICE_REFERENCE"), 32)
	if err != nil {
		log.Fatal(err)
	}
	reference := float32(r)

	// Ask price
	a, err := strconv.ParseFloat(ticker.XXBTZEUR.Ask[0], 32)
	if err != nil {
		log.Fatal(err)
	}
	ask := float32(a)

	// Current balance
	b, err := strconv.ParseFloat(os.Getenv("BTC_BALANCE"), 32)
	if err != nil {
		log.Fatal(err)
	}
	balance := float32(b)
	if balance == 0 {
		b, err := api.Balance()
		if err != nil {
			log.Fatal(err)
		}
		balance = b.XXBT
	}

	// Output
	potential := balance * (ask - reference)
	fmt.Printf("                                 \n")
	fmt.Printf(" Potential gain │  Current price \n")
	fmt.Printf("────────────────┼────────────────\n")
	fmt.Printf(" %14.2f │ %14.2f \n", potential, ask)
	fmt.Printf("                                 \n")
}
