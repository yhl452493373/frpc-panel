var loadProxyInfo = (function ($) {
    var i18n = {};

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
            data: dataList,
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
        layui.table.on('tool(proxyListTable)', function (obj) {
            var data = obj.data;

            switch (obj.event) {
                case 'update':
                    break;
                // updatePopup(data);
                case 'remove':
                    // removePopup(data);
                    break;
            }
        });
    }

    return loadProxyInfo;
})(layui.$);