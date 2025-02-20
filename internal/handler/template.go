package handler

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
)

// TemplateData テンプレートに渡すデータの構造体
type TemplateData struct {
	Title      string
	ActiveMenu string
	User       *models.User              // models.Userを使用
	Data       map[string]interface{}
	Flash      *Flash
	Meta       map[string]interface{}    // Metaフィールドを追加
}

// Flash フラッシュメッセージの構造体
type Flash struct {
	Type    string // success, error, warning, info
	Message string
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
		"now": time.Now,
	}
}

// LoadTemplates はテンプレートをロードします
func (tm *TemplateManager) LoadTemplates() error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	logger := log.New(os.Stdout, "[SuiminNisshi] ", log.LstdFlags|log.Lshortfile)

	// ベースレイアウトを読み込み
	logger.Printf("[START] Loading templates from %s", tm.basePath)  // デバッグ用
	layouts, err := filepath.Glob(filepath.Join(tm.basePath, "layouts/*.html"))
	if err != nil {
		return fmt.Errorf("[NG] error loading layouts: %v", err)
	}

	// パーシャルテンプレートを読み込み
	logger.Printf("[START] Loading partials from %s", tm.basePath)  // デバッグ用
	partials, err := filepath.Glob(filepath.Join(tm.basePath, "partials/*.html"))
	if err != nil {
		return fmt.Errorf("[NG] error loading partials: %v", err)
	}

	// 通常ページテンプレートを読み込み
	logger.Printf("[START] Loading pages from %s", tm.basePath)  // デバッグ用
	pages, err := filepath.Glob(filepath.Join(tm.basePath, "pages/*.html"))
	if err != nil {
		return fmt.Errorf("[NG] error loading pages: %v", err)
	}

	// エラーページテンプレートを読み込み
	logger.Printf("[START] Loading error pages from %s", tm.basePath)  // デバッグ用
	errorPages, err := filepath.Glob(filepath.Join(tm.basePath, "errors/*.html"))
	if err != nil {
		return fmt.Errorf("[NG] error loading error pages: %v", err)
	}

	// ログインと登録ページは独立したテンプレート
	logger.Printf("[START] Loading standalone pages from %s", tm.basePath)  // デバッグ用
	noTemplatePages := []string{"login.html", "register.html", "delete-account.html", "forgot-password.html", "reset-password.html"}
	for _, page := range noTemplatePages {
		fullPath := filepath.Join(tm.basePath, "pages", page)
		template, err := template.New(page).Funcs(tm.funcMap).ParseFiles(fullPath)
		if err != nil {
			return fmt.Errorf("[NG] error parsing %s: %v", page, err)
		}
		tm.templates[page] = template
	}

	// エラーページは独立したテンプレート
	logger.Printf("[START] Loading standalone error pages from %s", tm.basePath)  // デバッグ用
	for _, page := range errorPages {
		name := filepath.Base(page)
		template, err := template.New(name).Funcs(tm.funcMap).ParseFiles(page)
		if err != nil {
			return fmt.Errorf("error parsing error page %s: %v", name, err)
		}
		tm.templates[name] = template
		logger.Printf("[OK] Loaded error template: %s", name)  // デバッグ用
	}

	// その他のページはレイアウトとパーシャルを含む
	logger.Printf("[START] Loading normal pages from %s", tm.basePath)  // デバッグ用
	for i, page := range pages {
		name := filepath.Base(page)
		logger.Printf("[%d] Loading template: %s", i, name)  // デバッグ用
		// ログインと登録ページはスキップ
		if slices.Contains(noTemplatePages, name) {
			logger.Printf("[SKIP] Skipping standalone page: %s", name)  // デバッグ用
			continue
		}
		files := append(append(layouts, partials...), page)
		template, err := template.New("base.html").Funcs(tm.funcMap).ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("[NG] error parsing %s: %v", name, err)
		}
		tm.templates[name] = template
		logger.Printf("[OK] Loaded template: %s with files: %v", name, files)  // デバッグ用
	}

	logger.Printf("[OK] All loaded templates: %v", tm.GetTemplateNames(false))  // デバッグ用
	return nil
}

// GetTemplateNames は読み込まれたテンプレート名の一覧を返します
func (tm *TemplateManager) GetTemplateNames(isLock bool) []string {
	if (isLock) {
		tm.mutex.RLock()
		defer tm.mutex.RUnlock()	
	} 

	names := make([]string, 0, len(tm.templates))
	for name := range tm.templates {
		names = append(names, name)
	}
	return names
}

// Render はテンプレートをレンダリングします
func (tm *TemplateManager) Render(w http.ResponseWriter, name string, data *TemplateData) error {
	tm.mutex.RLock()
	template, exists := tm.templates[name]
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
		return template.Execute(w, data)
	}

	// その他のページはbase.htmlを使用
	return template.ExecuteTemplate(w, "base.html", data)
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
