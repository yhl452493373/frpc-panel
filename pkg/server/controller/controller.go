package controller

import (
	"bytes"
	"fmt"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
			"version":              c.Version,
			"OriginalNameKey":      OriginalNameKey,
			"showExit":             trimString(c.CommonInfo.AdminUser) != "" && trimString(c.CommonInfo.AdminPwd) != "",
			"FrpcPanel":            ginI18n.MustGetMessage(context, "Frpc Panel"),
			"ClientInfo":           ginI18n.MustGetMessage(context, "Client Info"),
			"Overview":             ginI18n.MustGetMessage(context, "Overview"),
			"Proxies":              ginI18n.MustGetMessage(context, "Proxies"),
			"ServerAddress":        ginI18n.MustGetMessage(context, "Server Address"),
			"ServerPort":           ginI18n.MustGetMessage(context, "Server Port"),
			"Protocol":             ginI18n.MustGetMessage(context, "Protocol"),
			"TCPMux":               ginI18n.MustGetMessage(context, "TCP Mux"),
			"User":                 ginI18n.MustGetMessage(context, "User"),
			"UserToken":            ginI18n.MustGetMessage(context, "User Token"),
			"AdminAddress":         ginI18n.MustGetMessage(context, "Admin Address"),
			"AdminPort":            ginI18n.MustGetMessage(context, "Admin Port"),
			"AdminUser":            ginI18n.MustGetMessage(context, "Admin User"),
			"AdminPwd":             ginI18n.MustGetMessage(context, "Admin Pwd"),
			"HeartbeatInterval":    ginI18n.MustGetMessage(context, "Heartbeat Interval"),
			"HeartbeatTimeout":     ginI18n.MustGetMessage(context, "Heartbeat Timeout"),
			"TLSEnable":            ginI18n.MustGetMessage(context, "TLS Enable"),
			"TLSKeyFile":           ginI18n.MustGetMessage(context, "TLS Key File"),
			"TLSCertFile":          ginI18n.MustGetMessage(context, "TLS Cert File"),
			"TLSTrustedCAFile":     ginI18n.MustGetMessage(context, "TLS Trusted CA File"),
			"NewProxy":             ginI18n.MustGetMessage(context, "New Proxy"),
			"RemoveProxy":          ginI18n.MustGetMessage(context, "Remove Proxy"),
			"Update":               ginI18n.MustGetMessage(context, "Update"),
			"Remove":               ginI18n.MustGetMessage(context, "Remove"),
			"Basic":                ginI18n.MustGetMessage(context, "Basic"),
			"Extra":                ginI18n.MustGetMessage(context, "Extra"),
			"ProxyName":            ginI18n.MustGetMessage(context, "Proxy Name"),
			"ProxyNameExample":     ginI18n.MustGetMessage(context, "Proxy Name Example"),
			"LocalIP":              ginI18n.MustGetMessage(context, "Local IP"),
			"LocalIPExample":       ginI18n.MustGetMessage(context, "Local IP Example"),
			"LocalPort":            ginI18n.MustGetMessage(context, "Local Port"),
			"LocalPortExample":     ginI18n.MustGetMessage(context, "Local Port Example"),
			"RemotePort":           ginI18n.MustGetMessage(context, "Remote Port"),
			"RemotePortExample":    ginI18n.MustGetMessage(context, "Remote Port Example"),
			"CustomDomains":        ginI18n.MustGetMessage(context, "Custom Domains"),
			"CustomDomainsExample": ginI18n.MustGetMessage(context, "Custom Domains Example"),
			"Subdomain":            ginI18n.MustGetMessage(context, "Subdomain"),
			"UseEncryption":        ginI18n.MustGetMessage(context, "Use Encryption"),
			"true":                 ginI18n.MustGetMessage(context, "true"),
			"UseCompression":       ginI18n.MustGetMessage(context, "Use Compression"),
			"ParamName":            ginI18n.MustGetMessage(context, "Param Name"),
			"ParamValue":           ginI18n.MustGetMessage(context, "Param Value"),
		})
	}
}

