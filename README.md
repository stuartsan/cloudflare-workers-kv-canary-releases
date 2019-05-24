# Canary releases at the edge with Cloudflare workers

This project demonstrates A) how to continuously deploy a static site directly to Cloudflare Workers KV store at edge locations and B) how to do canary releases (AKA staged rollouts AKA phased rollouts) at the edge.

The worker script defines the stages: whether they're based on ip address, geolocation, user segments, random samples / increasing percentages and cookie setting to make it sticky, whatever.

## Requirements
- serverless framework
- CloudFlare account
- golang

## Setup
The first step is to go into your Cloudflare account, visit the dashboard for the site on which you want to deploy the worker, enable workers, and enable workers KV.

Locate each of the following items and set environment variables for them on your local machine:
- API key => CLOUDFLARE_AUTH_KEY
- Email => CLOUDFLARE_AUTH_EMAIL
- Account ID => CLOUDFLARE_ACCOUNT
- Zone => CLOUDFLARE_ZONE

Deploy the worker, create the KV namespaces, and bind them to the worker: `serverless deploy`

### Deploy app
Deploy both the `current` and `next` versions so we have a "new" thing to release:

`go run deploy.go -dir app/current/ -deploy-id current`

`go run deploy.go -dir app/next/ -deploy-id next`

`current` and `next` are ids for demonstration purposes. Probably use a git commit SHA instead for deploy ids!

### Doing canary releases 
   
To send all production traffic to the "current" version of the app:
```
go run release.go -deploy-id current
```

And then do a canary release for next version, releasing it to stage 1:
```
go run release.go -deploy-id next -stage 1
```

And then releasing it to stage 2:
```
go run release.go -deploy-id next -stage 2
```

And then all production traffic to "next":
```
go run release.go -deploy-id next
```
