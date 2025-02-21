// テンプレートを管理するハンドラーです。
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

/*
	TemplateData テンプレートデータ
*/
type TemplateData struct {
	Title      string
	ActiveMenu string
	User       *models.User              // models.Userを使用
	Data       map[string]interface{}
	Flash      *Flash
	Meta       map[string]interface{}    // Metaフィールドを追加
}

/*
	Flash フラッシュメッセージ
*/
type Flash struct {
	Type    string // success, error, warning, info
	Message string
}

/*
	TemplateManager テンプレートマネージャ
*/
type TemplateManager struct {
	templates      map[string]*template.Template
	mutex          sync.RWMutex
	basePath       string
	funcMap        template.FuncMap
	embedFS        *embed.FS
	standalonePages []string
}

/*
	NewTemplateManager は TemplateManager を作成します。
*/
func NewTemplateManager(basePath string, embedFS *embed.FS) *TemplateManager {
	return &TemplateManager{
		templates:       make(map[string]*template.Template),
		basePath:        basePath,
		embedFS:         embedFS,
		funcMap:         makeTemplateFuncMap(),
		standalonePages: []string{
			"login.html",
			"register.html",
			"delete-account.html",
			"forgot-password.html",
			"reset-password.html",
			"404.html",
			"500.html",
			"403.html",
		},
	}
}

/*
	makeTemplateFuncMap はテンプレート関数のマップを作成します。
*/
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


/*
	LoadTemplates はテンプレートを読み込みます。
*/
func (tm *TemplateManager) LoadTemplates() error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	logger := log.New(os.Stdout, "[SuiminNisshi] ", log.LstdFlags|log.Lshortfile)

	// ベースレイアウトを読み込み
	logger.Printf("[START] Loading templates from %s", tm.basePath)
	layouts, err := filepath.Glob(filepath.Join(tm.basePath, "layouts/*.html"))
	if err != nil {
		return fmt.Errorf("[NG] error loading layouts: %v", err)
	}
	logger.Printf("[OK] Loaded layouts: %v", layouts)

	// パーシャルテンプレートを読み込み
	logger.Printf("[START] Loading partials from %s", tm.basePath)
	partials, err := filepath.Glob(filepath.Join(tm.basePath, "partials/*.html"))
	if err != nil {
		return fmt.Errorf("[NG] error loading partials: %v", err)
	}
	logger.Printf("[OK] Loaded partials: %v", partials)

	// 通常ページテンプレートを読み込み
	logger.Printf("[START] Loading pages from %s", tm.basePath)
	pages, err := filepath.Glob(filepath.Join(tm.basePath, "pages/*.html"))
	if err != nil {
		return fmt.Errorf("[NG] error loading pages: %v", err)
	}
	logger.Printf("[OK] Loaded pages: %v", pages)

	// エラーページテンプレートを読み込み
	logger.Printf("[START] Loading error pages from %s", tm.basePath)
	errorPages, err := filepath.Glob(filepath.Join(tm.basePath, "errors/*.html"))
	if err != nil {
		return fmt.Errorf("[NG] error loading error pages: %v", err)
	}
	logger.Printf("[OK] Loaded error pages: %v", errorPages)

	// スタンドアロンページの読み込み
	logger.Printf("[START] Loading standalone pages from %s", tm.basePath)
	for _, page := range tm.standalonePages {
		var fullPath string
		if slices.Contains([]string{"404.html", "500.html", "403.html"}, page) {
			fullPath = filepath.Join(tm.basePath, "errors", page)
		} else {
			fullPath = filepath.Join(tm.basePath, "pages", page)
		}
		template, err := template.New(page).Funcs(tm.funcMap).ParseFiles(fullPath)
		if err != nil {
			return fmt.Errorf("[NG] error parsing %s: %v", page, err)
		}
		tm.templates[page] = template
		logger.Printf("[OK] Loaded standalone template: %s", page)
	}

	// その他のページはレイアウトとパーシャルを含む
	logger.Printf("[START] Loading normal pages from %s", tm.basePath)
	for _, page := range pages {
		name := filepath.Base(page)
		// スタンドアロンページはスキップ
		if slices.Contains(tm.standalonePages, name) {
			logger.Printf("[SKIP] Skipping standalone page: %s", name)
			continue
		}
		files := append(append(layouts, partials...), page)
		template, err := template.New("base.html").Funcs(tm.funcMap).ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("[NG] error parsing %s: %v", name, err)
		}
		tm.templates[name] = template
		logger.Printf("[OK] Loaded template: %s with files: %v", name, files)
	}

	logger.Printf("[OK] All loaded templates: %v", tm.GetTemplateNames(false))
	return nil
}

/*
	Render はテンプレートをレンダリングします。
*/
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

	// スタンドアロンページの処理
	if slices.Contains(tm.standalonePages, name) {
		return template.Execute(w, data)
	}

	// その他のページはbase.htmlを使用
	return template.ExecuteTemplate(w, "base.html", data)
}

/*
	GetTemplateNames はテンプレート名のスライスを返します。
*/
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

/*
	ReloadTemplates はテンプレートを再読み込みします。
*/
func (tm *TemplateManager) ReloadTemplates() error {
	return tm.LoadTemplates()
}

/*
	contains はスライスに指定した要素が含まれているかどうかを返します。
*/
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
