# Minikube

This directory allows you to deploy the configured Prometheus Whistleblower to a Minikube cluster
alongside the kube-prometheus-stack helm chart.

## Usage

```bash
$ minikube start
$ terraform init
$ terraform apply
# You can now port-foward the prometheus-whistleblower or configure the ingress to work locally
# and access the UI from the container. 
```

