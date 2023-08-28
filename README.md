for live reload install air
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
go mod init main

generated certificates are stored at:

/etc/letsencrypt/live

TODO:
* we should begin a search with a user's previous buying history (ask: same kind
as last kind?)
* in uf_order/orderFunctions.go, we should figure out all orderObject.Code values
* should 'quit' be allowed to exit any part of fridayy? Should we have AimQuit?

notify if huge price change (30%)

what phone number to use for billing address
at sign up ask for safe word
Link GTP3 to Zinc
Picture sent with text - not sure if already done?
1 - send details before picture
2 - we send a safe word to user to confirm
Address
CC Info
"Front end" text flow
View messages front end
Train Model (sizing, too) 
Landing page linked to initial text message (free or paid plan) 


Phase 2:
Scraper for Shopify Stores + checkout bot
Add other stores
Holiday shopping manager


TO DISCUSS:
we should probably execute zinc before we charge the user. However, this puts us
at risk of the user's payment method failing. There is a chance that Stripe
allows us to preauthorize purchases. If we charge a user on Stripe then the
Zinc order fails, then I think we need to pay a fine to Stripe to reverse the
customer payment


places autocomplete


clothing:
male/female/child/unisex questions

user UI pages for more info


Common Problems:

* error: error parsing global-config-prod.yaml: error converting YAML to JSON: 
yaml: line 14: found character that cannot start any token" means that the line
starts with a tab; it should start with spaces

## Microservices

### To create a new microservice
1. add ports and domains to config.go, including database if needed
1. add the same info to uno/cmd/global/global-config-prod.yaml
1. create main.go in correct folder in uno/cmd
1. create uservice database folder and files if needed
1. create uservice repository on dockerhub
1. add uf_{name} to build-and-push-all-latest.sh, check-uf.sh, 
reapply-all-latest.sh

### To launch a new microservice
1. a source file is required at /fridayy/taxes/full.csv
1. ./build-and-push-all-latest.sh
1. export KUBECONFIG=/fridayy/fridayy-ga-kubeconfig.yaml
1. cd into unc/cmd/{name}/kube-static then run kubectl apply -f for each file
1. kubectl apply -f postgres-deployment.yaml if needed
1. execute database/kubeDbInit.sh if necessary
1. kubectl apply -f prod-latest.yaml OR run reapply-all-latest.sh


