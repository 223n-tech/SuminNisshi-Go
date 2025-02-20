# 睡眠日誌DB テーブル定義書

## 1. users（ユーザー）

### 1-1. テーブル定義

ユーザー情報を管理するテーブル

### 1-2. カラム定義

| No. | 物理名              | 論理名               | 型               | NOT NULL | デフォルト        | 備考                        |
| --- | ------------------- | -------------------- | ---------------- | -------- | ----------------- | --------------------------- |
| 1   | id                  | ユーザーID           | int(10) unsigned | YES      | AUTO_INCREMENT    | 主キー                      |
| 2   | email               | メールアドレス       | varchar(255)     | YES      | -                 | ユニーク制約                |
| 3   | display_name        | 表示名               | varchar(100)     | YES      | -                 |                             |
| 4   | password_hash       | パスワード(ハッシュ) | varchar(255)     | YES      | -                 |                             |
| 5   | last_login_datetime | 最終ログイン日時     | datetime         | NO       | NULL              |                             |
| 6   | created             | 作成日時             | datetime         | YES      | CURRENT_TIMESTAMP |                             |
| 7   | modified            | 更新日時             | datetime         | YES      | CURRENT_TIMESTAMP | ON UPDATE CURRENT_TIMESTAMP |
| 8   | deleted             | 削除日時             | datetime         | NO       | NULL              | 論理削除用                  |

### 1-3. インデックス

| No. | インデックス名     | カラム | 種類    | 備考                 |
| --- | ------------------ | ------ | ------- | -------------------- |
| 1   | PRIMARY            | id     | PRIMARY | クラスタインデックス |
| 2   | users_email_UNIQUE | email  | UNIQUE  | ユニーク制約用       |
| 3   | email_idx          | email  | INDEX   | 検索用               |

## 2. sleep_diaries（睡眠日誌）

### 2-1. テーブル定義

睡眠日誌の基本情報を管理するテーブル

### 2-2. カラム定義

| No. | 物理名     | 論理名     | 型               | NOT NULL | デフォルト        | 備考                        |
| --- | ---------- | ---------- | ---------------- | -------- | ----------------- | --------------------------- |
| 1   | id         | 睡眠日誌ID | int(10) unsigned | YES      | AUTO_INCREMENT    | 主キー                      |
| 2   | user_id    | ユーザーID | int(10) unsigned | YES      | -                 | 外部キー（users.id）        |
| 3   | start_date | 記録開始日 | date             | YES      | -                 |                             |
| 4   | end_date   | 記録終了日 | date             | YES      | -                 |                             |
| 5   | diary_name | 日誌名称   | varchar(100)     | YES      | -                 |                             |
| 6   | note       | 備考       | text             | NO       | NULL              |                             |
| 7   | created    | 作成日時   | datetime         | YES      | CURRENT_TIMESTAMP |                             |
| 8   | modified   | 更新日時   | datetime         | YES      | CURRENT_TIMESTAMP | ON UPDATE CURRENT_TIMESTAMP |
| 9   | deleted    | 削除日時   | datetime         | NO       | NULL              | 論理削除用                  |

### 2-3. インデックス

| No. | インデックス名           | カラム     | 種類        | 備考                 |
| --- | ------------------------ | ---------- | ----------- | -------------------- |
| 1   | PRIMARY                  | id         | PRIMARY     | クラスタインデックス |
| 2   | user_id_idx              | user_id    | INDEX       | 外部キー用           |
| 3   | start_date_idx           | start_date | INDEX       | 検索用               |
| 4   | fk_sleep_diaries_user_id | user_id    | FOREIGN KEY | users.id への参照    |

## 3. sleep_records（睡眠記録）

### 3-1. テーブル定義

睡眠記録の詳細を管理するテーブル

### 3-2. カラム定義

