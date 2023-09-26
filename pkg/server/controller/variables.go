package controller

import (
	"github.com/vaughan0/go-ini"
	"sort"
	"strings"
)

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
	OldNameKey       = "_old_name"
	SessionName      = "GOSESSION"
	AuthName         = "_PANEL_AUTH"
	LoginUrl         = "/login"
	LoginSuccessUrl  = "/"
	LogoutUrl        = "/logout"
	LogoutSuccessUrl = "/login"
)

var (
	clientCommon  ini.Section
	clientProxies map[string]ini.Section
)

func init() {
	clientCommon = ini.Section{}
	clientProxies = make(map[string]ini.Section)
}

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

type ClientProxies struct {
	Proxy ini.Section `json:"proxy"`
}

type SectionInfo struct {
	Name    string
	Section ini.Section
}

type Sections struct {
	sections map[string]ini.Section
}

func (s *Sections) sort() []SectionInfo {
	sectionInfos := make([]SectionInfo, 0)

	for key, value := range s.sections {
		sectionInfos = append(sectionInfos, SectionInfo{Name: key, Section: value})
	}

	sort.Slice(sectionInfos, func(i, j int) bool {
		typeCompare := strings.Compare(sectionInfos[i].Section["type"], sectionInfos[j].Section["type"])
		if typeCompare == -1 {
			return true
		} else if typeCompare == 0 {
			nameCompare := strings.Compare(sectionInfos[i].Name, sectionInfos[j].Name)
			return nameCompare == -1
		} else {
			return false
		}
	})

	return sectionInfos
}
