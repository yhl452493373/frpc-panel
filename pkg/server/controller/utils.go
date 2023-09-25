package controller

import (
	"crypto/tls"
	"fmt"
	"github.com/fatedier/frp/pkg/config"
	"github.com/vaughan0/go-ini"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func trimString(str string) string {
	return strings.TrimSpace(str)
}

func serializeSectionsToString() []byte {
	var build strings.Builder
	build.WriteString("[common]\n")
	for key, value := range clientCommon {
		build.WriteString(fmt.Sprintf("%s = %s\n", key, value))
	}
	build.WriteString("\n")

	for name, section := range clientProxies {
		build.WriteString(fmt.Sprintf("[%s]\n", name))
		for key, value := range section {
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

func (c *HandleController) parseConfigure(content, proxyType string) (interface{}, error) {
	currentProxies := make(map[string]ini.Section)
	clientProxies = make(map[string]ini.Section)
	common, err := config.UnmarshalClientConfFromIni(content)
	if err != nil {
		return nil, err
	}
	cfg, err := ini.Load(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	for name, section := range cfg {
		if name == "common" {
			clientCommon = section
			continue
		}
		if strings.ToLower(section["type"]) == strings.ToLower(proxyType) {
			currentProxies[name] = section
		}
		clientProxies[name] = section
		delete(clientProxies[name], NameKey)
	}

	if proxyType == "none" {
		return common, nil
	}

	return currentProxies, nil
}
