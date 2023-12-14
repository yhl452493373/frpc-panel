package controller

import (
	"crypto/tls"
	"fmt"
	"github.com/fatedier/frp/pkg/config"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/pelletier/go-toml/v2"
	"github.com/vaughan0/go-ini"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func trimString(str string) string {
	return strings.TrimSpace(str)
}

func equalIgnoreCase(source string, target string) bool {
	return strings.ToUpper(source) == strings.ToUpper(target)
}

func sortSectionKeys(object ini.Section) []string {
	var keys []string
	for key := range object {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func serializeSections() []byte {
	var build strings.Builder
	build.WriteString("[common]\n")

	for _, key := range sortSectionKeys(clientCommon) {
		build.WriteString(fmt.Sprintf("%s = %s\n", key, clientCommon[key]))
	}
	build.WriteString("\n")

	sections := Sections{clientProxies}

	for _, sectionInfo := range sections.sort() {
		name := sectionInfo.Name
		build.WriteString(fmt.Sprintf("[%s]\n", name))
		section := sectionInfo.Section

		for _, key := range sortSectionKeys(section) {
			value := section[key]
			if key == NameKey || key == OldNameKey || trimString(value) == "" {
				continue
			}
			build.WriteString(fmt.Sprintf("%s = %s\n", key, value))
		}
		build.WriteString("\n")
	}

	return []byte(build.String())
}

func (c *HandleController) buildRequestUrl(serverApi string) string {
	var protocol string

	if c.CommonInfo.DashboardTls {
		protocol = "https://"
	} else {
		protocol = "http://"
	}

	host := c.CommonInfo.DashboardAddr
	port := c.CommonInfo.DashboardPort

	host, _ = strings.CutPrefix(host, protocol)

	requestUrl := protocol + host + ":" + strconv.Itoa(port) + serverApi

	return requestUrl
}

func (c *HandleController) buildClient() *http.Client {
	var client *http.Client

	if c.CommonInfo.DashboardTls {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	} else {
		client = http.DefaultClient
	}

	return client
}

func (c *HandleController) getClientResponse(request *http.Request, client *http.Client) (*http.Response, error) {
	username := c.CommonInfo.DashboardUser
	password := c.CommonInfo.DashboardPwd
	if trimString(username) != "" && trimString(password) != "" {
		request.SetBasicAuth(username, password)
	}

	response, err := client.Do(request)
	return response, err
}

func (c *HandleController) parseResponse(res *ProxyResponse, response *http.Response) {
	res.Code = response.StatusCode
	body, err := io.ReadAll(response.Body)
	if err != nil {
		res.Success = false
		res.Message = err.Error()
	} else {
		bodyString := string(body)
		url := response.Request.URL
		if res.Code == http.StatusOK {
			res.Success = true
			res.Data = bodyString
			res.Message = fmt.Sprintf("Proxy to %s success", url)
		} else {
			res.Success = false
			if res.Code == http.StatusNotFound {
				res.Message = fmt.Sprintf("Proxy to %s error: url not found", url)
			} else if res.Code == http.StatusBadRequest {
				res.Code = ReloadFail
				res.Message = bodyString
			} else {
				res.Message = fmt.Sprintf("Proxy to %s error: %s", url, bodyString)
			}
		}
	}

	log.Printf(res.Message)
}

type MyProxy struct {
	Toml string `json:"toml"`
	*v1.ProxyBaseConfig
}

func (c *HandleController) parseConfigure(content, proxyType string) (interface{}, error) {
	clientConfig := v1.ClientConfig{}
	err := config.LoadConfigure([]byte(content), &clientConfig)
	if err != nil {
		return nil, err
	}

	if proxyType == "none" {
		return clientConfig, nil
	}

	allProxies := clientConfig.Proxies
	var filterProxies = make([]MyProxy, 0)
	for i := range allProxies {
		if equalIgnoreCase(allProxies[i].Type, proxyType) {
			baseConfig := allProxies[i].GetBaseConfig()
			marshal, _ := toml.Marshal(allProxies[i].ProxyConfigurer)
			proxy := MyProxy{
				string(marshal),
				baseConfig,
			}
			filterProxies = append(filterProxies, proxy)
		}
	}

	var ss = `
Name = 'ssh_random'
Type = 'tcp'
LocalIP = '192.168.31.100'
LocalPort = 22
RemotePort = 0

[Transport]
UseEncryption = false
UseCompression = false
BandwidthLimitMode = ''
ProxyProtocolVersion = ''

[Transport.BandwidthLimit]

[LoadBalancer]
Group = ''
GroupKey = ''

[HealthCheck]
Type = ''
TimeoutSeconds = 0
MaxFailed = 0
IntervalSeconds = 0
Path = ''

[Plugin]
Type = ''

`
	v := v1.ProxyType("tcp")
	r := v1.NewProxyConfigurerByType(v)
	toml.Unmarshal([]byte(ss), r)
	return filterProxies, nil
}
