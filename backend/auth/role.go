package auth

const (
	RoleAdmin  = "admin"
	RoleUser   = "user"
	RoleViewer = "viewer"
)

// 指定ロールがアクセス権を持つかチェック
func HasRole(userRole, requiredRole string) bool {
	// 簡単に文字列比較例。必要に応じて複雑な権限階層も実装可
	return userRole == requiredRole || userRole == RoleAdmin
}
