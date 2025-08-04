simplified API to get free proxies from ProxyScrape and GeoNodes

```go
// get working proxies
proxies, err := proxyscrape.WorkingProxies()

for _, proxy := range proxies {

	// client is of type *http.Client, already configured with proxy information
	// it can be used simply to execute requests
	// client.Do(request)
	client := p.CreateClient()
}
	

```
