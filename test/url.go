package main

import (
	"fmt"
	"net/url"
)

func main() {
	str, err := url.QueryUnescape("http://localhost:14000/app")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
	code := "0000"
	aurl := fmt.Sprintf("/token?grant_type=authorization_code&client_id=1234&client_secret=aabbccdd&state=xyz&redirect_uri=%s&code=%s",
		url.QueryEscape("http://localhost:14000/appauth/code"), url.QueryEscape(code))
	fmt.Println(aurl)
}
