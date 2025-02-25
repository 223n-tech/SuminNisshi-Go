// internal/service/mail_service.go

// package logger implements package that extends logging capabilities
package service

import "context"

//  新しいメールサービスを作成
func NewEmailService(s *Service) *EmailService {
    return &EmailService{s: s}
}

// ウェルカムメールを送信
func (s *EmailService) SendWelcomeEmail(ctx context.Context, email, name string) error {
    // 実装
    return nil
}
