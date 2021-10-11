package main

import (
	"errors"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type NetConf struct {
	prefix            net.IPNet
	preferredLifetime time.Duration
	validLifetime     time.Duration
	dnsServers        []net.IP
	dnsSearchList     []string
}

func (conf *NetConf) logFields() log.Fields {
	return log.Fields{
		"prefix":  conf.prefix.String(),
		"preflt":  conf.preferredLifetime,
		"validlt": conf.validLifetime,
		"dns":     conf.dnsServers,
		"search":  conf.dnsSearchList,
	}
}

func staticNetConf(prefix string) (NetConf, error) {
	_, net, err := net.ParseCIDR(prefix)
	if err != nil {
		return NetConf{}, err
	}
	len, total := net.Mask.Size()
	if total != 128 {
		return NetConf{}, errors.New("Please specify an IPv6 prefix")
	}
	if len < 48 || len > 96 {
		return NetConf{}, errors.New("Please specify an IPv6 prefix between /48 and /96 in length")
	}
	log.WithFields(log.Fields{"prefix": net}).
		Info("Using static IPv6 prefix")
	return NetConf{prefix: *net}, nil
}
