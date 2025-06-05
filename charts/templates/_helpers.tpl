{{- define "prometheus-whistleblower.name" -}}
prometheus-whistleblower
{{- end }}

{{- define "prometheus-whistleblower.fullname" -}}
{{ .Release.Name }}-{{ include "prometheus-whistleblower.name" . }}
{{- end }}