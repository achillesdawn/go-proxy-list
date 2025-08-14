package geonode

import (
	"fmt"
	"strings"
)

func (p Proxy) GetIP() string {
	return p.IP
}

func (p Proxy) GetAddress() string {
	return fmt.Sprintf("%s:%s", p.IP, p.Port)
}
func (p Proxy) GetProtocol() string {
	return strings.Join(p.Protocols, ",")
}
