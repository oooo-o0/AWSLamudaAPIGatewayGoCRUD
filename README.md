# 介護基盤ダミーシステム

## 対象ドメインとテーブル一覧

| テーブル名            | 主キー            | GSI例（検索用途）                           | 用途                    |
| ---------------- | -------------- | ------------------------------------ | --------------------- |
| `CarePlan`       | CarePlanID     | UserIDIndex                          | 利用者のケアプラン管理           |
| `User`           | UserID         | RoleIndex                            | 利用者・職員の区別、認証情報        |
| `Claim`          | ClaimID        | FacilityIDIndex / ClaimDateIndex     | 施設単位、請求日検索            |
| `Facility`       | FacilityID     | NameIndex                            | 施設名で検索                |
| `Notification`   | NotificationID | UserIDIndex / StatusIndex            | 通知の未読管理など             |
| `Report`         | ReportID       | ReportTypeIndex / GeneratedDateIndex | レポート検索                |
| `AuthToken`      | TokenID        | UserIDIndex / ExpiryIndex            | 認証トークン管理（JWT/Cognito） |
| `AuditLog`       | LogID          | EventTypeIndex / TimestampIndex      | 操作ログ管理                |
| `MigrationState` | MigrationID    | -                                    | マイグレーション実行履歴管理        |

# CRUD + 一般機能におけるlambda関数
| 種別     | 関数名                           | 説明            |
| ------ | ----------------------------- | ------------- |
| ユーザー   | `UserHandlerFunction`         | 登録、認証、参照など    |
| ケアプラン  | `CarePlanHandlerFunction`     | 作成、更新、削除      |
| 請求     | `ClaimHandlerFunction`        | 請求データ処理       |
| 施設     | `FacilityHandlerFunction`     | 施設情報 CRUD     |
| 通知     | `NotificationHandlerFunction` | メール・システム通知    |
| ログイン   | `AuthLoginFunction`           | ログイン処理（JWT発行） |
| サインアップ | `AuthSignupFunction`          | 新規登録処理        |
| トークン検証 | `AuthVerifyTokenFunction`     | JWT 検証        |
| プロフィール | `UserProfileFunction`         | プロフィール取得・編集   |
| ロール管理  | `UserRoleFunction`            | ユーザー権限設定・取得   |

# バッジにおけるlambda関数
| バッチ名                   | 実行タイミング        | 説明                  |
| ---------------------- | -------------- | ------------------- |
| `DailyBillingBatch`    | 毎日 01:00 JST   | 利用記録を元に請求情報を集計      |
| `WeeklyReportBatch`    | 毎週月曜 03:00 JST | 利用者別レポート集計を PDF 化など |
| `MonthlyCleanupBatch`  | 毎月1日 02:00 JST | 古い通知やログをクリーンアップ     |
| `NotificationDispatch` | 毎日 06:00 JST   | 未読通知を一括送信           |
| `InactiveUserCleanup`     | 毎週日曜 04:00 JST | 30日以上未ログインユーザーの一時停止／通知  |
| `FacilityUsageAggregator` | 毎日 02:30 JST   | 施設ごとの利用統計（稼働率、利用時間）集計   |
| `InvoiceEmailDispatcher`  | 毎月5日 07:00 JST | 月次請求書をメールで自動配信          |
| `CarePlanAutoExpire`      | 毎日 00:10 JST   | 有効期限切れの介護プランを「終了済み」へ更新  |
| `HealthCheckJob`          | 毎時 (:00)       | Lambdaの死活監視、システム健全性ログ記録 |
