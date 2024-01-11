package controller

const (
	Success int = iota
	ParamError
	FrpClientError
	ProxyExist
	ProxyNotExist
	ReloadFail
)

const (
	NameKey          = "name"
	OriginalNameKey  = "_original_name"
	SessionName      = "GOSESSION"
	AuthName         = "_PANEL_AUTH"
	LoginUrl         = "/login"
	LoginSuccessUrl  = "/"
	LogoutUrl        = "/logout"
	LogoutSuccessUrl = "/login"
)

type HTTPError struct {
	Code int
	Err  error
}

func (e *HTTPError) Error() string {
	return e.Err.Error()
}

type Common struct {
	Common CommonInfo
}

type CommonInfo struct {
	PluginAddr    string `toml:"plugin_addr"`
	PluginPort    int    `toml:"plugin_port"`
	AdminUser     string `toml:"admin_user"`
	AdminPwd      string `toml:"admin_pwd"`
	AdminKeepTime int    `toml:"admin_keep_time"`
	TlsMode       bool   `toml:"tls_mode"`
	TlsCertFile   string `toml:"tls_cert_file"`
	TlsKeyFile    string `toml:"tls_key_file"`
	DashboardAddr string `toml:"dashboard_addr"`
	DashboardPort int    `toml:"dashboard_port"`
	DashboardUser string `toml:"dashboard_user"`
	DashboardPwd  string `toml:"dashboard_pwd"`
	DashboardTls  bool
}

type OperationResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ProxyResponse struct {
	OperationResponse
	Data any `json:"data"`
}
