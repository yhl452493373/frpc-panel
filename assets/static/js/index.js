window.clientConfig = {};
window.i18n = {};
(function ($) {
    $(function () {
        function init() {
            var langLoading = layui.layer.load()
            $.getJSON('/lang.json').done(function (lang) {
                window.i18n = lang;
                $.ajaxSetup({
                    error: function (xhr,) {
                        if (xhr.status === 401) {
                            layui.layer.msg(lang['TokenInvalid'], function () {
                                window.location.reload();
                            });
                        }
                    },
                })

                layui.element.on('nav(leftNav)', function (elem) {
                    var id = elem.attr('id');
                    var title = elem.text();
                    if (id === 'clientInfo') {
                        loadClientInfo(title.trim());
                    } else if (id === 'overview') {
                        loadOverview(title.trim());
                    } else if (elem.closest('.layui-nav-item').attr('id') === 'proxies') {
                        if (id != null && id.trim() !== '') {
                            var suffix = elem.closest('.layui-nav-item').children('a').text().trim();
                            loadProxyInfo(title + " " + suffix, id);
                        }
                    }
                });

                $('#leftNav .layui-this > a').click();
            }).always(function () {
                layui.layer.close(langLoading);
            });
        }

        /**
         * add verify rule to layui.form
         */
        function initFormVerifyRule() {
            layui.form.verify({
                proxyName: function (value, elem) {
                    if (value.trim() === '') {
                        var nameI18n = $('#proxyName').closest('.layui-form-item').children('.layui-form-label').text();
                        return nameI18n + i18n['RequireNotEmpty'];
                    }
                },
                localPort: function (value, elem) {
                    if (value !== '' && !/^\d+$/.test(value)) {
                        var nameI18n = $('#localPort').closest('.layui-form-item').children('.layui-form-label').text();
                        return nameI18n + i18n['RequireNumber'];
                    }
                },
                domain: function (value, elem) {
                    var proxyType = $('#proxyType').val().toLowerCase();
                    var $customDomains = $('#customDomains');
                    var customDomains = $customDomains.val();
                    var $subdomain = $('#subdomain');
                    var subdomain = $subdomain.val();
                    if (proxyType === 'http' || proxyType === 'https') {
                        if (customDomains.trim() === '' && subdomain.trim() === '') {
                            var customDomainsNameI18n = $customDomains.closest('.layui-form-item').children('.layui-form-label').text();
                            var subdomainNameI18n = $subdomain.closest('.layui-form-item').children('.layui-form-label').text();
                            return customDomainsNameI18n + i18n['and'] + subdomainNameI18n + i18n['RequireNotAllEmpty'];
                        } else if (customDomains.trim() !== '') {
                            var nameI18n = $customDomains.closest('.layui-form-item').children('.layui-form-label').text();
                            if (!/^\s*\[.*]\s*$/.test(customDomains)) {
                                return nameI18n + i18n['RequireArray'];
                            }
                        }
                    }
                }
            });
        }

        function logout() {
            $.get("/logout", function (result) {
                window.location.reload();
            });
        }

        $(document).on('click.logout', '#logout', function () {
            logout();
        });

        init();
        initFormVerifyRule();
    });
})(layui.$);