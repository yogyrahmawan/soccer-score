package app

import (
	"golang.org/x/net/html"
	"net/http"
)

//GetParseableHTML get html root node from url
func GetParseableHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return root, nil
}
