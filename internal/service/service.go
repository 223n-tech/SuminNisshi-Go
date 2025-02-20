// ビジネスロジック層の構造体を定義
package service

import (
	"github.com/223n-tech/SuiminNisshi-Go/internal/repository"
)

/*
	Service サービスの構造体
*/
type Service struct {
	repo repository.RepositoryInterface
}

/*
	NewService サービスを作成
*/
func NewService(repo repository.RepositoryInterface) *Service {
	return &Service{
		repo: repo,
	}
}

/*
	Close サービスをクローズ
*/
func (s *Service) Close() error {
	return s.repo.Close()
}
