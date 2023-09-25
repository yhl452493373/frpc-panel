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
                renderCommonInfo(result.data);
            } else {
                layui.layer.msg(result.message);
            }
        }).always(function () {
            layui.layer.close(loading);
        });
    }

    function renderCommonInfo(data) {
        data.tcp_mux = i18n[data.tcp_mux];
        data.tls_enable = i18n[data.tls_enable];
        var html = layui.laytpl($('#clientInfoTemplate').html()).render(data);
        $('#content').html(html);
    }

    return loadClientInfo;
})(layui.$);