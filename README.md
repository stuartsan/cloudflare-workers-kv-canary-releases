# Canary releases at the edge with Cloudflare workers

## Setup

- go into CF account in dashboard and for a particular site, enable workers, and enable workers KV
- get CF API key and set CLOUDFLARE_EMAIL and CLOUDFLARE_TOKEN env vars. and CLOUDFLARE_ACCOUNT. see https://developers.cloudflare.com/workers/kv/writing-data/
- also, get zone name and id
- Deploy the "app" `./deploy.sh ce00382cb26f4d93964753ea0014da62 new.html app/new.html` and `./deploy.sh ce00382cb26f4d93964753ea0014da62 old.html app/old.html`
