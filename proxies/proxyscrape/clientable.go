package proxyscrape

func (p Proxy) GetIP() string {
	return p.Ip
}

func (p Proxy) GetAddress() string {
	return p.Proxy
}
func (p Proxy) GetProtocol() string {
	return p.Protocol
}
