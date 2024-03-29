package controller

import (
	"crypto/tls"
	"fmt"
	"github.com/fatedier/frp/pkg/config"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func trimString(str string) string {
	return strings.TrimSpace(str)
}

func equalIgnoreCase(source string, target string) bool {
	return strings.ToUpper(source) == strings.ToUpper(target)
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
	var filterProxies = make([]v1.TypedProxyConfig, 0)
	for i := range allProxies {
		if equalIgnoreCase(allProxies[i].Type, proxyType) {
			filterProxies = append(filterProxies, allProxies[i])
		}
	}

	return filterProxies, nil
}