| No. | 物理名         | 論理名     | 型               | NOT NULL | デフォルト        | 備考                         |
| --- | -------------- | ---------- | ---------------- | -------- | ----------------- | ---------------------------- |
| 1   | id             | 睡眠記録ID | int(10) unsigned | YES      | AUTO_INCREMENT    | 主キー                       |
| 2   | sleep_diary_id | 睡眠日誌ID | int(10) unsigned | YES      | -                 | 外部キー（sleep_diaries.id） |
| 3   | sleep_state_id | 睡眠状態ID | int(10) unsigned | YES      | -                 | 外部キー（sleep_states.id）  |
| 4   | record_date    | 記録日     | date             | YES      | -                 |                              |
| 5   | time_slot      | 時間枠     | time             | YES      | -                 | 30分単位（00:00-23:30）      |
| 6   | record_type    | 記録種別   | ENUM             | YES      | 'STATE'           | STATE/EVENT/MEAL             |
| 7   | meal_type_id   | 食事種別ID | int(10) unsigned | NO       | NULL              | 外部キー（meal_types.id）    |
| 8   | note           | 備考       | text             | NO       | NULL              |                              |
| 9   | created        | 作成日時   | datetime         | YES      | CURRENT_TIMESTAMP |                              |
| 10  | modified       | 更新日時   | datetime         | YES      | CURRENT_TIMESTAMP | ON UPDATE CURRENT_TIMESTAMP  |
| 11  | deleted        | 削除日時   | datetime         | NO       | NULL              | 論理削除用                   |

### 3-3. インデックス

| No. | インデックス名               | カラム                 | 種類        | 備考                      |
| --- | ---------------------------- | ---------------------- | ----------- | ------------------------- |
| 1   | PRIMARY                      | id                     | PRIMARY     | クラスタインデックス      |
| 2   | sleep_diary_id_idx           | sleep_diary_id         | INDEX       | 外部キー用                |
| 3   | sleep_state_id_idx           | sleep_state_id         | INDEX       | 外部キー用                |
| 4   | record_date_idx              | record_date            | INDEX       | 検索用                    |
| 5   | time_slot_idx                | record_date, time_slot | INDEX       | 時間枠検索用              |
| 6   | fk_sleep_records_sleep_diary | sleep_diary_id         | FOREIGN KEY | sleep_diaries.id への参照 |
| 7   | fk_sleep_records_sleep_state | sleep_state_id         | FOREIGN KEY | sleep_states.id への参照  |
| 8   | fk_sleep_records_meal_type   | meal_type_id           | FOREIGN KEY | meal_types.id への参照    |

## 4. sleep_states（睡眠状態）

### 4-1. テーブル定義

睡眠状態とイベントの種類を管理するマスターテーブル

### 4-2. カラム定義

| No. | 物理名            | 論理名     | 型               | NOT NULL | デフォルト        | 備考                        |
| --- | ----------------- | ---------- | ---------------- | -------- | ----------------- | --------------------------- |
| 1   | id                | 睡眠状態ID | int(10) unsigned | YES      | AUTO_INCREMENT    | 主キー                      |
| 2   | state_name        | 状態名     | varchar(50)      | YES      | -                 |                             |
| 3   | state_code        | 状態コード | varchar(20)      | YES      | -                 | ユニーク制約                |
| 4   | state_description | 状態の説明 | varchar(255)     | NO       | NULL              |                             |
| 5   | display_symbol    | 表示記号   | varchar(10)      | YES      | -                 | 表示用記号（Z、×など）      |
| 6   | display_order     | 表示順     | int(10) unsigned | YES      | 0                 |                             |
| 7   | created           | 作成日時   | datetime         | YES      | CURRENT_TIMESTAMP |                             |
| 8   | modified          | 更新日時   | datetime         | YES      | CURRENT_TIMESTAMP | ON UPDATE CURRENT_TIMESTAMP |
| 9   | deleted           | 削除日時   | datetime         | NO       | NULL              | 論理削除用                  |

### 4-3. インデックス

| No. | インデックス名                 | カラム     | 種類    | 備考                 |
| --- | ------------------------------ | ---------- | ------- | -------------------- |
| 1   | PRIMARY                        | id         | PRIMARY | クラスタインデックス |
| 2   | sleep_states_state_code_UNIQUE | state_code | UNIQUE  | ユニーク制約用       |

### 4-4. 初期データ

