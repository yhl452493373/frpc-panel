(function ($) {
    //param names in Basic tab
    var basicParams = [
        {
            name: 'name',
            defaultValue: '-'
        }, {
            name: 'type',
            defaultValue: '-'
        }, {
            name: 'localIP',
            defaultValue: '-'
        }, {
            name: 'localPort',
            defaultValue: '-'
        }, {
            name: 'customDomains',
            defaultValue: '-'
        }, {
            name: 'subdomain',
            defaultValue: '-'
        }, {
            name: 'remotePort',
            defaultValue: '-'
        }, {
            name: 'transport.useEncryption',
            defaultValue: false,
        }, {
            name: 'transport.useCompression',
            defaultValue: false,
        }
    ];
    var mapParams = [{
        name: 'plugin.ClientPluginOptions',
        map: 'plugin'
    }];
    var paramTypes = {
        number: [
            'healthCheck.timeoutSeconds',
            'healthCheck.maxFailed',
            'healthCheck.intervalSeconds',
            'localPort',
        ],
        boolean: [
            'transport.useEncryption',
            'transport.useCompression'
        ],
        array: [
            'customDomains',
            'locations',
            'allowUsers'
        ],
        map: [
            'metadatas',
            'requestHeaders.set',
            'plugin.ClientPluginOptions.requestHeaders.set'
        ]
    }

    /**
     *  a.b.c = 1
     *  a.b.d = [2,3]
     *  to
     *  {a: {
     *          b: {
     *              c: 1
     *              d: "[2,3]"
     *          }
     *      }
     *  }
     *
     * @param obj json object
     * @param keys all keys split from key string like 'a.b.c'
     * @param value default value
     * @param stringifyArray if true, when value is an array, it will stringify by JSON.stringify(value)
     */
    function expandJSONKeys(obj, keys, value, stringifyArray) {
        stringifyArray = stringifyArray == null ? true : stringifyArray;
        var currentIndex = this.index || 0;
        var childrenIndex = (currentIndex + 1) > keys.length ? null : (currentIndex + 1);
        var currentKey = keys[currentIndex], currentValue = {};
        var childrenKey = childrenIndex == null ? null : keys[childrenIndex];
        if (obj.hasOwnProperty(currentKey)) {
            currentValue = obj[currentKey];
        } else {
            obj[currentKey] = currentValue;
        }

        if (childrenKey != null) {
            this.index = childrenIndex;
            expandJSONKeys(currentValue, keys, value, stringifyArray);
        } else {
            if (value != null) {
                if (Array.isArray(value) && stringifyArray) {
                    obj[currentKey] = JSON.stringify(value);
                } else {
                    obj[currentKey] = value;
                }
            }
            this.index = 0;
        }
    }

    function expandJSON(obj) {
        var newObj = {};
        var inPopup = $('#proxyForm').length !== 0;
        for (var name in obj) {
            var value = obj[name];
            if (value === '') {
                continue;
            }
            var nameI18n = name;
            if (inPopup) {
                nameI18n = $('#proxyForm [name="' + name + '"]').closest('.layui-form-item').children('.layui-form-label').text();
                if (nameI18n === '') {
                    nameI18n = name;
                }
            }

            if (paramTypes.number.indexOf(name) !== -1) {
                value = parseInt(value);
                if (isNaN(value)) {
                    throw new Error(nameI18n + i18n['RequireNumber']);
                }
            } else if (paramTypes.boolean.indexOf(name) !== -1) {
                if (typeof value === "string" && (value === 'true' || value === 'false')) {
                    value = value === 'true';
                } else if (typeof value !== 'boolean') {
                    throw new Error(nameI18n + i18n['RequireBoolean']);
                }
            } else if (paramTypes.array.indexOf(name) !== -1) {
                try {
                    if (/^\s*\[.*]\s*$/.test(value)) {
                        value = eval('(' + value + ')') || [];
                    } else {
                        throw new Error('value format incorrect');
                    }
                } catch (e) {
                    throw new Error(nameI18n + i18n['RequireArray']);
                }
            } else {
                for (var i = 0; i < paramTypes.map.length; i++) {
                    var key = paramTypes.map[i];
                    if (name.startsWith(key)) {
                        var json = {};
                        json[name.substring(key.length + 1, name.length)] = value;
                        value = json;
                        name = name.substring(0, key.length)
                        break;
                    }
                }
            }

            expandJSONKeys(newObj, name.split("."), value, false);
        }
        return newObj;
    }

    /**
     * {a: {
     *          b: {
     *              c: 1
     *              d: [2,3]
     *          }
     *      }
     *  }
     *  to
     *  {
     *      'a.b.c': 1,
     *      'a.b.d': '[2,3]'
     *  }
     *
     * @param obj json object
     * @returns {*} flatted json key array
     */
    function flatJSON(obj) {
        var flat = function (obj, prentKey, flattedJSON) {
            flattedJSON = flattedJSON || {};
            prentKey = prentKey || '';
            if (prentKey !== '')
                prentKey = prentKey + '.';
            for (var key in obj) {
                var value = obj[key];
                if (typeof value === 'object' && Object.prototype.toString.call(value) === '[object Object]') {
                    flat(value, prentKey + key, flattedJSON);
                } else {
                    for (var mapParam of mapParams) {
                        if (prentKey.startsWith(mapParam.name)) {
                            prentKey = mapParam.map + '.';
                            break;
                        }
                    }
                    if (Array.isArray(value)) {
                        flattedJSON[prentKey + key] = JSON.stringify(value);
                    } else {
                        flattedJSON[prentKey + key] = value;
                    }
                }
            }
            return flattedJSON;
        }
        return flat(obj);
    }

    window.basicParams = basicParams;
    window.expandJSONKeys = expandJSONKeys;
    window.expandJSON = expandJSON;
    window.flatJSON = flatJSON;
})(layui.$);