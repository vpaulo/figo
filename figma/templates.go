package figma

const CssVariablesTemplate = `
{{- range $selector, $rules := . }}
{{if ne $selector ":root" -}}
.
{{- end -}}
{{$selector}} {
{{- range $rules }}
 	{{ . }}
{{- end }}
}
{{ end }}`

const CssComponentsTemplate = `
{{- define "component" }}
{{ if .Styles -}}
{{ .Selectors }} {
	{{- range $selector, $value := .Styles }}
	{{ $selector }}: {{ $value }};
	{{- end }}
}
{{- end -}}
{{- range .Children }}
{{ template "component" . }}
{{- end -}}
{{- end -}}
{{ template "component" . }}
`

const HtmlComponentsTemplate = `
{{- define "component" }}
{{- if gt (len .Variants) 0 }}
{{- template "component" (index .Children 0) -}}
{{- else }}
<div class="{{.Name}}">
	{{- range .Children}}{{template "component" .}}{{end -}}
</div>
{{ end -}}
{{ end -}}
{{template "component" .}}
`
