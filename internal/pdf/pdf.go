// internal/pdf/pdf.go
// pdfは、PDFを生成するための機能を提供します。

// Package pdf provides functionality to generate PDFs.
package pdf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/signintech/gopdf"
)

// PDFを生成するための構造体
type Generator struct {
	pdf      *gopdf.GoPdf
	fontPath string
}

// 新しいPDFジェネレーターを作成
func New(fontPath string) *Generator {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	
	return &Generator{
		pdf:      pdf,
		fontPath: fontPath,
	}
}

// 睡眠記録のPDFを生成
func (g *Generator) GenerateSleepRecord(data *SleepRecordData) (*bytes.Buffer, error) {
	// フォントの設定
	err := g.setupFont()
	if err != nil {
		return nil, fmt.Errorf("failed to setup font: %v", err)
	}

	// ページを追加
	g.pdf.AddPage()

	// ヘッダーを設定
	if err := g.pdf.SetFont("gothic", "", 16); err != nil {
		return nil, fmt.Errorf("failed to set font: %v", err)
	}
	if err := g.pdf.Text("睡眠記録"); err != nil {
		return nil, fmt.Errorf("failed to write text: %v", err)
	}
	g.pdf.SetY(g.pdf.GetY() + 15)

	// 基本情報
	if err := g.pdf.SetFont("gothic", "", 12); err != nil {
		return nil, fmt.Errorf("failed to set font: %v", err)
	}
	if err := g.writeBasicInfo(data); err != nil {
		return nil, fmt.Errorf("failed to write basic info: %v", err)
	}

	// 睡眠データ
	g.pdf.SetY(g.pdf.GetY() + 10)
	if err := g.writeSleepData(data); err != nil {
		return nil, fmt.Errorf("failed to write sleep data: %v", err)
	}

	// 統計データ
	g.pdf.SetY(g.pdf.GetY() + 10)
	if err := g.writeStatistics(data); err != nil {
		return nil, fmt.Errorf("failed to write statistics: %v", err)
	}

	// PDFをバッファに出力
	var buf bytes.Buffer
	if _, err = g.pdf.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %v", err)
	}

	return &buf, nil
}

// フォントの設定
func (g *Generator) setupFont() error {
	err := g.pdf.AddTTFFont("gothic", g.fontPath)
	if err != nil {
		return fmt.Errorf("failed to add font: %v", err)
	}
	return nil
}

// 基本情報を書き込み
func (g *Generator) writeBasicInfo(data *SleepRecordData) error {
	if err := g.pdf.SetFont("gothic", "", 14); err != nil {
		return err
	}
	if err := g.pdf.Text("基本情報"); err != nil {
		return err
	}
	g.pdf.SetY(g.pdf.GetY() + 10)

	if err := g.pdf.SetFont("gothic", "", 12); err != nil {
		return err
	}
	
	// 期間
	if err := g.pdf.Text("期間:"); err != nil {
		return err
	}
	g.pdf.SetX(60)
	periodText := fmt.Sprintf("%s 〜 %s", 
		data.StartDate.Format("2006/01/02"),
		data.EndDate.Format("2006/01/02"))
	if err := g.pdf.Text(periodText); err != nil {
		return err
	}
	g.pdf.SetY(g.pdf.GetY() + 8)

	// 記録日数
	if err := g.pdf.Text("記録日数:"); err != nil {
		return err
	}
	g.pdf.SetX(60)
	if err := g.pdf.Text(fmt.Sprintf("%d 日", data.TotalDays)); err != nil {
		return err
	}
	g.pdf.SetY(g.pdf.GetY() + 8)

	return nil
}

// 睡眠データを書き込み
func (g *Generator) writeSleepData(data *SleepRecordData) error {
	if err := g.pdf.SetFont("gothic", "", 14); err != nil {
		return err
	}
	if err := g.pdf.Text("睡眠データ"); err != nil {
		return err
	}
	g.pdf.SetY(g.pdf.GetY() + 10)

	// テーブルヘッダー
	if err := g.pdf.SetFont("gothic", "", 12); err != nil {
		return err
	}

	headers := []string{"日付", "就寝時刻", "起床時刻", "睡眠時間", "睡眠スコア"}
	xPositions := []float64{10, 50, 90, 130, 170}
	
	for i, header := range headers {
		g.pdf.SetX(xPositions[i])
		if err := g.pdf.Text(header); err != nil {
			return err
		}
	}
	g.pdf.SetY(g.pdf.GetY() + 8)

	// テーブルデータ
	for _, record := range data.Records {
		g.pdf.SetX(xPositions[0])
		if err := g.pdf.Text(record.Date.Format("2006/01/02")); err != nil {
			return err
		}
		
		g.pdf.SetX(xPositions[1])
		if err := g.pdf.Text(record.BedTime); err != nil {
			return err
		}
		
		g.pdf.SetX(xPositions[2])
		if err := g.pdf.Text(record.WakeTime); err != nil {
			return err
		}
		
		g.pdf.SetX(xPositions[3])
		if err := g.pdf.Text(fmt.Sprintf("%.1f時間", record.Duration)); err != nil {
			return err
		}
		
		g.pdf.SetX(xPositions[4])
		if err := g.pdf.Text(fmt.Sprintf("%d点", record.Score)); err != nil {
			return err
		}
		
		g.pdf.SetY(g.pdf.GetY() + 8)
	}

	return nil
}

// 統計データを書き込み
func (g *Generator) writeStatistics(data *SleepRecordData) error {
	if err := g.pdf.SetFont("gothic", "", 14); err != nil {
		return err
	}
	if err := g.pdf.Text("統計データ"); err != nil {
		return err
	}
	g.pdf.SetY(g.pdf.GetY() + 10)

	if err := g.pdf.SetFont("gothic", "", 12); err != nil {
		return err
	}
	
	stats := []struct {
		label string
		value string
	}{
		{"平均睡眠時間:", fmt.Sprintf("%.1f時間", data.AverageDuration)},
		{"平均就寝時刻:", data.AverageBedTime},
		{"平均起床時刻:", data.AverageWakeTime},
		{"平均睡眠スコア:", fmt.Sprintf("%.1f点", data.AverageScore)},
	}

	for _, stat := range stats {
		if err := g.pdf.Text(stat.label); err != nil {
			return err
		}
		g.pdf.SetX(100)
		if err := g.pdf.Text(stat.value); err != nil {
			return err
		}
		g.pdf.SetY(g.pdf.GetY() + 8)
	}

	return nil
}

// PDFに出力する睡眠記録データ
type SleepRecordData struct {
	StartDate       time.Time
	EndDate         time.Time
	TotalDays       int
	Records         []SleepRecord
	AverageDuration float64
	AverageBedTime  string
	AverageWakeTime string
	AverageScore    float64
}

// 個別の睡眠記録
type SleepRecord struct {
	Date     time.Time
	BedTime  string
	WakeTime string
	Duration float64
	Score    int
}
