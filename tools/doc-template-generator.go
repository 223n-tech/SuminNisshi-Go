// パッケージmainは、Goパッケージのドキュメントを自動生成するためのコマンドラインツールを提供します。
// このスクリプトは、抽象構文木（AST）を解析し、パッケージの構造体、関数、型情報からMarkdownドキュメントを生成します。
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// 解析されたGoパッケージの情報を保持する構造体です。
// パッケージ名、ドキュメント、型、関数の情報を格納します。
type PackageInfo struct {
	Name       string
	Doc        string
	Types      []TypeInfo
	Functions  []FunctionInfo
}

// 解析された型の情報を保持する構造体です。
// 型の名前、ドキュメント、フィールド情報を格納します。
type TypeInfo struct {
	Name    string
	Doc     string
	Fields  []FieldInfo
}

// 構造体のフィールド情報を保持する構造体です。
// フィールド名と型情報を格納します。
type FieldInfo struct {
	Name string
	Type string
}

// 解析された関数の情報を保持する構造体です。
// 関数名、ドキュメント、シグネチャ情報を格納します。
type FunctionInfo struct {
	Name       string
	Doc        string
	Signature  string
}

// テキストドキュメントの空行を整理し、余分な空白を削除する関数です。
// 連続する空行を1行に削減し、先頭と末尾の空行を取り除きます。
func cleanupDoc(text string) string {
	// 連続する空行を1行に、先頭と末尾の空行を削除
	lines := strings.Split(text, "\n")
	var cleanLines []string
	var lastLineEmpty bool
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if !lastLineEmpty {
				cleanLines = append(cleanLines, "")
				lastLineEmpty = true
			}
		} else {
			cleanLines = append(cleanLines, trimmed)
			lastLineEmpty = false
		}
	}

	// 先頭と末尾の空行を削除
	for len(cleanLines) > 0 && cleanLines[0] == "" {
		cleanLines = cleanLines[1:]
	}
	for len(cleanLines) > 0 && cleanLines[len(cleanLines)-1] == "" {
		cleanLines = cleanLines[:len(cleanLines)-1]
	}

	return strings.Join(cleanLines, "\n")
}

// 生成されたMarkdownドキュメントの余分な空行や空白を削除する関数です。
// 連続する空行を2行以下に制限し、行末の不要な空白を取り除きます。
func cleanupMarkdown(content string) string {
	// 連続する空行を削除する正規表現
	re := regexp.MustCompile(`\n{3,}`)
	cleaned := re.ReplaceAllString(content, "\n\n")
	
	// 行末の空白を削除
	re = regexp.MustCompile(`[ \t]+$`)
	cleaned = re.ReplaceAllString(cleaned, "")
	
	// ファイルの先頭と末尾の空行を削除
	cleaned = strings.TrimSpace(cleaned) + "\n"
	
	return cleaned
}

// AST型式を読みやすい文字列形式に変換する関数です。
// ポインタ型、配列型、パッケージ修飾型などを適切にフォーマットします。
func formatTypeString(t ast.Expr) string {
	switch typ := t.(type) {
	case *ast.Ident:
		return fmt.Sprintf("`%s`", typ.Name)
	case *ast.StarExpr:
		// ポインタ型の場合
		innerType := formatTypeString(typ.X)
		// 内側の型から先頭と末尾のバッククォートを削除
		innerType = strings.Trim(innerType, "`")
		return fmt.Sprintf("`*%s`", innerType)
	case *ast.SelectorExpr:
		// パッケージ修飾された型
		pkgPart := formatTypeString(typ.X)
		// パッケージ部分から先頭と末尾のバッククォートを削除
		pkgPart = strings.Trim(pkgPart, "`")
		return fmt.Sprintf("`%s.%s`", pkgPart, typ.Sel.Name)
	case *ast.ArrayType:
		// 配列型の場合
		innerType := formatTypeString(typ.Elt)
		// 内側の型から先頭と末尾のバッククォートを削除
		innerType = strings.Trim(innerType, "`")
		return fmt.Sprintf("`[]%s`", innerType)
	default:
		return fmt.Sprintf("`%T`", t)
	}
}

