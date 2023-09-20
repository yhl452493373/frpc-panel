package controller

import (
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

func (c *HandleController) reloadFrpc() error {
	return nil
}
