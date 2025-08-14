package proxifly

func (p Proxy) GetIP() string {
	return p.IP
}

func (p Proxy) GetAddress() string {
	return p.Proxy
}
func (p Proxy) GetProtocol() string {
	return p.Protocol
}
