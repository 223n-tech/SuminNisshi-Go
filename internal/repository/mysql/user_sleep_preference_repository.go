// internal/repository/mysql/user_sleep_preference_repository.go
// user_sleep_preference_repositoryは、睡眠設定のリポジトリを提供します。

// Package mysql provides MySQL repository implementations.
package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/223n-tech/SuiminNisshi-Go/internal/models"
)

// UserSleepPreferenceRepositoryのMySQL実装
type UserSleepPreferenceRepository struct {
	repo *MySQLRepository
}

// ユーザーIDで睡眠設定を検索
func (r *UserSleepPreferenceRepository) GetByUserID(ctx context.Context, userID int64) (*models.UserSleepPreference, error) {
	query := `
		SELECT id, user_id, preferred_bedtime, preferred_wakeup_time, sleep_goal_hours, is_reminder_enabled, created, modified, deleted
		FROM users_sleep_preferences
		WHERE user_id = ? AND deleted IS NULL
	`

	pref := &models.UserSleepPreference{}
	err := r.repo.getDB().(*sql.DB).QueryRowContext(ctx, query, userID).Scan(
		&pref.ID,
		&pref.UserID,
		&pref.PreferredBedtime,
		&pref.PreferredWakeupTime,
		&pref.SleepGoalHours,
		&pref.IsReminderEnabled,
		&pref.Created,
		&pref.Modified,
		&pref.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return pref, nil
}

//新規睡眠設定を作成
func (r *UserSleepPreferenceRepository) Create(ctx context.Context, pref *models.UserSleepPreference) error {
	query := `
		INSERT INTO users_sleep_preferences (
			user_id, preferred_bedtime, preferred_wakeup_time,
			sleep_goal_hours, is_reminder_enabled, created, modified
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		pref.UserID,
		pref.PreferredBedtime,
		pref.PreferredWakeupTime,
		pref.SleepGoalHours,
		pref.IsReminderEnabled,
		now,
		now,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	pref.ID = id
	pref.Created = now
	pref.Modified = now

	return nil
}

// 睡眠設定を更新
func (r *UserSleepPreferenceRepository) Update(ctx context.Context, pref *models.UserSleepPreference) error {
	query := `
		UPDATE users_sleep_preferences
		SET preferred_bedtime = ?, preferred_wakeup_time = ?, sleep_goal_hours = ?, is_reminder_enabled = ?, modified = ?
		WHERE id = ? AND deleted IS NULL
	`

	now := time.Now()
	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		pref.PreferredBedtime,
		pref.PreferredWakeupTime,
		pref.SleepGoalHours,
		pref.IsReminderEnabled,
		now,
		pref.ID,
	)

	if err != nil {
		return err
	}

	pref.Modified = now
	return nil
}

// 睡眠設定を論理削除
func (r *UserSleepPreferenceRepository) Delete(ctx context.Context, userID int64) error {
	query := `
		UPDATE users_sleep_preferences
		SET deleted = ?
		WHERE user_id = ? AND deleted IS NULL
	`

	_, err := r.repo.getDB().(*sql.DB).ExecContext(ctx, query,
		time.Now(),
		userID,
	)

	return err
}

// デフォルトの睡眠設定を生成
func (r *UserSleepPreferenceRepository) GetDefaultPreference(userID int64) *models.UserSleepPreference {
	defaultBedtime, _ := time.Parse("15:04", "23:00")
	defaultWakeupTime, _ := time.Parse("15:04", "07:00")

	return &models.UserSleepPreference{
		UserID:              userID,
		PreferredBedtime:    defaultBedtime,
		PreferredWakeupTime: defaultWakeupTime,
		SleepGoalHours:      8,
		IsReminderEnabled:   true,
	}
}
