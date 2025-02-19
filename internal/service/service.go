package service

import (
	"github.com/223n-tech/SuiminNisshi-Go/internal/repository"
)

// Service ビジネスロジック層の構造体
type Service struct {
	repo repository.RepositoryInterface
}

// NewService サービスインスタンスを作成
func NewService(repo repository.RepositoryInterface) *Service {
	return &Service{
		repo: repo,
	}
}

// Close サービスのクリーンアップ処理
func (s *Service) Close() error {
	return s.repo.Close()
}
