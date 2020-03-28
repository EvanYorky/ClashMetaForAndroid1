package config

import (
	"fmt"
	"github.com/Dreamacro/clash/component/fakeip"
	"github.com/Dreamacro/clash/config"
	"net/url"
)

var (
	OptionalDnsPatch  *config.RawDNS
	DnsPatch          *config.RawDNS
	NameServersAppend []string

	cachedPool *fakeip.Pool
)

func patchRawConfig(rawConfig *config.RawConfig) {
	rawConfig.DNS.FakeIPRange = "198.18.0.0/16"
	rawConfig.Experimental.Interface = ""
	rawConfig.ExternalUI = ""
	rawConfig.ExternalController = ""

	if len(rawConfig.Rule) != 0 {
		rawConfig.Rule = append([]string{fmt.Sprintf("IP-CIDR,%s,REJECT,no-resolve", tunAddress)}, rawConfig.Rule...)
	} else {
		rawConfig.RuleOld = append([]string{fmt.Sprintf("IP-CIDR,%s,REJECT,no-resolve", tunAddress)}, rawConfig.RuleOld...)
	}

	if d := DnsPatch; d != nil {
		rawConfig.DNS = *d
	} else if d := OptionalDnsPatch; d != nil {
		if !rawConfig.DNS.Enable {
			rawConfig.DNS = *d
		}
	}

	if nameServersAppend := NameServersAppend; len(nameServersAppend) > 0 {
		d := &rawConfig.DNS
		nameServers := make([]string, len(nameServersAppend)+len(d.NameServer))
		copy(nameServers, nameServersAppend)
		copy(nameServers[len(nameServersAppend):], d.NameServer)

		d.NameServer = nameServers
	}

	providers := rawConfig.ProxyProvider

	if len(rawConfig.ProxyProvider) == 0 {
		providers = rawConfig.ProxyProviderOld
	}

	for _, provider := range providers {
		path, ok := provider["path"].(string)
		if !ok {
			continue
		}

		provider["path"] = url.QueryEscape(path)
	}
}

func patchConfig(config *config.Config) {
	if config.DNS.FakeIPRange != nil {
		if c := cachedPool; c != nil {
			if config.DNS.FakeIPRange.Gateway().String() == c.Gateway().String() {
				config.DNS.FakeIPRange = c
			}
		} else {
			cachedPool = config.DNS.FakeIPRange
		}
	}
}
