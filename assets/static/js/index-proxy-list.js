var loadProxyInfo = (function ($) {
    var currentProxyType, currentTitle;

    /**
     * get proxy info
     * @param title page title
     * @param proxyType proxy type
     */
    function loadProxyInfo(title, proxyType) {
        if (title != null)
            currentTitle = title;
        if (proxyType != null)
            currentProxyType = proxyType.toLowerCase();
        $("#title").text(currentTitle);
        $('#content').empty();
        var loading = layui.layer.load();

        $.getJSON('/proxy/api/config', {
            type: currentProxyType
        }).done(function (result) {
            if (result.success) {
                $('#content').html($('#proxyListTableTemplate').html());
                renderProxyListTable(result.data);
                loadFrpcConfig();
            } else {
                layui.layer.msg(result.message);
            }
        }).always(function () {
            layui.layer.close(loading);
        });
    }

    function loadFrpcConfig() {
        $.getJSON('/proxy/api/config', {
            type: 'none'
        }).done(function (result) {
            if (result.success) {
                var proxies = [];
                result.data.proxies.forEach(function (proxy) {
                    var items = flatJSON(proxy['ProxyConfigurer']);
                    proxies.push(expandJSON(items))
                })
                var visitors = [];
                result.data.visitors.forEach(function (visitor) {
                    var items = flatJSON(visitor['VisitorConfigurer']);
                    visitors.push(expandJSON(items))
                });

                window.clientConfig = $.extend(true, {}, result.data);
                window.clientConfig.proxies = proxies;
                window.clientConfig.visitors = visitors;
            } else {
                window.clientConfig = {};
                layui.layer.msg(result.message);
            }
        });
    }

    /**
     * render proxy list table
     * @param data {[Map<string,string>]} proxy data
     */
    function renderProxyListTable(data) {
        var proxies = [];
        data.forEach(function (temp) {
            var proxy = temp.ProxyConfigurer;
            basicParams.forEach(function (basicParam) {
                var name = basicParam.name;
                var defaultValue = basicParam.defaultValue;
                var value = null;
                try {
                    value = eval('proxy.' + name);
                    if (value == null) {
                        value = defaultValue;
                    }
                } catch (e) {
                    value = defaultValue;
                }
                eval('proxy.' + name + ' = value');

            });
            proxies.push(proxy);
        });

        var $section = $('#content > section');
        var cols = [
            {type: 'checkbox'},
            {field: 'name', title: i18n['Name'], sort: true},
            {field: 'type', title: i18n['Type'], width: 110, sort: true},
            {field: 'localIP', title: i18n['LocalIP'], width: 150, sort: true},
            {field: 'localPort', title: i18n['LocalPort'], width: 120, sort: true},
        ];

        if (currentProxyType === 'tcp' || currentProxyType === 'udp') {
            cols.push({field: 'remotePort', title: i18n['RemotePort'], width: 130, sort: true});
        } else if (currentProxyType === 'http' || currentProxyType === 'https') {
            cols.push({field: 'customDomains', title: i18n['CustomDomains'], sort: true});
            cols.push({field: 'subdomain', title: i18n['Subdomain'], width: 150, sort: true});
        }

        cols.push({
            field: 'useEncryption', title: i18n['UseEncryption'], width: 170, templet: function (d) {
                return i18n[d.transport.useEncryption]
            }, sort: true
        });
        cols.push({
            field: 'useCompression', title: i18n['UseCompression'], width: 170, templet: function (d) {
                return i18n[d.transport.useCompression]
            }, sort: true
        });
        cols.push({title: i18n['Operation'], width: 150, toolbar: '#proxyListOperationTemplate'});

        var proxyListTable = layui.table.render({
            elem: '#proxyListTable',
            height: $section.height(),
            text: {none: i18n['EmptyData']},
            cols: [cols],
            page: {
                limitTemplet: function (item) {
                    return item + i18n['PerPage'];
                },
                skipText: [i18n['Goto'], '', i18n['Confirm']],
                countText: [i18n['Total'], i18n['Items']]
            },
            toolbar: '#proxyListToolbarTemplate',
            defaultToolbar: false,
            data: proxies,
            initSort: {
                field: 'name',
                type: 'asc'
            }
        });

        window.onresize = function () {
            proxyListTable.resize();
        }

        bindFormEvent();
    }

    /**
     * bind event of {{@link layui.form}}
     */
    function bindFormEvent() {
        layui.table.on('toolbar(proxyListTable)', function (obj) {
            var id = obj.config.id;
            var checkStatus = layui.table.checkStatus(id);
            var data = checkStatus.data;

            for (var key in data) {
                data[key] = data[key] === '-' ? '' : data[key];
            }

            switch (obj.event) {
                case 'add':
                    proxyPopup({
                        type: currentProxyType
                    }, false);
                    break
                case 'remove':
                    batchRemovePopup(data);
                    break
            }
        });

        layui.table.on('tool(proxyListTable)', function (obj) {
            var data = obj.data;

            for (var key in data) {
                data[key] = data[key] === '-' ? '' : data[key];
            }

            switch (obj.event) {
                case 'update':
                    proxyPopup(data, true);
                    break;
                case 'remove':
                    batchRemovePopup([data]);
                    break;
            }
        });
    }

    /**
     * addOrUpdate proxy popup
     * @param data {Map<string,object>} proxy data
     * @param update update flag. true - update, false - add
     */
    function proxyPopup(data, update) {
        var basicData = {};
        var extraData = [];
        if (data != null) {
            var tempData = $.extend(true, {}, data);

            basicParams.forEach(function (basicName) {
                var name = basicName.name;
                if (name.indexOf('.') !== -1) {
                    var keys = name.split('.');
                    expandJSONKeys(tempData, keys, null);
                    expandJSONKeys(basicData, keys, basicName.defaultValue);
                }

                eval('basicData.' + name + ' = ' + 'tempData.' + name);
                eval('delete tempData.' + name)
            });

            var flatted = flatJSON(tempData);
            for (var key in flatted) {
                var value = flatted[key];
                if (value == null || value === '')
                    continue;
                extraData.push({
                    name: key,
                    value: value
                });
            }
        }

        var html = document.getElementById('proxyFormTemplate').innerHTML;
        var content = layui.laytpl(html).render({
            type: currentProxyType,
            extraData: extraData
        });
        layui.layer.open({
            type: 1,
            title: false,
            skin: 'proxy-popup',
            area: ['550px', '400px'],
            content: content,
            btn: [i18n['Confirm'], i18n['Cancel']],
            btn1: function (index) {
                if (layui.form.validate('#proxyForm')) {
                    var formData = layui.form.val('proxyForm');
                    var $items = $('#proxyForm .extra-param-tab-item .extra-param-item');
                    $items.each(function () {
                        var name = $(this).find('input').first().val();
                        var value = $(this).find('input').last().val();
                        formData[name] = value;
                    });
                    addOrUpdate(formData, index, update);
                }
            },
            btn2: function (index) {
                layui.layer.close(index);
            },
            success: function (layero, index, that) {
                //get and set old name for update form
                var originalNameKey = layero.find('#originalNameKey').attr('name');
                basicData[originalNameKey] = basicData.name;
                layui.form.val('proxyForm', flatJSON(basicData));
                proxyPopupSuccess();
            }
        });
    }

    function proxyPopupSuccess() {
        layui.form.render(null, 'proxyForm');
        layui.form.on('input-affix(addition)', function (obj) {
            var $paramValue = $(obj.elem);
            var $paramName = $paramValue.closest('.layui-form-item').find('input[type=text]');
            var name = $paramName.first().val();
            var value = $paramValue.val();
            var html = document.getElementById('extraParamAddedTemplate').innerHTML;
            var formItem = layui.laytpl(html).render({
                name: name,
                value: value
            });
            $paramValue.closest('.layui-tab-item').append(formItem);
            $paramName.val('');
            $paramValue.val('');

            layui.form.render();
        });
        layui.form.on('input-affix(subtraction)', function (obj) {
            var $elem = $(obj.elem);
            $elem.closest('.layui-form-item').remove();
        });
    }

    /**
     * addOrUpdate proxy action
     * @param data proxy data
     * @param index popup index
     * @param update update flag. true - update, false - add
     */
    function addOrUpdate(data, index, update) {
        try {
            var originalNameKey = $('#originalNameKey').attr('name');
            var proxies = clientConfig.proxies;
            if (update) {
                for (var i = 0; i < proxies.length; i++) {
                    if (data[originalNameKey] === proxies[i].name) {
                        delete data[originalNameKey];
                        proxies[i] = expandJSON(data);
                    }
                }
            } else {
                proxies.push(expandJSON(data));
            }
        } catch (e) {
            layui.layer.msg(e.message);
            return;
        }

        updateFrpcConfig(index);
    }

    /**
     * batch remove proxy popup
     * @param data {[Map<string,string>]} proxy data list
     */
    function batchRemovePopup(data) {
        if (data.length === 0) {
            layui.layer.msg(i18n['ShouldCheckProxy']);
            return;
        }
        layui.layer.confirm(i18n['ConfirmRemoveProxy'], {
            title: i18n['OperationConfirm'],
            btn: [i18n['Confirm'], i18n['Cancel']]
        }, function (index) {
            layui.layer.close(index);
            var proxies = clientConfig.proxies;
            for (var i = 0; i < data.length; i++) {
                proxies.forEach(function (proxy, j) {
                    if (data[i].name === proxy.name) {
                        proxies.splice(j, 1);
                    }
                });
            }

            updateFrpcConfig(index);
        });
    }

    /**
     * update frpc's config
     * @param index popup index
     */
    function updateFrpcConfig(index) {
        var loading = layui.layer.load();
        var content = TOML.stringify(clientConfig);
        $.ajax({
            url: '/update',
            type: 'post',
            contentType: 'text/plain',
            data: content,
            success: function (result) {
                if (result.success) {
                    layui.layer.close(index);
                    reloadTable();
                    layui.layer.msg(i18n['OperateSuccess']);
                } else {
                    errorMsg(result);
                    if (result.code === 5) {
                        layui.layer.close(index);
                        reloadTable();
                    }
                }
            },
            complete: function () {
                layui.layer.close(loading);
            }
        });
    }

    /**
     * reload proxy table
     */
    function reloadTable() {
        loadProxyInfo(null, null);
    }

    /**
     * show error message popup
     * @param result
     */
    function errorMsg(result) {
        var reason = i18n['OtherError'];
        if (result.code === 1)
            reason = i18n['ParamError'];
        else if (result.code === 2)
            reason = i18n['FrpClientError'];
        else if (result.code === 3)
            reason = i18n['ProxyExist'];
        else if (result.code === 4)
            reason = i18n['ProxyNotExist'];
        if (result.code === 5) {
            layui.layer.alert(result.message, {
                title: i18n['ClientTips'],
                maxWidth: 350,
                btn: [i18n['Confirm']]
            });
        } else {
            layui.layer.msg(i18n['OperateFailed'] + ',' + reason);
        }
    }

    return loadProxyInfo;
})(layui.$);