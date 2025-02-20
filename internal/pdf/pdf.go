// PDFファイルを作成するためのパッケージ
package pdf

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

/*
	PDFの構成
*/
const (
	// ページ設定
	pageWidth     = 297.0  // A4横
	pageHeight    = 210.0  // A4縦
	marginTop     = 20.0
	marginBottom  = 20.0
	marginLeft    = 15.0
	marginRight   = 15.0

	// グリッドの設定
	cellWidth     = 10.0   // 30分あたりの幅
	cellHeight    = 8.0    // セルの高さ
	timeRowHeight = 6.0    // 時刻行の高さ

	// フォントサイズ
	titleFontSize   = 14
	normalFontSize  = 10
	symbolFontSize  = 12
)

/*
	SleepDiaryPDF PDFファイルを作成するための構造体
*/
type SleepDiaryPDF struct {
	pdf *gofpdf.Fpdf
}

/*
	NewSleepDiaryPDF PDFインスタンスの作成
*/
func NewSleepDiaryPDF() *SleepDiaryPDF {
	pdf := gofpdf.New("L", "mm", "A4", "")
	
	// フォントの追加
	pdf.AddFont("IPAGothic", "", "ipag.ttf")
	
	return &SleepDiaryPDF{
		pdf: pdf,
	}
}

/*
	SetHeader ヘッダーの設定
*/
func (s *SleepDiaryPDF) SetHeader(userName string) {
	s.pdf.AddPage()
	s.pdf.SetFont("IPAGothic", "", titleFontSize)
	
	// タイトル
	s.pdf.Cell(40, 10, "睡眠日誌")
	
	// ユーザー名
	s.pdf.SetFont("IPAGothic", "", normalFontSize)
	s.pdf.Cell(40, 10, fmt.Sprintf("お名前：%s", userName))
}

/*
	DrawTimeGrid 時間軸グリッドの描画
*/
func (s *SleepDiaryPDF) DrawTimeGrid(x, y float64) {
	s.pdf.SetFont("IPAGothic", "", normalFontSize)
	
	// 時刻の描画（00-23時）
	for hour := 0; hour < 24; hour++ {
		xPos := x + float64(hour*2)*cellWidth
		s.pdf.Text(xPos, y, fmt.Sprintf("%02d", hour))
	}

	// グリッド線の描画
	for i := 0; i <= 48; i++ { // 30分毎、48マス
		xPos := x + float64(i)*cellWidth
		
		// 時間単位の区切り線は太く
		if i%2 == 0 {
			s.pdf.SetLineWidth(0.3)
		} else {
			s.pdf.SetLineWidth(0.1)
		}
		
		s.pdf.Line(xPos, y+timeRowHeight, xPos, y+timeRowHeight+4*cellHeight)
	}

	// 横線の描画（4段分）
	for i := 0; i <= 4; i++ {
		yPos := y + timeRowHeight + float64(i)*cellHeight
		s.pdf.SetLineWidth(0.3)
		s.pdf.Line(x, yPos, x+48*cellWidth, yPos)
	}
}

/*
	DrawLegend 凡例の描画
*/
func (s *SleepDiaryPDF) DrawLegend(x, y float64) {
	s.pdf.SetFont("IPAGothic", "", normalFontSize)
	
	legends := []struct {
		symbol string
		text   string
	}{
		{"■", "睡眠中"},
		{"╱", "床で覚醒"},
		{"□", "通常覚醒"},
		{"Z", "強い眠気"},
		{"×", "睡眠薬服用"},
		{"▲", "朝食"},
		{"●", "昼食"},
		{"■", "夕食"},
		{"○", "軽食"},
	}

	for i, legend := range legends {
		yPos := y + float64(i)*6
		s.pdf.SetFont("IPAGothic", "", symbolFontSize)
		s.pdf.Text(x, yPos, legend.symbol)
		s.pdf.SetFont("IPAGothic", "", normalFontSize)
		s.pdf.Text(x+6, yPos, legend.text)
	}
}

/*
	DrawSleepState 睡眠状態の描画
*/
func (s *SleepDiaryPDF) DrawSleepState(x, y float64, stateSymbol string, timeSlot time.Time) {
	hourFloat := float64(timeSlot.Hour()) + float64(timeSlot.Minute())/60.0
	xPos := x + hourFloat*2*cellWidth
	
	s.pdf.SetFont("IPAGothic", "", symbolFontSize)
	s.pdf.Text(xPos+1, y, stateSymbol)
}

/*
	DrawEvent イベントの描画
*/
func (s *SleepDiaryPDF) DrawEvent(x, y float64, eventSymbol string, timeSlot time.Time) {
	hourFloat := float64(timeSlot.Hour()) + float64(timeSlot.Minute())/60.0
	xPos := x + hourFloat*2*cellWidth
	
	s.pdf.SetFont("IPAGothic", "", symbolFontSize)
	s.pdf.Text(xPos+1, y, eventSymbol)
}

/*
	DrawMeal 食事の描画
*/
func (s *SleepDiaryPDF) DrawMeal(x, y float64, mealSymbol string, timeSlot time.Time) {
	hourFloat := float64(timeSlot.Hour()) + float64(timeSlot.Minute())/60.0
	xPos := x + hourFloat*2*cellWidth
	
	s.pdf.SetFont("IPAGothic", "", symbolFontSize)
	s.pdf.Text(xPos+1, y, mealSymbol)
}

/*
	DrawNote 備考欄の描画
*/
func (s *SleepDiaryPDF) DrawNote(x, y float64, note string) {
	s.pdf.SetFont("IPAGothic", "", normalFontSize)
	s.pdf.Text(x, y, "備考：")
	s.pdf.SetFont("IPAGothic", "", normalFontSize)
	s.pdf.MultiCell(pageWidth-marginLeft-marginRight, 5, note, "", "", false)
}

/*
	SavePDF PDFの保存
*/
func (s *SleepDiaryPDF) SavePDF(filename string) error {
	return s.pdf.OutputFileAndClose(filename)
}

// Usage example:
/*
func main() {
	pdf := NewSleepDiaryPDF()
	pdf.SetHeader("山田太郎")
	
	// グリッドの描画開始位置
	x := marginLeft
	y := marginTop + 20
	
	pdf.DrawTimeGrid(x, y)
	pdf.DrawLegend(pageWidth-marginRight-30, y)
	
	// 睡眠状態の記録
	pdf.DrawSleepState(x, y+timeRowHeight+cellHeight, "■", time.Date(2024, 2, 21, 23, 0, 0, 0, time.Local))
	
	// イベントの記録
	pdf.DrawEvent(x, y+timeRowHeight+2*cellHeight, "Z", time.Date(2024, 2, 21, 15, 30, 0, 0, time.Local))
	
	// 食事の記録
	pdf.DrawMeal(x, y+timeRowHeight+3*cellHeight, "▲", time.Date(2024, 2, 21, 7, 0, 0, 0, time.Local))
	
	// 備考の記録
	pdf.DrawNote(x, y+timeRowHeight+5*cellHeight, "夕飯で晩酌 3回トイレに起きた")
	
	pdf.SavePDF("sleep_diary.pdf")
}
*/