func (c *HandleController) MakeLangFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"Proxies":            ginI18n.MustGetMessage(context, "Proxies"),
			"EmptyData":          ginI18n.MustGetMessage(context, "Empty data"),
			"true":               ginI18n.MustGetMessage(context, "true"),
			"false":              ginI18n.MustGetMessage(context, "false"),
			"Name":               ginI18n.MustGetMessage(context, "Name"),
			"Type":               ginI18n.MustGetMessage(context, "Type"),
			"LocalAddress":       ginI18n.MustGetMessage(context, "Local Address"),
			"Plugin":             ginI18n.MustGetMessage(context, "Plugin"),
			"RemoteAddress":      ginI18n.MustGetMessage(context, "Remote Address"),
			"Status":             ginI18n.MustGetMessage(context, "Status"),
			"Info":               ginI18n.MustGetMessage(context, "Info"),
			"running":            ginI18n.MustGetMessage(context, "running"),
			"start error":        ginI18n.MustGetMessage(context, "start error"),
			"new":                ginI18n.MustGetMessage(context, "new"),
			"LocalIP":            ginI18n.MustGetMessage(context, "Local IP"),
			"LocalPort":          ginI18n.MustGetMessage(context, "Local Port"),
			"RemotePort":         ginI18n.MustGetMessage(context, "Remote Port"),
			"UseEncryption":      ginI18n.MustGetMessage(context, "Use Encryption"),
			"UseCompression":     ginI18n.MustGetMessage(context, "Use Compression"),
			"CustomDomains":      ginI18n.MustGetMessage(context, "Custom Domains"),
			"Subdomain":          ginI18n.MustGetMessage(context, "Subdomain"),
			"Operation":          ginI18n.MustGetMessage(context, "Operation"),
			"Confirm":            ginI18n.MustGetMessage(context, "Confirm"),
			"Cancel":             ginI18n.MustGetMessage(context, "Cancel"),
			"TokenInvalid":       ginI18n.MustGetMessage(context, "Token Invalid"),
			"OperateSuccess":     ginI18n.MustGetMessage(context, "Operate Success"),
			"OperateFailed":      ginI18n.MustGetMessage(context, "Operate Failed"),
			"ShouldCheckProxy":   ginI18n.MustGetMessage(context, "Should Check Proxy"),
			"ConfirmRemoveProxy": ginI18n.MustGetMessage(context, "Confirm Remove Proxy"),
			"OperationConfirm":   ginI18n.MustGetMessage(context, "Operation Confirm"),
			"ParamError":         ginI18n.MustGetMessage(context, "Param Error"),
			"FrpClientError":     ginI18n.MustGetMessage(context, "Frp Client Error"),
			"ProxyExist":         ginI18n.MustGetMessage(context, "Proxy Exist"),
			"ProxyNotExist":      ginI18n.MustGetMessage(context, "Proxy Not Exist"),
			"ClientTips":         ginI18n.MustGetMessage(context, "Client Tips"),
			"and":                ginI18n.MustGetMessage(context, "and"),
			"RequireNotAllEmpty": ginI18n.MustGetMessage(context, "Require Not All Empty"),
			"RequireNotEmpty":    ginI18n.MustGetMessage(context, "Require Not Empty"),
			"RequireNumber":      ginI18n.MustGetMessage(context, "Require Number"),
			"RequireBoolean":     ginI18n.MustGetMessage(context, "Require Boolean"),
			"RequireArray":       ginI18n.MustGetMessage(context, "Require Array"),
			"Total":              ginI18n.MustGetMessage(context, "Total"),
			"Items":              ginI18n.MustGetMessage(context, "Items"),
			"Goto":               ginI18n.MustGetMessage(context, "Go to"),
			"PerPage":            ginI18n.MustGetMessage(context, "Per Page"),
		})
	}
}

func (c *HandleController) MakeUpdateProxyFunc() func(context *gin.Context) {
	return func(context *gin.Context) {

		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "proxy update success",
		}

		data, err := context.GetRawData()
		if err != nil {
			response.Success = false
			response.Code = http.StatusNoContent
			response.Message = fmt.Sprintf("update failed, error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		res := c.UpdateFrpcConfig(data)
		if !res.Success {
			response.Success = false
			response.Code = res.Code
			response.Message = fmt.Sprintf("update failed, error : %v", res.Message)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(http.StatusOK, &response)
	}
}

func (c *HandleController) MakeProxyFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		res := ProxyResponse{}
		serverApi := context.Param("serverApi")
		requestUrl := c.buildRequestUrl(serverApi)
		request, _ := http.NewRequest("GET", requestUrl, nil)
		response, err := c.getClientResponse(request, c.buildClient())

		if err != nil {
			res.Code = FrpClientError
			res.Success = false
			res.Message = err.Error()
			log.Print(err)
			context.JSON(http.StatusOK, &res)
			return
		}

		c.parseResponse(&res, response)

		if serverApi == "/api/config" {
			proxyType, _ := context.GetQuery("type")
			content := fmt.Sprintf("%s", res.Data)
			configure, err := c.parseConfigure(content, trimString(proxyType))

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

func (c *HandleController) UpdateFrpcConfig(tomlStr []byte) ProxyResponse {
	res := ProxyResponse{}

	requestUrl := c.buildRequestUrl("/api/config")
	request, _ := http.NewRequest("PUT", requestUrl, bytes.NewReader(tomlStr))
	response, err := c.getClientResponse(request, c.buildClient())

	if err != nil {
		res.Code = FrpClientError
		res.Success = false
		res.Message = err.Error()
		log.Print(err)
		return res
	}

	c.parseResponse(&res, response)

	if res.Success {
		c.ReloadFrpcConfig(&res)
	}
	return res
}

func (c *HandleController) ReloadFrpcConfig(res *ProxyResponse) {
	requestUrl := c.buildRequestUrl("/api/reload")
	request, _ := http.NewRequest("GET", requestUrl, nil)
	response, err := c.getClientResponse(request, c.buildClient())

	if err != nil {
		res.Code = FrpClientError
		res.Success = false
		res.Message = err.Error()
		log.Print(err)
		return
	}

	c.parseResponse(res, response)
}
