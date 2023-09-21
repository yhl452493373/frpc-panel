var loadProxyInfo = (function ($) {
    var i18n = {};
    //param names in Basic tab
    var basicParamNames = ['name', 'type', 'local_ip', 'local_port', 'custom_domains', 'subdomain', 'remote_port', 'use_encryption', 'use_compression'];

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
                    proxyPopup(type, null);
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
                    proxyPopup(type, data);
                    break;
                case 'remove':
                    // removePopup(data);
                    break;
            }
        });
    }

    /**
     * add proxy popup
     * @param type proxy type
     * @param data {Map<string,object>} proxy data
     */
    function proxyPopup(type, data) {
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
                    var formData = layui.form.val('addProxyTemplate');
                    add(formData, index);
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
            $paramValue.closest('.layui-tab-item').prepend(formItem);
            $paramName.val('');
            $paramValue.val('');

            layui.form.render();

            var tabContent = $paramValue.closest('.layui-tab-content');
            var scrollHeight = tabContent.prop("scrollHeight");
            tabContent.scrollTop(scrollHeight)
        });
        layui.form.on('input-affix(subtraction)', function (obj) {
            var $elem = $(obj.elem);
            $elem.closest('.layui-form-item').remove();
        });
    }

    /**
     * add proxy action
     * @param data {{user:string, token:string, comment:string, enable:boolean, ports:[string|number], domains:[string], subdomains:[string]}} user data
     * @param index popup index
     */
    function add(data, index) {
        var loading = layui.layer.load();
        $.ajax({
            url: '/add',
            type: 'post',
            contentType: 'application/json',
            data: JSON.stringify(data),
            success: function (result) {
                // if (result.success) {
                //     reloadTable();
                //     layui.layer.close(index);
                //     layui.layer.msg(i18n['OperateSuccess'], function (index) {
                //         layui.layer.close(index);
                //     });
                // } else {
                //     errorMsg(result);
                // }
            },
            complete: function () {
                layui.layer.close(loading);
            }
        });
    }

    /**
     * update proxy action
     * @param before {{user:string, token:string, comment:string, enable:boolean, ports:[string|number], domains:[string], subdomains:[string]}} data before update
     * @param after {{user:string, token:string, comment:string, enable:boolean, ports:[string|number], domains:[string], subdomains:[string]}} data after update
     */
    function update(before, after) {
        before.ports.forEach(function (port, index) {
            if (/^\d+$/.test(String(port))) {
                before.ports[index] = parseInt(String(port));
            }
        });
        after.ports.forEach(function (port, index) {
            if (/^\d+$/.test(String(port)) && typeof port === "string") {
                after.ports[index] = parseInt(String(port));
            }
        });
        var loading = layui.layer.load();
        $.ajax({
            url: '/update',
            type: 'post',
            contentType: 'application/json',
            data: JSON.stringify({
                before: before,
                after: after,
            }),
            success: function (result) {
                if (result.success) {
                    layui.layer.msg(i18n['OperateSuccess']);
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

    return loadProxyInfo;
})(layui.$);