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
