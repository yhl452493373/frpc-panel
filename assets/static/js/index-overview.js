var loadOverview = (function ($) {
    var i18n = {};

    /**
     * get proxy info
     * @param lang {{}} language json
     * @param title page title
     */
    function loadOverview(lang, title) {
        i18n = lang;
        $("#title").text(title);
        $('#content').empty();
        var loading = layui.layer.load();

        $.getJSON('/proxy/api/status').done(function (result) {
            if (result.success) {
                $('#content').html($('#overviewTableTemplate').html());
                renderOverviewTable(JSON.parse(result.data));
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
    function renderOverviewTable(data, proxyType) {
        var dataList = [];
        for (var type in data) {
            var temp = data[type];
            dataList = dataList.concat(temp);
        }

        var $section = $('#content > section');
        var cols = [
            {field: 'name', title: 'Name', sort: true},
            {field: 'type', title: 'Type', width: 100, sort: true},
            {
                field: 'local_addr',
                title: 'Local Address',
                templet: '<span>{{= d.local_addr || "-" }}</span>',
                width: 220,
                sort: true
            },
            {field: 'plugin', title: 'plugin', templet: '<span>{{= d.plugin || "-" }}</span>', sort: true},
            {field: 'remote_addr', title: 'Remote Address', sort: true},
            {field: 'status', title: 'Status', width: 100, sort: true},
            {field: 'err', title: 'Info',templet: '<span>{{= d.err || "-" }}</span>', width: 200}
        ];

        var overviewTable = layui.table.render({
            elem: '#overviewTable',
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
            overviewTable.resize();
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

    return loadOverview;
})(layui.$);