package controller

import (
	"crypto/tls"
	"fmt"
	"github.com/fatedier/frp/pkg/config"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/vaughan0/go-ini"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (c *HandleController) MakeLoginFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		if context.Request.Method == "GET" {
			if c.LoginAuth("", "", context) {
				if context.Request.RequestURI == LoginUrl {
					context.Redirect(http.StatusTemporaryRedirect, LoginSuccessUrl)
				}
				return
			}
			context.HTML(http.StatusOK, "login.html", gin.H{
				"version":             c.Version,
				"FrpcPanel":           ginI18n.MustGetMessage(context, "Frpc Panel"),
				"Username":            ginI18n.MustGetMessage(context, "Username"),
				"Password":            ginI18n.MustGetMessage(context, "Password"),
				"Login":               ginI18n.MustGetMessage(context, "Login"),
				"PleaseInputUsername": ginI18n.MustGetMessage(context, "Please input username"),
				"PleaseInputPassword": ginI18n.MustGetMessage(context, "Please input password"),
			})
		} else if context.Request.Method == "POST" {
			username := context.PostForm("username")
			password := context.PostForm("password")
			if c.LoginAuth(username, password, context) {
				context.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": ginI18n.MustGetMessage(context, "Login success"),
				})
			} else {
				context.JSON(http.StatusOK, gin.H{
					"success": false,
					"message": ginI18n.MustGetMessage(context, "Username or password incorrect"),
				})
			}
		}
	}
}

func (c *HandleController) MakeLogoutFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		ClearAuth(context)
		context.Redirect(http.StatusTemporaryRedirect, LogoutSuccessUrl)
	}
}

func (c *HandleController) MakeIndexFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{
			"version":   c.Version,
			"FrpcPanel": ginI18n.MustGetMessage(context, "Frpc Panel"),
			"showExit":  trimString(c.CommonInfo.AdminUser) != "" && trimString(c.CommonInfo.AdminPwd) != "",
		})
	}
}

func (c *HandleController) MakeLangFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"EmptyData": ginI18n.MustGetMessage(context, "Empty data"),
		})
	}
}

func (c *HandleController) MakeAddProxyFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		proxyType, exist := context.GetQuery("type")
		if !exist {

		}

		proxy, exist := proxyConfTypeMap[proxyType]
		if !exist {

		}

		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "user add success",
		}

		err := context.BindJSON(&proxy)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("user add failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		err = c.reloadFrpc()
		if err != nil {
			response.Success = false
			response.Code = SaveError
			response.Message = fmt.Sprintf("add failed, error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(0, &response)
	}
}

func (c *HandleController) MakeUpdateProxyFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		proxyType, exist := context.GetQuery("type")
		if !exist {

		}

		proxy, exist := proxyConfTypeMap[proxyType]
		if !exist {

		}

		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "user update success",
		}

		err := context.BindJSON(&proxy)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("update failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		err = c.reloadFrpc()
		if err != nil {
			response.Success = false
			response.Code = SaveError
			response.Message = fmt.Sprintf("user update failed, error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(http.StatusOK, &response)
	}
}

func (c *HandleController) MakeRemoveProxyFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		proxyType, exist := context.GetQuery("type")
		if !exist {

		}

		proxy, exist := proxyConfTypeMap[proxyType]
		if !exist {

		}

		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "proxy remove success",
		}

		err := context.BindJSON(&proxy)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("user remove failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		err = c.reloadFrpc()
		if err != nil {
			response.Success = false
			response.Code = SaveError
			response.Message = fmt.Sprintf("user update failed, error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(http.StatusOK, &response)
	}
}

func (c *HandleController) MakeProxyFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		var client *http.Client
		var protocol string

		if c.CommonInfo.DashboardTls {
			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			}
			protocol = "https://"
		} else {
			client = http.DefaultClient
			protocol = "http://"
		}

		res := ProxyResponse{}
		host := c.CommonInfo.DashboardAddr
		port := c.CommonInfo.DashboardPort
		serverApi := context.Param("serverApi")
		requestUrl := protocol + host + ":" + strconv.Itoa(port) + serverApi
		request, _ := http.NewRequest("GET", requestUrl, nil)
		username := c.CommonInfo.DashboardUser
		password := c.CommonInfo.DashboardPwd
		if trimString(username) != "" && trimString(password) != "" {
			request.SetBasicAuth(username, password)
			log.Printf("Proxy to %s", requestUrl)
		}

		response, err := client.Do(request)

		if err != nil {
			res.Code = FrpServerError
			res.Success = false
			res.Message = err.Error()
			log.Print(err)
			context.JSON(http.StatusOK, &res)
			return
		}

		res.Code = response.StatusCode
		body, err := io.ReadAll(response.Body)

		if err != nil {
			res.Success = false
			res.Message = err.Error()
		} else {
			if res.Code == http.StatusOK {
				res.Success = true
				res.Data = string(body)
				res.Message = fmt.Sprintf("Proxy to %s success", requestUrl)
			} else {
				res.Success = false
				if res.Code == http.StatusNotFound {
					res.Message = fmt.Sprintf("Proxy to %s error: url not found", requestUrl)
				} else {
					res.Message = fmt.Sprintf("Proxy to %s error: %s", requestUrl, string(body))
				}
			}
		}
		log.Printf(res.Message)

		if serverApi == "/api/config" {
			proxyType, _ := context.GetQuery("type")
			content := fmt.Sprintf("%s", res.Data)
			configure, err := parseConfigure(content, trimString(proxyType))

			if err != nil {
				res.Success = false
				res.Message = err.Error()
			} else {
				res.Data = configure
			}
		}

		context.JSON(http.StatusOK, &res)
	}
}

func parseConfigure(content, proxyType string) (interface{}, error) {

	if proxyType == "none" {
		common, err := config.UnmarshalClientConfFromIni(content)

		if err != nil {
			return nil, err
		}

		return common, nil
	}

	cfg, err := ini.Load(strings.NewReader(content))
	proxyList := make(map[string]ini.Section)
	for name, section := range cfg {
		if name == "common" {
			continue
		}
		if proxyType != "" && strings.ToLower(section["type"]) != strings.ToLower(proxyType) {
			continue
		}
		proxyList[name] = section
	}

	if err != nil {
		return nil, err
	} else {
		return proxyList, nil
	}
}
