var loadClientInfo = (function ($) {
    var i18n = {};

    /**
     * get client info
     * @param lang {Map<string,string>} language json
     * @param title {string} page title
     */
    function loadClientInfo(lang, title) {
        i18n = lang;
        $("#title").text(title);
        $('#content').empty();
        var loading = layui.layer.load();

        $.getJSON('/proxy/api/config', {
            type: 'none'
        }).done(function (result) {
            if (result.success) {
                renderClientInfo(result.data);
            } else {
                layui.layer.msg(result.message);
            }
        }).always(function () {
            layui.layer.close(loading);
        });
    }

    function renderClientInfo(data) {
        data.transport.tcpMux = i18n[data.transport.tcpMux];
        data.transport.tls.enable = i18n[data.transport.tls.enable];
        var html = layui.laytpl($('#clientInfoTemplate').html()).render(data);
        $('#content').html(html);
    }

    return loadClientInfo;
})(layui.$);