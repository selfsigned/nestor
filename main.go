package main

import (
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
)

var (
	api_url = "https://api-nestor.com/"
	menu_route = "menu/"
)

type Menu struct {
}

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	url := api_url + menu_route + os.Args[1]
	resp, err := http.Get(url)
	if err != nil {
		panic("Something went wrong")
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "%s Error:%d\n",
			url, resp.StatusCode)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(resp.Body)
	println(string(body))
	defer resp.Body.Close()
}
