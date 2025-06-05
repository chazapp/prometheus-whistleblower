resource "kubernetes_namespace" "tools_namespace" {
  metadata {
    name = "monitoring"
  }
}

resource "helm_release" "kube-prometheus-stack" {
  name       = "kube-prometheus-stack"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"
  version    = "73.2.0"
  
  namespace  = "monitoring"

  values = [
    "${file("${path.module}/configs/kube-prometheus-stack.yaml")}"
  ]
}

resource "helm_release" "prometheus-whistleblower" {
  name       = "prometheus-whistleblower"
  repository = "oci://ghcr.io/chazapp/helm-charts/"
  chart      = "prometheus-whistleblower"
  version    = "0.1.0"
  
  namespace  = "monitoring"

  values = [
    "${file("${path.module}/configs/prometheus-whistleblower.yaml")}"
  ]
  depends_on = [ helm_release.kube-prometheus-stack ]
}

