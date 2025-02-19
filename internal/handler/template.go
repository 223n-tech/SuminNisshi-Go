package handler

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// TemplateData はテンプレートに渡すデータの構造体です
type TemplateData struct {
	Title      string                 // ページタイトル
	ActiveMenu string                 // アクティブなメニュー項目
	User       *User                  // ログインユーザー情報
	Data       interface{}            // ページ固有のデータ
	Flash      *Flash                 // フラッシュメッセージ
	CSRF       string                 // CSRFトークン
	Meta       map[string]interface{} // メタデータ
}

// Flash はフラッシュメッセージの構造体です
type Flash struct {
	Type    string // success, info, warning, danger
	Message string
}

// User はユーザー情報の構造体です
type User struct {
	ID       int64
	Email    string
	Name     string
	IsAdmin  bool
	Settings map[string]interface{}
}

// TemplateManager はテンプレート管理を担当する構造体です
type TemplateManager struct {
	templates map[string]*template.Template
	mutex     sync.RWMutex
	basePath  string
	funcMap   template.FuncMap
	embedFS   *embed.FS
}

// NewTemplateManager は新しいTemplateManagerインスタンスを作成します
func NewTemplateManager(basePath string, embedFS *embed.FS) *TemplateManager {
	return &TemplateManager{
		templates: make(map[string]*template.Template),
		basePath:  basePath,
		embedFS:   embedFS,
		funcMap:   makeTemplateFuncMap(),
	}
}

// makeTemplateFuncMap はテンプレートで使用する関数マップを作成します
func makeTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"formatDate": func(date interface{}) string {
			return fmt.Sprintf("%v", date)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"seq": func(start, end int) []int {
			var result []int
			for i := start; i < end; i++ {
				result = append(result, i)
			}
			return result
		},
	}
}

// LoadTemplates はテンプレートをロードします
func (tm *TemplateManager) LoadTemplates() error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// ベースレイアウトを読み込み
	layouts, err := filepath.Glob(filepath.Join(tm.basePath, "layouts/*.html"))
	if err != nil {
		return fmt.Errorf("error loading layouts: %v", err)
	}

	// パーシャルテンプレートを読み込み
	partials, err := filepath.Glob(filepath.Join(tm.basePath, "partials/*.html"))
	if err != nil {
		return fmt.Errorf("error loading partials: %v", err)
	}

	// ページテンプレートを読み込み
	pages, err := filepath.Glob(filepath.Join(tm.basePath, "pages/*.html"))
	if err != nil {
		return fmt.Errorf("error loading pages: %v", err)
	}

	// ログインと登録ページは独立したテンプレート
	standalone := []string{"login.html", "register.html"}
	for _, page := range standalone {
		fullPath := filepath.Join(tm.basePath, "pages", page)
		tmpl, err := template.New(page).Funcs(tm.funcMap).ParseFiles(fullPath)
		if err != nil {
			return fmt.Errorf("error parsing %s: %v", page, err)
		}
		tm.templates[page] = tmpl
	}

	// その他のページはレイアウトとパーシャルを含む
	for _, page := range pages {
		name := filepath.Base(page)
		// スタンドアロンページはスキップ
		if contains(standalone, name) {
			continue
		}

		files := append(append(layouts, partials...), page)
		tmpl, err := template.New("base.html").Funcs(tm.funcMap).ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("error parsing %s: %v", name, err)
		}
		tm.templates[name] = tmpl
	}

	// エラーページの読み込み
	errorPages, err := filepath.Glob(filepath.Join(tm.basePath, "errors/*.html"))
	if err != nil {
		return fmt.Errorf("error loading error pages: %v", err)
	}

	for _, page := range errorPages {
		name := filepath.Base(page)
		tmpl, err := template.New(name).ParseFiles(page)
		if err != nil {
			return fmt.Errorf("error parsing error page %s: %v", name, err)
		}
		tm.templates[name] = tmpl
	}

	return nil
}

// Render はテンプレートをレンダリングします
func (tm *TemplateManager) Render(w http.ResponseWriter, name string, data *TemplateData) error {
	tm.mutex.RLock()
	tmpl, exists := tm.templates[name]
	tm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("template %s not found", name)
	}

	// デフォルト値の設定
	if data == nil {
		data = &TemplateData{}
	}
	if data.Meta == nil {
		data.Meta = make(map[string]interface{})
	}

	// スタンドアロンページと独立したテンプレートの処理
	standalone := []string{"login.html", "register.html", "404.html", "500.html", "403.html"}
	if contains(standalone, name) {
		return tmpl.Execute(w, data)
	}

	// その他のページはbase.htmlを使用
	return tmpl.ExecuteTemplate(w, "base.html", data)
}

// ReloadTemplates はテンプレートを再読み込みします
func (tm *TemplateManager) ReloadTemplates() error {
	return tm.LoadTemplates()
}

// contains は文字列スライスに特定の文字列が含まれているかチェックします
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
