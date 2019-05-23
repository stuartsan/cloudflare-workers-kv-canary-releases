# Canary releases at the edge with Cloudflare workers

This project demonstrates A) how to continuously deploy a static site directly to Cloudflare Workers KV store at edge locations and B) how to do canary releases (AKA staged rollouts AKA phased rollouts) at the edge.

The worker script defines the stages: whether they're based on ip address, geolocation, user segments, random samples / increasing percentages and cookie setting to make it sticky, whatever.

## Requirements
- serverless framework
- CloudFlare account
- golang

## Setup
- go into Cloudflare account and in the dashboard for a particular site, enable workers, and enable workers KV
- find your CF API key and set `CLOUDFLARE_AUTH_KEY` and `CLOUDFLARE_AUTH_EMAIL` environment variables to appease serverless. And `CLOUDFLARE_ACCOUNT`.
- deploy the worker, create the KV namespaces, and bind them to the worker: `serverless deploy`

### Deploy app
Deploy both the `current` and `next` versions:

`go run deploy.go -dir app/current/ -cf-account $CLOUDFLARE_ACCOUNT -cf-api-key $CLOUDFLARE_AUTH_KEY -cf-email $CLOUDFLARE_AUTH_EMAIL -cf-kv-namespace APP_DEPLOYS -deploy-id current`

`go run deploy.go -dir app/next/ -cf-account $CLOUDFLARE_ACCOUNT -cf-api-key $CLOUDFLARE_AUTH_KEY -cf-email $CLOUDFLARE_AUTH_EMAIL -cf-kv-namespace APP_DEPLOYS -deploy-id next`

`current` and `next` are ids for demonstration purposes. Probably use a git commit SHA instead!

### Doing canary releases 
   
To send all production traffic to the "current" version of the app:
```
go run release.go -cf-account $CLOUDFLARE_ACCOUNT -cf-api-key $CLOUDFLARE_AUTH_KEY -cf-email $CLOUDFLARE_AUTH_EMAIL -cf-kv-namespace RELEASE_STATE -deploy-id current
```

And then do a canary release for next version, releasing it to stage 1:
```
go run release.go -cf-account $CLOUDFLARE_ACCOUNT -cf-api-key $CLOUDFLARE_AUTH_KEY -cf-email $CLOUDFLARE_AUTH_EMAIL -cf-kv-namespace RELEASE_STATE -deploy-id next -stage 1
```

And then releasing it to stage 2:
```
go run release.go -cf-account $CLOUDFLARE_ACCOUNT -cf-api-key $CLOUDFLARE_AUTH_KEY -cf-email $CLOUDFLARE_AUTH_EMAIL -cf-kv-namespace RELEASE_STATE -deploy-id next -stage 2
```

And then all production traffic to "next":
```
go run release.go -cf-account $CLOUDFLARE_ACCOUNT -cf-api-key $CLOUDFLARE_AUTH_KEY -cf-email $CLOUDFLARE_AUTH_EMAIL -cf-kv-namespace RELEASE_STATE -deploy-id next
```
