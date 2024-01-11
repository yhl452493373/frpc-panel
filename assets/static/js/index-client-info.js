var loadClientInfo = (function ($) {
    /**
     * get client info
     * @param title {string} page title
     */
    function loadClientInfo( title) {
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