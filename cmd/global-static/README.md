## Set up nginx-ingress to handle requests from LoadBalancer

1. `helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx`

1. `helm repo update`

1. Install nginx controller to kube `helm install ingress-nginx ingress-nginx/ingress-nginx`


## Create subdomain entries

1. `kubectl -n default get services -o wide ingress-nginx-controller`

1. Create A records pointing to external domain

1. Test that IP has propogated: `dig +short blog.example.com`

1. finish TLS setup at https://www.linode.com/docs/guides/how-to-configure-load-balancing-with-tls-encryption-on-a-kubernetes-cluster/



## How to setup TLS key

1. Generate a TLS key and certificate using a TLS toolkit like OpenSSL. Be sure to change the CN and O values to those of your own website domain.

```
 openssl req -newkey rsa:4096 \     
   -x509 \     
   -sha256 \     
   -days 3650 \     
   -nodes \     
   -out tls.crt \     
   -keyout tls.key \     
   -subj "/CN=mywebsite.com/O=mywebsite.com"
```

1. Create the secret using the create secret tls command. Ensure you substitute $SECRET_NAME for the name you’d like to give to your secret. This will be how you reference the secret in your Service manifest.

`kubectl create secret tls $SECRET_NAME --cert tls.crt --key tls.key`

1. You can check to make sure your Secret has been successfully stored by using describe:

`kubectl describe secret $SECRET_NAME`
