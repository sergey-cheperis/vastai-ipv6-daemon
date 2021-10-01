package main

import (
	"context"
	"errors"
	"log"
	"time"
)

type DhcpKeeper struct {
	ifname  string
	netConf NetConf
	ctx     context.Context
}

func startDhcp(ctx context.Context, ifname string) (NetConf, error) {
	keeper := DhcpKeeper{
		ifname: ifname,
		ctx:    ctx,
	}
	err := keeper.renew()
	if err != nil {
		return NetConf{}, err
	}
	go keeper.renewLoop()
	return keeper.netConf, nil
}

func (c *DhcpKeeper) renew() error {
	netConf, err := receiveConfWithDhcp(c.ctx, c.ifname)
	if err != nil {
		return err
	}
	log.Printf("Received network configuration: %s", netConf.String())

	len, _ := netConf.prefix.Mask.Size()
	if len < 48 || len > 96 {
		return errors.New("Delegated prefix must be between /48 and /96 in length")
	}

	c.netConf = netConf
	return nil
}

func (c *DhcpKeeper) renewLoop() {
	// TODO what to do if the prefix changes or expires?
	for {
		time.Sleep(c.netConf.preferredLifetime)
		for {
			err := c.renew()
			if err == nil {
				break
			}
			delay := 15 * time.Minute
			log.Printf("Error: %v, will retry in %s", err, delay.String())
			time.Sleep(delay)
		}
	}
}
