configs to deploy for prod:

1. global/global-config-prod.yaml
2. uf_\*/kube-static/\*.yaml
3. uf_\*/postgres-deployment.yaml
4. uf_\*/deploy-latest.yaml

DO NOT redeploy the load balancer in uf_public/kube-static. It has custom settings in linode
