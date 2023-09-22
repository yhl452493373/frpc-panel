var loadProxyInfo = (function ($) {
    var i18n = {};
    //param names in Basic tab
    var basicParamNames = ['name', 'type', 'local_ip', 'local_port', 'custom_domains', 'subdomain', 'remote_port', 'use_encryption', 'use_compression'];
    //param need to convert type
    var intParamNames = ['local_port', 'health_check_timeout_s', 'health_check_max_failed', 'health_check_interval_s'];
    var booleanParamNames = ['use_encryption', 'use_compression'];
    var stringArrayParamNames = ['custom_domains', 'locations']
    var mapParamPrefixes = [
        {name: 'metas', prefix: 'meta_'},
        {name: 'plugin_params', prefix: 'plugin_'},
        {name: 'headers', prefix: 'header_'}
    ];

    /**
     * get proxy info
     * @param lang {{}} language json
     * @param title page title
     * @param proxyType proxy type
     */
    function loadProxyInfo(lang, title, proxyType) {
        i18n = lang;
        $("#title").text(title);
        $('#content').empty();
        var loading = layui.layer.load();

        $.getJSON('/proxy/api/config', {
            type: proxyType
        }).done(function (result) {
            if (result.success) {
                $('#content').html($('#proxyListTableTemplate').html());
                renderProxyListTable(result.data, proxyType);
            } else {
                layui.layer.msg(result.message);
            }
        }).always(function () {
            layui.layer.close(loading);
        });
    }

    /**
     * render proxy list table
     * @param data {Map<string,Map<string,string>>} proxy data
     * @param proxyType proxy type
     */
    function renderProxyListTable(data, proxyType) {
        var dataList = [];
        for (var key in data) {
            var temp = data[key];
            temp.name = key;
            temp.local_ip = temp.local_ip || '-';
            temp.local_port = temp.local_port || '-';
            dataList.push(temp);
        }

        var $section = $('#content > section');
        var cols = [
            {type: 'checkbox'},
            {field: 'name', title: 'Name', sort: true},
            {field: 'type', title: 'Type', sort: true},
            {field: 'local_ip', title: 'Local Ip', sort: true},
            {field: 'local_port', title: 'Local Port', sort: true},
            {title: 'Operation', width: 150, toolbar: '#proxyListOperationTemplate'}
        ];

        var proxyListTable = layui.table.render({
            elem: '#proxyListTable',
            height: $section.height(),
            text: {none: i18n['EmptyData']},
            cols: [cols],
            page: navigator.language.indexOf("zh") !== -1,
            toolbar: '#proxyListToolbarTemplate',
            defaultToolbar: false,
            data: dataList,
            initSort: {
                field: 'name',
                type: 'asc'
            }
        });

        window.onresize = function () {
            proxyListTable.resize();
        }

        bindFormEvent(proxyType);
    }

    /**
     * bind event of {{@link layui.form}}
     *
     * @param type proxy type
     */
    function bindFormEvent(type) {
        layui.table.on('toolbar(proxyListTable)', function (obj) {
            var id = obj.config.id;
            var checkStatus = layui.table.checkStatus(id);
            var data = checkStatus.data;

            for (var key in data) {
                data[key] = data[key] === '-' ? '' : data[key];
            }

            switch (obj.event) {
                case 'add':
                    proxyPopup(type, {
                        type: type
                    }, false);
                    break
                case 'remove':
                    // batchRemovePopup(data);
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
                    data.oldName = data.name;
                    proxyPopup(type, data, true);
                    break;
                case 'remove':
                    // removePopup(data);
                    break;
            }
        });
    }

    /**
     * addOrUpdate proxy popup
     * @param type proxy type
     * @param data {Map<string,object>} proxy data
     * @param update update flag. true - update, false - add
     */
    function proxyPopup(type, data, update) {
        var basicData = {};
        var extraData = [];
        if (data != null) {
            var tempData = $.extend(true, {}, data);
            basicParamNames.forEach(function (basicName) {
                if (data.hasOwnProperty(basicName)) {
                    basicData[basicName] = tempData[basicName];
                    delete tempData[basicName];
                }
            });
            for (var key in tempData) {
                extraData.push({
                    name: key,
                    value: tempData[key]
                });
            }
        }
        var html = document.getElementById('addProxyTemplate').innerHTML;
        var content = layui.laytpl(html).render({
            type: type,
            extraData: extraData
        });
        layui.layer.open({
            type: 1,
            title: false,
            skin: 'proxy-popup',
            area: ['500px', '400px'],
            content: content,
            btn: ['Confirm', 'Cancel'],
            btn1: function (index) {
                if (layui.form.validate('#addProxyTemplate')) {
                    var formData = layui.form.val('addProxyForm');
                    var $items = $('#addProxyForm .extra-param-tab-item .extra-param-item');
                    $items.each(function () {
                        var name = $(this).find('input').first().val();
                        var value = $(this).find('input').last().val();
                        formData[name] = value;
                    });

                    addOrUpdate(type, formData, index, update);
                }
            },
            btn2: function (index) {
                layui.layer.close(index);
            },
            success: function (layero, index, that) {
                layui.form.val('addProxyForm', basicData);
                proxyPopupSuccess(layero, index, that, basicData);
            }
        });
    }


    /**
     * repack form data
     * @param formData
     * @returns {*}
     */
    function repackData(formData) {
        mapParamPrefixes.forEach(function (temp) {
            var name = temp.name;
            var prefix = temp.prefix;
            for (var key in formData) {
                if (key !== name && key.startsWith(prefix)) {
                    formData[name] = formData[name] || {};
                    var newKey = key.replace(prefix, '');
                    formData[name][newKey] = formData[key];
                    delete formData[key];
                }
            }
        });
        intParamNames.forEach(function (paramName) {
            if (formData.hasOwnProperty(paramName) && formData[paramName] !== '') {
                formData[paramName] = parseInt(formData[paramName]);
            }
        });
        booleanParamNames.forEach(function (paramName) {
            if (formData.hasOwnProperty(paramName) && formData[paramName] !== '') {
                formData[paramName] = formData[paramName] === 'true';
            }
        });
        stringArrayParamNames.forEach(function (paramName) {
            if (formData.hasOwnProperty(paramName) && formData[paramName] !== '') {
                formData[paramName] = formData[paramName].split(',');
            }
        });
        return formData;
    }

    function proxyPopupSuccess(layero, index, that) {
        layui.form.render(null, 'addProxyForm');
        layui.form.on('input-affix(addition)', function (obj) {
            var $paramValue = $(obj.elem);
            var $paramName = $paramValue.closest('.layui-form-item').find('input');
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

            // var tabContent = $paramValue.closest('.layui-tab-content');
            // var scrollHeight = tabContent.prop("scrollHeight");
            // tabContent.scrollTop(scrollHeight)
        });
        layui.form.on('input-affix(subtraction)', function (obj) {
            var $elem = $(obj.elem);
            $elem.closest('.layui-form-item').remove();
        });
    }

    /**
     * addOrUpdate proxy action
     * @param type proxy type
     * @param data proxy data
     * @param index popup index
     * @param update update flag. true - update, false - add
     */
    function addOrUpdate(type, data, index, update) {
        var loading = layui.layer.load();
        var url = '';
        if (update) {
            url = '/update';
        } else {
            url = '/add?type=' + type;
        }
        $.ajax({
            url: url,
            type: 'post',
            contentType: 'application/json',
            data: JSON.stringify(data),
            success: function (result) {
                if (result.success) {
                    reloadTable();
                    layui.layer.close(index);
                    layui.layer.msg('OperateSuccess', function (index) {
                        layui.layer.close(index);
                    });
                } else {
                    errorMsg(result);
                }
            },
            complete: function () {
                layui.layer.close(loading);
            }
        });
    }

    /**
     * batch remove proxy popup
     * @param data {[{user:string, token:string, comment:string, enable:boolean, ports:[string|number], domains:[string], subdomains:[string]}]} user data list
     */
    function batchRemovePopup(data) {
        if (data.length === 0) {
            layui.layer.msg(i18n['ShouldCheckUser']);
            return;
        }
        layui.layer.confirm(i18n['ConfirmRemoveUser'], {
            title: i18n['OperationConfirm'],
            btn: [i18n['Confirm'], i18n['Cancel']]
        }, function (index) {
            operate(apiType.Remove, data, index);
        });
    }

    /**
     * reload user table
     */
    function reloadTable() {
        // var searchData = layui.form.val('searchForm');
        var searchData = {};
        layui.table.reloadData('tokenTable', {
            where: searchData
        }, true)
    }

    /**
     * show error message popup
     * @param result
     */
    function errorMsg(result) {
        layui.layer.msg(result.message);
        // var reason = i18n['OtherError'];
        // if (result.code === 1)
        //     reason = i18n['ParamError'];
        // else if (result.code === 2)
        //     reason = i18n['SaveError'];
        // else if (result.code === 3)
        //     reason = i18n['FrpServerError'];
        // else if (result.code === 4)
        //     reason = i18n['ProxyExist'];
        // else if (result.code === 5)
        //     reason = i18n['ProxyNotExist'];
        // layui.layer.msg(i18n['OperateFailed'] + ',' + reason)
    }

    return loadProxyInfo;
})(layui.$);