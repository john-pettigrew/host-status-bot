package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SiteResult represents a site from the database
type SiteResult struct {
	id  int
	url string
}

func startLoop() error {
	var sites []SiteResult
	downCh := make(chan string, 5)

	go watchReport(downCh)

	//connect to db
	db, err := sql.Open("sqlite3", "./status-bot.db")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS sites (id INTEGER PRIMARY KEY, url TEXT REQUIRED)")
	if err != nil {
		return err
	}

	//start loop
	for {
		sites, err = getSitesFromDB(db)
		checkSites(sites, downCh)
		time.Sleep(time.Second * 3)
	}
}

func checkSites(sites []SiteResult, downCh chan string) {
	var status int
	var err error
	fmt.Println("loop")
	fmt.Println(sites)

	// Check each site
	for i := 0; i < len(sites); i++ {
		status, err = checkSite(sites[i].url)
		if err != nil {
			downCh <- "Unable to check " + sites[i].url + "."
		}
		if status < 200 || status >= 300 {
			fmt.Println(status)
			downCh <- "site \"" + sites[i].url + "\" is returning a \"" + strconv.Itoa(status) + "\" code."
		}
	}
}

func getSitesFromDB(db *sql.DB) ([]SiteResult, error) {
	q := "SELECT id, url FROM sites"
	rows, err := db.Query(q)
	if err != nil {
		return []SiteResult{}, err
	}
	defer rows.Close()
	var results []SiteResult
	var currentResult SiteResult
	for rows.Next() {
		currentResult = SiteResult{}
		err = rows.Scan(&currentResult.id, &currentResult.url)
		if err != nil {
			return []SiteResult{}, err
		}
		results = append(results, currentResult)
	}

	return results, nil
}

func checkSite(site string) (int, error) {
	res, err := http.Get(site)
	if err != nil {
		return 0, err
	}
	return res.StatusCode, nil
}

func watchReport(downCh chan string) {
	for errM := range downCh {
		fmt.Println("Sending to IRC")
		fmt.Println(errM)
	}
}
