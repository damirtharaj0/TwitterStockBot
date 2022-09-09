package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/get-stock-tweets", GetStockTweets)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// homepage
func homePage(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

// recieves stock ticker in query parameter
func GetStockTweets(w http.ResponseWriter, r *http.Request) {
	stockticker := r.FormValue("ticker")
	log.Println(stockticker)
	if stockticker == "" {
		fmt.Fprintf(w, "Enter a stock as a query\n")
		fmt.Fprintf(w, "Ex. localhost:8080/get-stock-tweets?ticker=aapl")
		return
	}

	respStr := ScrapeTwitter(stockticker)

	type TwitterResp struct {
		Data []struct {
			Text string `json:"text"`
		} `json:"data"`
	}

	resp := TwitterResp{}
	json.Unmarshal([]byte(respStr), &resp)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// scrape twitter for
func ScrapeTwitter(stockName string) string {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, `https://api.twitter.com/2/tweets/search/recent`, nil)
	checkErr(err)

	q := req.URL.Query()
	q.Add("query", "#"+stockName+" stock -wallstreetbet lang:en -is:retweet")
	req.URL.RawQuery = q.Encode()
	apiToken := ""
	req.Header.Set("Authorization", "Bearer " + apiToken)
	resp, err := client.Do(req)
	checkErr(err)

	b, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	return string(b)
}

// checks error
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
