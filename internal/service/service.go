// internal/service/service.go
// serviceは、アプリケーションのサービス層を提供します。

// Package service provides application services.
package service

import (
	"context"
	"log"

	"github.com/223n-tech/SuiminNisshi-Go/internal/repository"
)

// アプリケーションのサービス層を表す構造体
type Service struct {
    repo   repository.Repository
    logger *LoggerService
    user   *UserService
    diary  *SleepDiaryService
    record *SleepRecordService
    pdf    *PDFService
    email  *EmailService
}

// メール送信サービス
type EmailService struct {
    s *Service
}

// 新しいサービスインスタンスを作成
func NewService(repo repository.Repository, level LogLevel, logger *log.Logger) *Service {
    s := &Service{
        repo:   repo,
    }
    s.user = NewUserService(s)
    s.diary = NewSleepDiaryService(s)
    s.record = NewSleepRecordService(s)
    s.pdf = NewPDFService(s)
    s.email = NewEmailService(s)
    s.logger = NewLoggerService(level, logger)
    return s
}

// メール関連のサービスを取得
func (s *Service) Email() *EmailService {
    return s.email
}

// ログ関連のサービスを取得
func (s *Service) Logger() *LoggerService {
    return s.logger
}

// ユーザー関連のサービスを取得
func (s *Service) User() *UserService {
	return s.user
}

// 睡眠日誌関連のサービスを取得
func (s *Service) Diary() *SleepDiaryService {
	return s.diary
}

// 睡眠記録関連のサービスを取得
func (s *Service) Record() *SleepRecordService {
	return s.record
}

// PDF出力関連のサービスを取得
func (s *Service) PDF() *PDFService {
	return s.pdf
}

// トランザクションを実行
func (s *Service) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return s.repo.Transaction(ctx, func(_ repository.Repository) error {
		return fn(ctx)
	})
}
