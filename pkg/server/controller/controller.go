package controller

import (
	"crypto/tls"
	"fmt"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"sort"
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
		context.JSON(http.StatusOK, gin.H{})
	}
}

func (c *HandleController) MakeQueryTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {

		search := TokenSearch{}
		search.Limit = 0

		err := context.BindQuery(&search)
		if err != nil {
			return
		}

		var tokenList []TokenInfo
		for _, tokenInfo := range c.Tokens {
			tokenList = append(tokenList, tokenInfo)
		}
		sort.Slice(tokenList, func(i, j int) bool {
			return strings.Compare(tokenList[i].User, tokenList[j].User) < 0
		})

		var filtered []TokenInfo
		for _, tokenInfo := range tokenList {
			if filter(tokenInfo, search.TokenInfo) {
				filtered = append(filtered, tokenInfo)
			}
		}
		if filtered == nil {
			filtered = []TokenInfo{}
		}

		count := len(filtered)
		if search.Limit > 0 {
			start := max((search.Page-1)*search.Limit, 0)
			end := min(search.Page*search.Limit, len(filtered))
			filtered = filtered[start:end]
		}

		context.JSON(http.StatusOK, &TokenResponse{
			Code:  0,
			Msg:   "query Tokens success",
			Count: count,
			Data:  filtered,
		})
	}
}

func (c *HandleController) MakeAddTokenFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		info := TokenInfo{
			Enable: true,
		}
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "user add success",
		}
		err := context.BindJSON(&info)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("user add failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		result := c.verifyToken(info, TOKEN_ADD)

		if !result.Success {
			context.JSON(http.StatusOK, &result)
			return
		}

		info.Comment = cleanString(info.Comment)
		info.Ports = cleanPorts(info.Ports)
		info.Domains = cleanStrings(info.Domains)
		info.Subdomains = cleanStrings(info.Subdomains)

		c.Tokens[info.User] = info

		err = c.saveToken()
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

func (c *HandleController) MakeUpdateTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "user update success",
		}
		update := TokenUpdate{}
		err := context.BindJSON(&update)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("update failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		before := update.Before
		after := update.After

		if before.User != after.User {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("update failed, user should be same : before -> %v, after -> %v", before.User, after.User)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		result := c.verifyToken(after, TOKEN_UPDATE)

		if !result.Success {
			context.JSON(http.StatusOK, &result)
			return
		}

		after.Comment = cleanString(after.Comment)
		after.Ports = cleanPorts(after.Ports)
		after.Domains = cleanStrings(after.Domains)
		after.Subdomains = cleanStrings(after.Subdomains)

		c.Tokens[after.User] = after

		err = c.saveToken()
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

func (c *HandleController) MakeRemoveTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "user remove success",
		}
		remove := TokenRemove{}
		err := context.BindJSON(&remove)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("user remove failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		for _, user := range remove.Users {
			result := c.verifyToken(user, TOKEN_REMOVE)

			if !result.Success {
				context.JSON(http.StatusOK, &result)
				return
			}
		}

		for _, user := range remove.Users {
			delete(c.Tokens, user.User)
		}

		err = c.saveToken()
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

func (c *HandleController) MakeDisableTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "remove success",
		}
		disable := TokenDisable{}
		err := context.BindJSON(&disable)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("disable failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		for _, user := range disable.Users {
			result := c.verifyToken(user, TOKEN_DISABLE)

			if !result.Success {
				context.JSON(http.StatusOK, &result)
				return
			}
		}

		for _, user := range disable.Users {
			token := c.Tokens[user.User]
			token.Enable = false
			c.Tokens[user.User] = token
		}

		err = c.saveToken()

		if err != nil {
			response.Success = false
			response.Code = SaveError
			response.Message = fmt.Sprintf("disable failed, error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		context.JSON(http.StatusOK, &response)
	}
}

func (c *HandleController) MakeEnableTokensFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		response := OperationResponse{
			Success: true,
			Code:    Success,
			Message: "remove success",
		}
		enable := TokenEnable{}
		err := context.BindJSON(&enable)
		if err != nil {
			response.Success = false
			response.Code = ParamError
			response.Message = fmt.Sprintf("enable failed, param error : %v", err)
			log.Printf(response.Message)
			context.JSON(http.StatusOK, &response)
			return
		}

		for _, user := range enable.Users {
			result := c.verifyToken(user, TOKEN_ENABLE)

			if !result.Success {
				context.JSON(http.StatusOK, &result)
				return
			}
		}

		for _, user := range enable.Users {
			token := c.Tokens[user.User]
			token.Enable = true
			c.Tokens[user.User] = token
		}

		err = c.saveToken()

		if err != nil {
			log.Printf("enable failed, error : %v", err)
			response.Success = false
			response.Code = SaveError
			response.Message = "enable failed"
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
		requestUrl := protocol + host + ":" + strconv.Itoa(port) + context.Param("serverApi")
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
				res.Message = "Proxy to " + requestUrl + " success"
			} else {
				res.Success = false
				res.Message = "Proxy to " + requestUrl + " error: " + string(body)
			}
		}
		log.Printf(res.Message)
		context.JSON(http.StatusOK, &res)
	}
}
