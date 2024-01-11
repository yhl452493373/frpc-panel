# frpc-panel(Support FRP >= 0.52.0)

[中文文档](README.md) | [README](README_en.md)

frpc-panel is a client tool of https://github.com/fatedier/frp , it's used to show client info friendly, and manage client proxy info.

frps-panel will run as one single process and manage frpc proxy info with frpc's api.

## Since version 2.0.0,this plugin only support frp version >= v0.52.0

## Features

+ Show frpc basic info
+ Show frpc proxies overview
+ Show frpc proxies list group by proxy type
+ Add proxy in each proxy type
+ Support multiple language,you can translate your own language by add language file in folder `assets/lang/`
+ Automatic darkmode

## Usage

1.add config in frpc's config file '`frpc.toml`:

```toml
webServer.addr = "127.0.0.1"
webServer.port = 7400
webServer.user = "admin"
webServer.password = "admin"
```
or
```toml
[webServer]
addr = "127.0.0.1"
port = 7400
user = "admin"
password = "admin"
```

2.`frpc-panel.toml`:

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

+ `plugin_addr` -- application's listen addr.If you need to visit your frpc-panel with internet, you should change it to `0.0.0.0`
+ `admin_user` -- username used to login
+ `admin_pwd` -- password for `admin_user`
+ `admin_keep_time` -- login session idle time  
+ `tls_mode` -- enable https. If `tls_cert_file` and `tls_key_file` is empty, even this is `true`, it will still run with http
+ `tls_cert_file` -- path of https cert file
+ `tls_key_file` -- path of https cert's key file
+ `dashboard_addr` -- `frpc` ip or domain of your frpc
+ `dashboard_port` -- `admin_port` in your `frpc.ini`
+ `dashboard_user` -- `admin_user` in your `frpc.ini`
+ `dashboard_pwd` -- `admin_pwd` in your `frpc.ini`

3.run with command:
```shell
./frpc-panel -c ./frpc-panel.toml
```

4.Manage your proxies in browser via:`http://127.0.0.1:7300` or `https://127.0.0.1:7300`

## Download

Download frpc-panel binary file from [Release](../../releases).

## Screenshots

![client_info.png](screenshots%2Fclient_info.png)
![client_info_i18n.png](screenshots%2Fclient_info_i18n.png)
![darkmode.png](screenshots%2Fdarkmode.png)
![extra_params.png](screenshots%2Fextra_params.png)
![login.png](screenshots%2Flogin.png)
![new_proxy.png](screenshots%2Fnew_proxy.png)
![proxy_list.png](screenshots%2Fproxy_list.png)
![proxy_overview.png](screenshots%2Fproxy_overview.png)

## Issues & Ideas

If you have any issues or ideas, put it on [issues](https://github.com/yhl452493373/frpc-panel/issues). I will try my best to achieve it.

## Credits

+ [frp](https://github.com/fatedier/frp)
+ [fp-multiuser](https://github.com/gofrp/fp-multiuser)
+ [layui](https://github.com/layui/layui)
+ [layui-theme-dark](https://github.com/Sight-wcg/layui-theme-dark)
+ modified version of [toml](https://github.com/yhl452493373/toml), forked from [toml](https://github.com/flourd/toml)
