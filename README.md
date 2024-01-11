# frpc-panel(支持 FRP >= 0.52.0)

[中文文档](README.md) | [README](README_en.md)

frpc-panel 是 https://github.com/fatedier/frp 的一个客户端工具，用于更好的展示客户端信息，以及管理客户端代理信息。

frps-panel 会以一个单独的进程运行，通过后台调用frpc的接口实现对frpc的操作。

## 从版本2.0.0开始，本插件只支持版本号大于等于v0.52.0的frp

## 功能

+ 展示客户端基础配置信息
+ 展示客户端代理概览
+ 分类展示客户端代理连接
+ 在各个分类下添加相应代理
+ 国际化
+ 自动深色模式

## 使用方法

1、frpc的配置文件`frpc.toml`中，增加如下内容：

```toml
webServer.addr = "127.0.0.1"
webServer.port = 7400
webServer.user = "admin"
webServer.password = "admin"
```
或
```toml
[webServer]
addr = "127.0.0.1"
port = 7400
user = "admin"
password = "admin"
```

2、frpc-panel的配置文件`frpc-panel.toml`的配置如下：

```toml
# basic options
[common]
# frps panel config info
plugin_addr = "127.0.0.1"
plugin_port = 7300
#admin_user = "admin"
#admin_pwd = "admin"
# specified login state keep time in secends
admin_keep_time = 0

# enable tls
tls_mode = false
#tls_cert_file = "cert.crt"
#tls_key_file = "cert.key"

# frpc dashboard info
dashboard_addr = "127.0.0.1"
dashboard_port = 7400
dashboard_user = "admin"
dashboard_pwd = "admin"
```

+ `plugin_addr`指定监听地址。如需要外网访问，则配成`0.0.0.0`
+ `admin_user`指定登录时的账户。如果需要鉴权登录，则去掉`admin_user`和`admin_pwd`前面的`#`
+ `admin_pwd`指定登录时的密码。如果需要鉴权登录，则去掉`admin_user`和`admin_pwd`前面的`#`
+ `admin_keep_time`指定登陆后session的空闲时间，单位为秒。0表示完全关闭浏览器后登录失效；大于0表示空闲超过此时间后，登录失效
+ `tls_mode`启用https。如果未配置`tls_cert_file`和`tls_key_file`，即使这里为`true`，仍然以http的方式运行
+ `tls_cert_file`https的证书文件路径
+ `tls_key_file`https的证书的密钥文件路径
+ `dashboard_addr`为`frpc`客户端地址。如果`frpc-panel`和`frpc`在同一机器，则可以配置为`127.0.0.1`，否则为对应的ip或域名。如果frpc的地址为https，则填写`https://xxx.yyy.zzz`即可
+ `dashboard_port`为`frpc`客户端管理端口`admin_port`
+ `dashboard_user`为`frpc`客户端管理账户`admin_user`
+ `dashboard_pwd`为`frpc`客户端管理密码`admin_pwd`

3、通过在控制台或终端中执行`./frpc-panel -c ./frpc-panel.toml`启动

4、浏览器中输入`http://127.0.0.1:7300`或`https://127.0.0.1:7300`访问面板

## 下载

通过 [Release](../../releases) 页面下载对应系统版本的二进制文件到本地。

## 预览截图

![client_info.png](screenshots%2Fclient_info.png)
![client_info_i18n.png](screenshots%2Fclient_info_i18n.png)
![darkmode.png](screenshots%2Fdarkmode.png)
![extra_params.png](screenshots%2Fextra_params.png)
![login.png](screenshots%2Flogin.png)
![new_proxy.png](screenshots%2Fnew_proxy.png)
![proxy_list.png](screenshots%2Fproxy_list.png)
![proxy_overview.png](screenshots%2Fproxy_overview.png)

如果使用中有问题或者有其他想法，在[issues](https://github.com/yhl452493373/frpc-panel/issues)上提出来。 如果我能搞定的话，我尽量搞。

## 致谢

+ [frp](https://github.com/fatedier/frp)
+ [fp-multiuser](https://github.com/gofrp/fp-multiuser)
+ [layui](https://github.com/layui/layui)
+ [layui-theme-dark](https://github.com/Sight-wcg/layui-theme-dark)
+ 修改过的JavaScript版的[toml](https://github.com/yhl452493373/toml),原始地址[toml](https://github.com/flourd/toml)
