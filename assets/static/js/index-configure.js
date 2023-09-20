var loadCommon = (function ($) {
    var i18n = {};

    /**
     * get client info
     * @param lang {Map<string,string>} language json
     * @param title {string} page title
     */
    function loadCommon(lang, title) {
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
        var html = layui.laytpl($('#commonTemplate').html()).render(data);
        $('#content').html(html);
    }

    return loadCommon;
})(layui.$);