```sql
INSERT INTO sleep_states (state_name, state_code, display_symbol, display_order) VALUES
('睡眠中', 'SLEEPING', '■', 1),
('床で覚醒', 'AWAKE_IN_BED', '╱', 2),
('通常覚醒', 'AWAKE', '□', 3),
('強い眠気', 'DROWSINESS', 'Z', 4),
('睡眠薬服用', 'MEDICATION', '×', 5);
```

## 5. meal_types（食事種別）

### 5-1. テーブル定義

食事の種類を管理するマスターテーブル

### 5-2. カラム定義

| No. | 物理名         | 論理名     | 型               | NOT NULL | デフォルト        | 備考                        |
| --- | -------------- | ---------- | ---------------- | -------- | ----------------- | --------------------------- |
| 1   | id             | 食事種別ID | int(10) unsigned | YES      | AUTO_INCREMENT    | 主キー                      |
| 2   | type_name      | 種別名     | varchar(50)      | YES      | -                 |                             |
| 3   | type_code      | 種別コード | varchar(20)      | YES      | -                 | ユニーク制約                |
| 4   | display_symbol | 表示記号   | varchar(10)      | YES      | -                 | 表示用記号（▲、●など）      |
| 5   | display_order  | 表示順     | int(10) unsigned | YES      | 0                 |                             |
| 6   | created        | 作成日時   | datetime         | YES      | CURRENT_TIMESTAMP |                             |
| 7   | modified       | 更新日時   | datetime         | YES      | CURRENT_TIMESTAMP | ON UPDATE CURRENT_TIMESTAMP |
| 8   | deleted        | 削除日時   | datetime         | NO       | NULL              | 論理削除用                  |

### 5-3. インデックス

| No. | インデックス名              | カラム    | 種類    | 備考                 |
| --- | --------------------------- | --------- | ------- | -------------------- |
| 1   | PRIMARY                     | id        | PRIMARY | クラスタインデックス |
| 2   | meal_types_type_code_UNIQUE | type_code | UNIQUE  | ユニーク制約用       |

### 5-4. 初期データ

```sql
INSERT INTO meal_types (type_name, type_code, display_symbol, display_order) VALUES
('朝食', 'BREAKFAST', '▲', 1),
('昼食', 'LUNCH', '●', 2),
('夕食', 'DINNER', '■', 3),
('軽食', 'SNACK', '○', 4);
```

## 6. users_sleep_preferences（ユーザー睡眠設定）

### 6-1. テーブル定義

ユーザーごとの睡眠設定を管理するテーブル

### 6-2. カラム定義

| No. | 物理名                | 論理名               | 型               | NOT NULL | デフォルト        | 備考                        |
| --- | --------------------- | -------------------- | ---------------- | -------- | ----------------- | --------------------------- |
| 1   | id                    | 睡眠設定ID           | int(10) unsigned | YES      | AUTO_INCREMENT    | 主キー                      |
| 2   | user_id               | ユーザーID           | int(10) unsigned | YES      | -                 | 外部キー（users.id）        |
| 3   | preferred_bedtime     | 目標就寝時刻         | time             | YES      | -                 |                             |
| 4   | preferred_wakeup_time | 目標起床時刻         | time             | YES      | -                 |                             |
| 5   | sleep_goal_hours      | 目標睡眠時間         | int(3)           | YES      | 8                 |                             |
| 6   | is_reminder_enabled   | リマインダー設定有無 | boolean          | YES      | true              |                             |
| 7   | created               | 作成日時             | datetime         | YES      | CURRENT_TIMESTAMP |                             |
| 8   | modified              | 更新日時             | datetime         | YES      | CURRENT_TIMESTAMP | ON UPDATE CURRENT_TIMESTAMP |
| 9   | deleted               | 削除日時             | datetime         | NO       | NULL              | 論理削除用                  |

### 6-3. インデックス

| No. | インデックス名                     | カラム  | 種類        | 備考                 |
| --- | ---------------------------------- | ------- | ----------- | -------------------- |
| 1   | PRIMARY                            | id      | PRIMARY     | クラスタインデックス |
| 2   | user_id_idx                        | user_id | INDEX       | 外部キー用           |
| 3   | fk_users_sleep_preferences_user_id | user_id | FOREIGN KEY | users.id への参照    |