// 指定されたパッケージパスからパッケージ情報を抽出する関数です。
// 抽象構文木（AST）を解析し、パッケージの型、関数、ドキュメント情報を収集します。
func extractPackageInfo(pkgPath string) (*PackageInfo, error) {
	fset := token.NewFileSet()
	
	pkgs, err := parser.ParseDir(fset, pkgPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var packageInfo PackageInfo
	
	for _, pkg := range pkgs {
		packageInfo.Name = pkg.Name
		
		// パッケージのコメントを追加
		for _, file := range pkg.Files {
			if file.Doc != nil {
				packageInfo.Doc = cleanupDoc(file.Doc.Text())
				break
			}
		}

		astPkg := &ast.Package{
			Name:  pkg.Name,
			Files: pkg.Files,
		}
		docPkg := doc.New(astPkg, pkgPath, doc.AllDecls)

		// 型の情報を抽出
		for _, t := range docPkg.Types {
			typeInfo := TypeInfo{
				Name: t.Name,
				Doc:  cleanupDoc(t.Doc),
			}

			// フィールドの情報
			for _, spec := range t.Decl.Specs {
				if structType, ok := spec.(*ast.TypeSpec).Type.(*ast.StructType); ok {
					for _, field := range structType.Fields.List {
						if len(field.Names) > 0 {
							fieldName := field.Names[0].Name
							fieldType := formatTypeString(field.Type)
							
							typeInfo.Fields = append(typeInfo.Fields, FieldInfo{
								Name: fieldName,
								Type: fieldType,
							})
						}
					}
				}
			}

			packageInfo.Types = append(packageInfo.Types, typeInfo)
		}

		// 関数の情報を抽出
		for _, f := range docPkg.Funcs {
			packageInfo.Functions = append(packageInfo.Functions, FunctionInfo{
				Name: f.Name,
				Doc:  cleanupDoc(f.Doc),
			})
		}

		break // 最初のパッケージのみ処理
	}

	return &packageInfo, nil
}

// パッケージ情報からMarkdownドキュメントを生成する関数です。
// 抽出したパッケージ情報とテンプレートを使用して、ドキュメントファイルを作成します。
func generateDocumentation(pkgPath, templatePath, outputPath string) error {
	// パッケージ情報の抽出
	pkgInfo, err := extractPackageInfo(pkgPath)
	if err != nil {
		return err
	}

	// テンプレートの読み込み
	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	// カスタム関数マップを作成
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	// テンプレートの解析
	tmpl, err := template.New("package-doc").Funcs(funcMap).Parse(string(tmplBytes))
	if err != nil {
		return err
	}

	// 出力ディレクトリ作成
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	// ドキュメントの生成
	var docContent bytes.Buffer
	if err := tmpl.Execute(&docContent, pkgInfo); err != nil {
		return err
	}

	// Markdownのクリーンアップ
	cleanedContent := cleanupMarkdown(docContent.String())

	// ファイルに書き出し
	return os.WriteFile(outputPath, []byte(cleanedContent), 0644)
}

// コマンドライン引数からパッケージパス、テンプレートファイル、出力ファイルを受け取り、
// ドキュメント生成プロセスを実行します。
func main() {
	if len(os.Args) != 4 {
		fmt.Println("使用法: go run doc-template-generator.go <パッケージのパス> <テンプレートファイル> <出力ファイル>")
		os.Exit(1)
	}

	pkgPath := os.Args[1]
	templatePath := os.Args[2]
	outputPath := os.Args[3]

	err := generateDocumentation(pkgPath, templatePath, outputPath)
	if err != nil {
		fmt.Println("エラー:", err)
		os.Exit(1)
	}

	fmt.Printf("ドキュメントを %s に生成しました\n", outputPath)
}
