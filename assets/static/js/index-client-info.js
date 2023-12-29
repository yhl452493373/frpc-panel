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
                var proxies = [];
                result.data.proxies.forEach(function (proxy){
                    var items = flatJSON(proxy.ProxyConfigurer);
                    proxies.push(expandJSON(items))
                })
                var visitors = [];
                result.data.visitors.forEach(function (visitor){
                    var items = flatJSON(visitor.VisitorConfigurer);
                    visitors.push(expandJSON(items))
                })

                var newD = $.extend({},result.data,true);
                newD.proxies = proxies;
                newD.visitors = visitors;
                console.log(TOML.stringify(newD))
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