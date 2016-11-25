package main

import "net/http"

func checkSite(site string) (int, error) {
	res, err := http.Get(site)
	if err != nil {
		return 0, err
	}

	return res.StatusCode, nil
}
