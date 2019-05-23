package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cloudflare/cloudflare-go"
)

func main() {
	dir := flag.String("dir", "", "Directory whose contents will be deployed to Workers KV")
	deployId := flag.String("deploy-id", "", "ID for deployment (key prefix). Why not use a git SHA?")
	cfAPIKey := flag.String("cf-api-key", "", "CloudFlare API key")
	cfEmail := flag.String("cf-email", "", "CloudFlare email")
	cfAccount := flag.String("cf-account", "", "CloudFlare account id")
	cfNamespace := flag.String("cf-kv-namespace", "APP_DEPLOYS", "CloudFlare workers KV namespace name to write to")

	flag.Parse()

	if *dir == "" {
		log.Fatal("-dir required")
	}

	if *deployId == "" {
		log.Fatal("-deploy-id required")
	}

	api, err := cloudflare.New(*cfAPIKey, *cfEmail, cloudflare.UsingOrganization(*cfAccount))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := api.ListWorkersKVNamespaces(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var nsId string

	for _, ns := range resp.Result {
		if ns.Title == *cfNamespace {
			nsId = ns.ID
			break
		}
	}

	err = os.Chdir(*dir)
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			key := fmt.Sprintf("%s/%s", *deployId, path)

			value, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			_, err = api.WriteWorkersKV(context.Background(), nsId, key, value)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
