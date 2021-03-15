package main

import "github.com/braianpablodiaz/meli-proxy/proxy"

func main() {
	proxy := proxy.NewProxy()
	proxy.StartProxy()
}