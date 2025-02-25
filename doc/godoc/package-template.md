# パッケージ: {{.Name}}

## 1. 概要

{{.Doc}}

## 2. 型
{{if .Types}}

{{range $typeIndex, $type := .Types}}
### 2-{{add $typeIndex 1}}. `{{$type.Name}}`

{{$type.Doc}}
{{if $type.Fields}}
#### 2-{{add $typeIndex 1}}-1. フィールド

{{range $type.Fields}}
- `{{.Name}}`: {{.Type}}
{{end}}
{{end}}
{{end}}
{{else}}

なし
{{end}}

## 3. 関数
{{if .Functions}}

{{range $funcIndex, $func := .Functions}}
### 3-{{add $funcIndex 1}}. `{{$func.Name}}`

{{$func.Doc}}
{{end}}
{{else}}

なし
{{end}}
