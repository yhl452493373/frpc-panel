package controller

import (
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/consts"
	"reflect"
)

const (
	Success int = iota
	ParamError
	SaveError
	FrpServerError
)

const (
	ProxyAdd int = iota
	ProxyUpdate
	ProxyRemove
)

const (
	SessionName      = "GOSESSION"
	AuthName         = "_PANEL_AUTH"
	LoginUrl         = "/login"
	LoginSuccessUrl  = "/"
	LogoutUrl        = "/logout"
	LogoutSuccessUrl = "/login"
)

var (
	proxyConfTypeMap map[string]reflect.Type
)

func init() {
	proxyConfTypeMap = make(map[string]reflect.Type)
	proxyConfTypeMap[consts.TCPProxy] = reflect.TypeOf(config.TCPProxyConf{})
	proxyConfTypeMap[consts.TCPMuxProxy] = reflect.TypeOf(config.TCPMuxProxyConf{})
	proxyConfTypeMap[consts.UDPProxy] = reflect.TypeOf(config.UDPProxyConf{})
	proxyConfTypeMap[consts.HTTPProxy] = reflect.TypeOf(config.HTTPProxyConf{})
	proxyConfTypeMap[consts.HTTPSProxy] = reflect.TypeOf(config.HTTPSProxyConf{})
	proxyConfTypeMap[consts.STCPProxy] = reflect.TypeOf(config.STCPProxyConf{})
	proxyConfTypeMap[consts.XTCPProxy] = reflect.TypeOf(config.XTCPProxyConf{})
	proxyConfTypeMap[consts.SUDPProxy] = reflect.TypeOf(config.SUDPProxyConf{})
}

type HTTPError struct {
	Code int
	Err  error
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

func (e *HTTPError) Error() string {
	return e.Err.Error()
}
