package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetStargazers(t *testing.T) {
	filePath := "testdata/simple_stars.json"
	expectedTimestamps := []Stargazer{{"2015-10-31T10:00:00Z"}, {"2015-10-31T11:00:00Z"}, {"2015-10-31T12:00:00Z"}}
	body, err := ioutil.ReadFile(filePath)
	batch := 1
	if err != nil {
		t.Fatalf("An error occured while reading the file %s: %v\n", filePath, err)
	}

	serverUrl := ""
	handler := func(w http.ResponseWriter, r *http.Request) {
		linkFormat := BuildLinksFormat(serverUrl)
		w.Header().Add("Link", fmt.Sprintf(linkFormat, 1, 1))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	serverUrl = server.URL + "?per_page=" + strconv.Itoa(batch)

	timestamps, _, err := GetStargazers(server.URL, "token")
	if err != nil {
		t.Fatalf("An error occured while requesting the timestamps: %v\n", err)
	}
	if len(timestamps) != len(expectedTimestamps) {
		t.Fatalf("The expected timestamps %v and the actual ones %v"+
			" don't have the same size\n", expectedTimestamps, timestamps)
	}
	equals := true
	for i, v := range timestamps {
		if expectedTimestamps[i] != v {
			equals = false
			break
		}
	}
	if !equals {
		t.Fatalf("The expected timestamps %v and the actual ones %v"+
			" don't have the same values\n", expectedTimestamps, timestamps)
	}
}
