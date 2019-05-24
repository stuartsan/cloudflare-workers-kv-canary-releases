package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"

	"github.com/cloudflare/cloudflare-go"
)

type ReleaseState struct {
	Current *string `json:"current"`
	Next    *string `json:"next"`
	Stage   *string `json:"stage"`
}

func main() {
	deployId := flag.String("deploy-id", "", "ID for deployment (key prefix). Why not use a git SHA?")
	stage := flag.String("stage", "", "Release stage for deploy-id")
	cfNamespace := flag.String("cf-kv-namespace", "RELEASE_STATE", "CloudFlare workers KV namespace name storing the release state")
	cfAPIKey := flag.String("cf-api-key", os.Getenv("CLOUDFLARE_AUTH_KEY"), "CloudFlare API key")
	cfEmail := flag.String("cf-email", os.Getenv("CLOUDFLARE_AUTH_EMAIL"), "CloudFlare email")
	cfAccount := flag.String("cf-account", os.Getenv("CLOUDFLARE_ACCOUNT"), "CloudFlare account id")

	flag.Parse()

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

	var state *ReleaseState

	// this is NOT a staged release, just ship it!
	if *stage == "" {
		state = &ReleaseState{
			Current: deployId,
			Next:    nil,
			Stage:   nil,
		}
		// here we need to fetch existing state first!
	} else {
		bytes, err := api.ReadWorkersKV(context.Background(), nsId, "state")
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(bytes, &state)
		if err != nil {
			log.Fatal(err)
		}
		state.Next = deployId
		state.Stage = stage
	}

	releaseState, err := json.Marshal(state)
	if err != nil {
		log.Fatal(err)
	}

	_, err = api.WriteWorkersKV(context.Background(), nsId, "state", []byte(string(releaseState)))
	if err != nil {
		log.Fatal(err)
	}
}
