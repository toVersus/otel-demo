{{/*
Expand the name of the chart.
*/}}
{{- define "otel-demo.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "otel-demo.frontend.name" -}}
{{- printf "%s-frontend" (include "otel-demo.name" .) }}
{{- end }}

{{- define "otel-demo.order.name" -}}
{{- printf "%s-order" (include "otel-demo.name" .) }}
{{- end }}

{{- define "otel-demo.payment.name" -}}
{{- printf "%s-payment" (include "otel-demo.name" .) }}
{{- end }}

{{- define "otel-demo.users.name" -}}
{{- printf "%s-users" (include "otel-demo.name" .) }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "otel-demo.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{- define "otel-demo.frontend.fullname" -}}
{{- printf "%s-frontend" (include "otel-demo.fullname" .) }}
{{- end }}

{{- define "otel-demo.order.fullname" -}}
{{- printf "%s-order" (include "otel-demo.fullname" .) }}
{{- end }}

{{- define "otel-demo.payment.fullname" -}}
{{- printf "%s-payment" (include "otel-demo.fullname" .) }}
{{- end }}

{{- define "otel-demo.users.fullname" -}}
{{- printf "%s-users" (include "otel-demo.fullname" .) }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "otel-demo.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "otel-demo.frontend.labels" -}}
helm.sh/chart: {{ include "otel-demo.chart" . }}
{{ include "otel-demo.frontend.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "otel-demo.order.labels" -}}
helm.sh/chart: {{ include "otel-demo.chart" . }}
{{ include "otel-demo.order.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "otel-demo.payment.labels" -}}
helm.sh/chart: {{ include "otel-demo.chart" . }}
{{ include "otel-demo.payment.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "otel-demo.users.labels" -}}
helm.sh/chart: {{ include "otel-demo.chart" . }}
{{ include "otel-demo.users.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "otel-demo.frontend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "otel-demo.frontend.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "otel-demo.order.selectorLabels" -}}
app.kubernetes.io/name: {{ include "otel-demo.order.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "otel-demo.payment.selectorLabels" -}}
app.kubernetes.io/name: {{ include "otel-demo.payment.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "otel-demo.users.selectorLabels" -}}
app.kubernetes.io/name: {{ include "otel-demo.users.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "otel-demo.frontend.serviceAccountName" -}}
{{- if .Values.frontend.serviceAccount.create }}
{{- default (include "otel-demo.frontend.fullname" .) .Values.frontend.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.frontend.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "otel-demo.order.serviceAccountName" -}}
{{- if .Values.order.serviceAccount.create }}
{{- default (include "otel-demo.order.fullname" .) .Values.order.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.order.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "otel-demo.payment.serviceAccountName" -}}
{{- if .Values.payment.serviceAccount.create }}
{{- default (include "otel-demo.payment.fullname" .) .Values.payment.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.payment.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "otel-demo.users.serviceAccountName" -}}
{{- if .Values.users.serviceAccount.create }}
{{- default (include "otel-demo.users.fullname" .) .Values.users.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.users.serviceAccount.name }}
{{- end }}
{{- end }}
