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
