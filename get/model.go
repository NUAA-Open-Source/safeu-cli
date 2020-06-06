package get

type ValidationRequest struct {
	Password string `json:"password"`
}

type ValidationResponse struct {
	Token string         `json:"token"`
	Items []ItemDownload `json:"items"`
}

type DownloadRequest struct {
	Full  bool           `json:"full"`
	Items []ItemDownload `json:"items"`
}

type DownloadResponse struct {
	URL string `json:"url"`
}

type MinusDownCountRequest struct {
	Bucket string `json:"bucket"`
	Path   string `json:"path"`
}

type ItemDownload struct {
	OriginalName string `json:"original_name"`
	Protocol     string `json:"protocol"`
	Bucket       string `json:"bucket"`
	Endpoint     string `json:"endpoint"`
	Path         string `json:"path"`
}

// 下载包结构体
type DownloadModel struct {
	Dir          string         // 下载目录（默认为当前目录）
	Filepath     string         // 文件下载地址
	UserRecode   string         // 提取码
	UserPassword string         // 密码
	Csrf         string         // CSRF 口令
	Cookie       string         // Cookie
	Token        string         // 临时认证 Token
	Items        []ItemDownload // 需要下载的文件信息
	URL          string         // 已签名的下载链接
}
