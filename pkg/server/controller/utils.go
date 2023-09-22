package controller

import (
	"github.com/fatedier/frp/pkg/config"
	"github.com/vaughan0/go-ini"
	"strings"
)

func trimString(str string) string {
	return strings.TrimSpace(str)
}

func cleanString(originalString string) string {
	return trimString(originalString)
}

func stringContains(element string, data []string) bool {
	for _, v := range data {
		if element == v {
			return true
		}
	}
	return false
}

func (c *HandleController) parseConfigure(content, proxyType string) (interface{}, error) {
	currentProxies := make(map[string]ini.Section)
	clientProxies = make(map[string]ini.Section)
	common, err := config.UnmarshalClientConfFromIni(content)
	if err != nil {
		return nil, err
	}
	cfg, err := ini.Load(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	for name, section := range cfg {
		if name == "common" {
			clientCommon = section
			continue
		}
		if strings.ToLower(section["type"]) == strings.ToLower(proxyType) {
			currentProxies[name] = section
		}
		clientProxies[name] = section
		delete(clientProxies[name], "name")
	}

	if proxyType == "none" {
		return common, nil
	}

	return currentProxies, nil
}
