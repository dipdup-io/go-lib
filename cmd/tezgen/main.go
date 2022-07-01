package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/dipdup-net/go-lib/tzkt/api"
	"github.com/dipdup-net/go-lib/tzkt/data"
	"github.com/go-playground/validator/v10"
)

type args struct {
	Contract string `validate:"required_with=URL,len=36"`
	URL      string `validate:"url,required_with=Contract"`
	Name     string `validate:"omitempty"`
	Output   string `validate:"omitempty,dir"`
	File     string `validate:"required_without=URL"`
}

func main() {
	var args args

	flag.StringVar(&args.Contract, "c", "", "Contract address `KT1...`")
	flag.StringVar(&args.URL, "u", "https://api.tzkt.io/", "TzKT base URL")
	flag.StringVar(&args.Name, "n", "my_contract", "Contract name")
	flag.StringVar(&args.Output, "o", "", "Output directory")
	flag.StringVar(&args.File, "f", "", "Path to JSON schema file")

	flag.Parse()

	if err := validator.New().Struct(args); err != nil {
		log.Panic(err)
	}

	var response data.ContractJSONSchema
	var err error

	if args.File == "" {

		tzkt := api.New(args.URL)
		response, err = tzkt.GetContractJSONSchema(context.Background(), args.Contract)
		if err != nil {
			log.Panic(err)
		}

	} else {

		file, err := os.Open(args.File)
		if err != nil {
			log.Panic(err)
		}
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&response); err != nil {
			log.Panic(err)
		}
	}

	if err := Generate(response, args.Name, args.Contract, args.Output); err != nil {
		log.Panic(err)
	}
}
