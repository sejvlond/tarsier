package net

import (
	"fmt"

	tarsier_plugins "github.com/sejvlond/tarsier/plugins"
)

func NewDial(lgr tarsier_plugins.Logger, container *ConnContainer) *Dial {
	return &Dial{
		lgr:       lgr,
		container: container,
	}
}

type Dial struct {
	lgr       tarsier_plugins.Logger
	container *ConnContainer
}

type DialData struct {
	Network string
	Address string
	Amount  uint
}

func (cmd *Dial) Name() string {
	return NAME + "/dial"
}

func (cmd *Dial) Description() string {
	return `Dial will open as many connections as you say<br>
Data:
<pre>
	network: Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only), "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and "unixpacket".
	address: For TCP and UDP networks, addresses have the form host:port. If host is a literal IPv6 address it must be enclosed in square brackets as in "[::1]:80" or "[ipv6-host%zone]:80". The functions JoinHostPort and SplitHostPort manipulate addresses in this form. If the host is empty, as in ":80", the local system is assumed. 
	amount: how many connections to open
</pre>`
}

func (cmd *Dial) DataStruct() interface{} {
	return &DialData{}
}

func (cmd *Dial) Execute(raw interface{}) (string, error) {
	data, ok := raw.(*DialData)
	if !ok {
		return "", fmt.Errorf("Tarsier, those data are not mine!")
	}
	var (
		cnt uint
		err error
	)
	cmd.lgr.Infof("Opening %v connections", data.Amount)
	for cnt = 0; cnt < data.Amount; cnt++ {
		if err = cmd.container.Dial(data.Network, data.Address); err != nil {
			break
		}
	}
	if cnt < data.Amount {
		return fmt.Sprintf("Could not open more than %v connections: %v",
			cnt, err), err
	}
	return fmt.Sprintf("%v connections was opened", cnt), nil
}
