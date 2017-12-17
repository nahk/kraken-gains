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

// Wallet holds all the currencies
type Wallet struct {
	BTC Currency
	XRP Currency
}

// Currency wraps a ticker response and holds a balance
type Currency struct {
	Ask, Reference, Potential float32
	Balance                   float32
}

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
	var err error
	var w Wallet

	api := krakenapi.New(os.Getenv("KRAKEN_API_KEY"), os.Getenv("KRAKEN_PRIVATE_KEY"))

	ticker, err := api.Ticker(krakenapi.XXBTZEUR, krakenapi.XXRPZEUR)
	if err != nil {
		log.Fatal(err)
	}

	err = w.getReference()
	if err != nil {
		log.Fatal(err)
	}

	err = w.updateAsks(ticker)
	if err != nil {
		log.Fatal(err)
	}

	err = w.updateBalances(api)
	if err != nil {
		log.Fatal(err)
	}

	err = w.calculatePotential()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(w)

	// Output
	fmt.Printf("                                 \n")
	fmt.Printf(" Potential gain │  Current price \n")
	fmt.Printf("────────────────┼────────────────\n")
	fmt.Printf(" %14.2f │ %14.2f \n", w.BTC.Potential, w.BTC.Ask)
	fmt.Printf("                                 \n")
}

func (w Wallet) getReference() (err error) {
	var r float64
	r, err = strconv.ParseFloat(os.Getenv("BTC_PRICE_REFERENCE"), 32)
	if err != nil {
		return err
	}

	w.BTC.Reference = float32(r)

	return nil
}

func (w Wallet) updateAsks(ticker *krakenapi.TickerResponse) (err error) {
	var a float64

	a, err = strconv.ParseFloat(ticker.XXBTZEUR.Ask[0], 32)
	if err != nil {
		return err
	}
	w.BTC.Ask = float32(a)

	a, err = strconv.ParseFloat(ticker.XXRPZEUR.Ask[0], 32)
	if err != nil {
		return err
	}
	w.XRP.Ask = float32(a)

	return nil
}

func (w Wallet) updateBalances(api *krakenapi.KrakenApi) (err error) {
	var b float64

	b, err = strconv.ParseFloat(os.Getenv("BTC_BALANCE"), 32)
	if err != nil {
		return err
	}
	w.BTC.Balance = float32(b)

	b, err = strconv.ParseFloat(os.Getenv("XRP_BALANCE"), 32)
	if err != nil {
		return err
	}
	w.XRP.Balance = float32(b)

	if w.BTC.Balance == 0 && w.XRP.Balance == 0 {
		b, err := api.Balance()
		if err != nil {
			return err
		}
		w.BTC.Balance = b.XXBT
		w.XRP.Balance = b.XXRP
	}

	log.Println(w)

	return nil
}

func (w Wallet) calculatePotential() (err error) {
	// BTC (base currency)
	w.BTC.Potential = w.BTC.Balance * w.BTC.Ask

	// XRP
	w.XRP.Potential = w.XRP.Balance * w.XRP.Ask

	return nil
}
