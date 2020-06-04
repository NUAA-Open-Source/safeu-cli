package get

// 下载包结构体
type ItemDownload struct {
	userRecode   string // 提取码
	userPassword string // 密码
	csrf         string // CSRF 口令
	cookie       string // Cookie
	token        string // 临时认证 Token
}
