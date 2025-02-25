// Package handler provides HTTP handlers for the application.
package handler

// internal/handler/template.go
// templateは、テンプレートのハンドリングを提供します。

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
	"github.com/223n-tech/SuiminNisshi-Go/internal/service"
)

// テンプレートに渡すデータの構造体
type TemplateData struct {
	Title      string
	ActiveMenu string
	User       *models.User
	Data       map[string]interface{}
	Flash      *Flash
	Meta       map[string]interface{}
}

// フラッシュメッセージの構造体
type Flash struct {
	Type    string // success, info, warning, danger
	Message string
}

// テンプレートを管理する構造体
type TemplateManager struct {
	templates       map[string]*template.Template
	mutex          sync.RWMutex
	basePath       string
	funcMap        template.FuncMap
	embedFS        *embed.FS
	standalonePages []string
	logger         *log.Logger
	service        *service.Service
}

// 新しいTemplateManagerを作成
func NewTemplateManager(basePath string, embedFS *embed.FS, logger *log.Logger, svc *service.Service) *TemplateManager {
	return &TemplateManager{
		templates:       make(map[string]*template.Template),
		basePath:        basePath,
		embedFS:         embedFS,
		funcMap:         makeTemplateFuncMap(),
		standalonePages: []string{
			"login.html",
			"register.html",
			"forgot-password.html",
			"reset-password.html",
			"404.html",
			"500.html",
			"403.html",
		},
		logger:  logger,
		service: svc,
	}
}

// テンプレート関数を作成
func makeTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		// 数値操作
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
		
		// 日時操作
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"formatTime": func(t time.Time) string {
			return t.Format("15:04")
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04")
		},
		"now": time.Now,
		
		// 文字列操作
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"excerpt": func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},
		
		// スライス操作
		"seq": func(start, end int) []int {
			seq := make([]int, end-start)
			for i := range seq {
				seq[i] = start + i
			}
			return seq
		},
		"contains": func(slice []string, item string) bool {
			return slices.Contains(slice, item)
		},
		
		// 睡眠記録用
		"sleepQualityClass": func(score int) string {
			switch {
			case score >= 90:
				return "bg-success"
			case score >= 70:
				return "bg-info"
			case score >= 50:
				return "bg-warning"
			default:
				return "bg-danger"
			}
		},
		"formatDuration": func(duration float64) string {
			hours := int(duration)
			minutes := int((duration - float64(hours)) * 60)
			return fmt.Sprintf("%d時間%d分", hours, minutes)
		},
	}
}

// テンプレートをロード
func (tm *TemplateManager) LoadTemplates() error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// レイアウトファイルの読み込み
	layouts, err := filepath.Glob(filepath.Join(tm.basePath, "layouts/*.html"))
	if err != nil {
		return fmt.Errorf("レイアウトファイルの読み込みに失敗: %v", err)
	}

	// パーシャルファイルの読み込み
	partials, err := filepath.Glob(filepath.Join(tm.basePath, "partials/*.html"))
	if err != nil {
		return fmt.Errorf("パーシャルファイルの読み込みに失敗: %v", err)
	}

	// ページファイルの読み込み
	pages, err := filepath.Glob(filepath.Join(tm.basePath, "pages/*.html"))
	if err != nil {
		return fmt.Errorf("ページファイルの読み込みに失敗: %v", err)
	}

	// エラーページの読み込み
	_, err = filepath.Glob(filepath.Join(tm.basePath, "errors/*.html"))
	if err != nil {
		return fmt.Errorf("エラーページの読み込みに失敗: %v", err)
	}

	// スタンドアロンページのロード
	for _, page := range tm.standalonePages {
		var fullPath string
		if slices.Contains([]string{"404.html", "500.html", "403.html"}, page) {
			fullPath = filepath.Join(tm.basePath, "errors", page)
		} else {
			fullPath = filepath.Join(tm.basePath, "pages", page)
		}

		tmpl, err := template.New(page).Funcs(tm.funcMap).ParseFiles(fullPath)
		if err != nil {
			return fmt.Errorf("スタンドアロンページのパースに失敗 %s: %v", page, err)
		}
		tm.templates[page] = tmpl
		tm.logger.Printf("スタンドアロンページをロード: %s", page)
	}

	// 通常のページをレイアウトとパーシャルと共にロード
	for _, page := range pages {
		name := filepath.Base(page)
		if slices.Contains(tm.standalonePages, name) {
			continue
		}

		files := append(append(layouts, partials...), page)
		tmpl, err := template.New("base.html").Funcs(tm.funcMap).ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("テンプレートのパースに失敗 %s: %v", name, err)
		}
		tm.templates[name] = tmpl
		tm.logger.Printf("テンプレートをロード: %s", name)
	}

	return nil
}

// テンプレートをレンダリング
func (tm *TemplateManager) Render(w http.ResponseWriter, name string, data *TemplateData) error {
	tm.mutex.RLock()
	tmpl, exists := tm.templates[name]
	tm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("テンプレートが見つかりません: %s", name)
	}

	// デフォルト値の設定
	if data == nil {
		data = &TemplateData{}
	}
	if data.Data == nil {
		data.Data = make(map[string]interface{})
	}
	if data.Meta == nil {
		data.Meta = make(map[string]interface{})
	}

	// Content-Typeの設定
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// スタンドアロンページの処理
	if slices.Contains(tm.standalonePages, name) {
		return tmpl.Execute(w, data)
	}

	// 通常のページ処理
	return tmpl.ExecuteTemplate(w, "base.html", data)
}

// テンプレート名の一覧を取得
func (tm *TemplateManager) GetTemplateNames() []string {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	names := make([]string, 0, len(tm.templates))
	for name := range tm.templates {
		names = append(names, name)
	}
	return names
}

// テンプレートを再読み込み
func (tm *TemplateManager) ReloadTemplates() error {
	tm.logger.Println("テンプレートの再読み込みを開始")
	if err := tm.LoadTemplates(); err != nil {
		tm.logger.Printf("テンプレートの再読み込みに失敗: %v", err)
		return err
	}
	tm.logger.Println("テンプレートの再読み込みが完了")
	return nil
}

// カスタムテンプレート関数を追加
func (tm *TemplateManager) AddCustomFunc(name string, fn interface{}) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tm.funcMap[name] = fn
	tm.logger.Printf("カスタム関数を追加: %s", name)
}

// テンプレートキャッシュをクリア
func (tm *TemplateManager) ClearCache() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tm.templates = make(map[string]*template.Template)
	tm.logger.Println("テンプレートキャッシュをクリア")
}
