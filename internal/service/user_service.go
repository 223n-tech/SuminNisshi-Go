// internal/service/user_service.go
// user_serviceは、ユーザー関連のサービスを提供します。

// Package service provides application services.
package service

import (
	"context"
	"errors"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
    ErrEmptyPassword = errors.New("password cannot be empty")
    ErrPasswordMismatch = errors.New("passwords do not match")
    ErrPasswordTooShort = errors.New("password must be at least 8 characters long")
    ErrInvalidPasswordFormat = errors.New("password format is invalid")
    ErrEmptyDisplayName = errors.New("display name cannot be empty")
    ErrEmptyEmail = errors.New("email cannot be empty")
    ErrInvalidEmail = errors.New("email format is invalid")
    ErrEmptyTimezone = errors.New("timezone cannot be empty")
    ErrEmailAlreadyExists = errors.New("email already exists")
)

// ユーザー関連のサービス
type UserService struct {
	s *Service
}

// 新しいUserServiceを作成
func NewUserService(s *Service) *UserService {
	return &UserService{s: s}
}

// ユーザー登録
func (s *UserService) Register(ctx context.Context, email, displayName, password string) (*models.User, error) {
	// メールアドレスの重複チェック
	existing, err := s.s.repo.User().GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	// パスワードのハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		DisplayName:  displayName,
		PasswordHash: string(hash),
	}

	// ユーザーの作成
	if err := s.s.repo.User().Create(ctx, user); err != nil {
		return nil, err
	}

	// デフォルトの睡眠設定を作成
	pref := s.s.repo.UserSleepPreference().GetDefaultPreference(user.ID)
	if err := s.s.repo.UserSleepPreference().Create(ctx, pref); err != nil {
		return nil, err
	}

	return user, nil
}

// ユーザーIDからユーザーを取得
func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	return s.s.repo.User().GetByID(ctx, userID)
}

// ユーザー認証
func (s *UserService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.s.repo.User().GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := s.s.repo.User().UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, err
	}

	return user, nil
}

// ユーザーの睡眠設定を取得
func (s *UserService) GetSleepPreference(ctx context.Context, userID int64) (*models.UserSleepPreference, error) {
	pref, err := s.s.repo.UserSleepPreference().GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if pref == nil {
		// デフォルト設定を作成して返す
		pref = s.s.repo.UserSleepPreference().GetDefaultPreference(userID)
		if err := s.s.repo.UserSleepPreference().Create(ctx, pref); err != nil {
			return nil, err
		}
	}
	return pref, nil
}

// ユーザーの睡眠設定を更新
func (s *UserService) UpdateSleepPreference(ctx context.Context, pref *models.UserSleepPreference) error {
	return s.s.repo.UserSleepPreference().Update(ctx, pref)
}

// ユーザープロフィールを更新
func (s *UserService) UpdateProfile(ctx context.Context, user *models.User) error {
	return s.s.repo.User().Update(ctx, user)
}

// パスワードを更新
func (s *UserService) UpdatePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error {
	user, err := s.s.repo.User().GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hash)
	return s.s.repo.User().Update(ctx, user)
}

// アカウントを削除
func (s *UserService) DeleteAccount(ctx context.Context, userID int64) error {
	return s.s.Transaction(ctx, func(ctx context.Context) error {
		if err := s.s.repo.UserSleepPreference().Delete(ctx, userID); err != nil {
			return err
		}
		return s.s.repo.User().Delete(ctx, userID)
	})
}

// パスワードリセットの開始
func (s *UserService) InitiatePasswordReset(ctx context.Context, email string) (string, error) {
    // TODO: 実装
    // トークン生成とデータベースへの保存
    return "reset_token", nil
}

// リセットトークンの検証
func (s *UserService) ValidateResetToken(ctx context.Context, token string) (bool, error) {
    // TODO: 実装
    // トークンの検証
    return true, nil
}

// パスワードリセットの実行
func (s *UserService) CompletePasswordReset(ctx context.Context, token, newPassword string) error {
    // この関数が未実装のため追加
    // ...実装...
    return nil
}
