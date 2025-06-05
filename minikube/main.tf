terraform {
    required_providers {
        kubernetes = {
            source = "hashicorp/kubernetes"
             version = "2.37.1"
        }
        helm = {
            source = "hashicorp/helm"
            version = "2.17.0"
        }
    }
}

provider "helm" {
    kubernetes {
        config_path = "~/.kube/config"
        config_context = "minikube"
    }
}

provider "kubernetes" {
    config_path = "~/.kube/config"
    config_context = "minikube"
}

