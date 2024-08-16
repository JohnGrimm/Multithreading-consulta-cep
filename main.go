package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	figure "github.com/common-nighthawk/go-figure"
)

const (
	title   = "Consulta CEP"
	version = "1.0"
)

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Err         string `json:"erro,omitempty"`
}

type Result struct {
	API     string
	Address interface{}
	Err     error
}

func main() {
	// Setup for graceful exit on CTRL+C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		fmt.Println("\nExiting... Bye!")
		os.Exit(0)
	}()

	myFigure := figure.NewColorFigure(title+" V"+version, "doom", "cyan", false)
	myFigure.Print()
	myFigure = figure.NewColorFigure("by John Grimm", "bubble", "green", false)
	myFigure.Print()
	fmt.Println("")

	for {
		var cep string
		fmt.Print("Enter CEP (or press CTRL+C to exit): ")
		fmt.Scanln(&cep)
		cep = strings.TrimSpace(cep)

		// Check if input is valid
		if len(cep) != 8 {
			fmt.Println("Invalid CEP. Please enter an 8-digit CEP.")
			continue
		}

		timeout := time.Second * 1
		ch := make(chan Result)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		go fetchFromAPI(ctx, ch, "https://brasilapi.com.br/api/cep/v1/"+cep, "BrasilAPI")
		go fetchFromAPI(ctx, ch, "http://viacep.com.br/ws/"+cep+"/json/", "ViaCEP")

		select {
		case res := <-ch:
			if res.Err != nil {
				fmt.Println("Error:", res.Err)
			} else {
				addressJson, _ := json.MarshalIndent(res.Address, "", "  ")
				fmt.Printf("API: %s\nResult:\n%s\n", res.API, string(addressJson))
			}
		case <-ctx.Done():
			fmt.Println("Error: Timeout")
		}
	}
}

func fetchFromAPI(ctx context.Context, ch chan<- Result, url string, apiName string) {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		ch <- Result{API: apiName, Err: err}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- Result{API: apiName, Err: fmt.Errorf("API %s returned status code %d", apiName, resp.StatusCode)}
		return
	}

	if apiName == "ViaCEP" {
		var address ViaCEP
		if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
			ch <- Result{API: apiName, Err: err}
			return
		}

		if address.Err == "true" {
			ch <- Result{API: apiName, Err: fmt.Errorf("API %s returned: cep not found", apiName)}
			return
		}

		ch <- Result{API: apiName, Address: address}
	}

	var address BrasilAPI
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		ch <- Result{API: apiName, Err: err}
		return
	}

	ch <- Result{API: apiName, Address: address}

}
