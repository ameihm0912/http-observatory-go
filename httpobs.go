// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com

// httpobsgo is a library to provide a simple interface to run scans using
// the Mozilla HTTP Observatory, and collect the results.
package httpobsgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// The API URL to use for requests
var APIUrl = "https://http-observatory.security.mozilla.org/api/v1/"

// Maximum amount of time we will wait for a scan to complete
var MaxWait = time.Second * 300

// Time between polls
var PollInterval = time.Second * 5

type HTTPObsTime struct{ time.Time }

// We require custom unmarshalling to handle the data formats returned by the
// HTTP observatory API
func (t *HTTPObsTime) UnmarshalJSON(b []byte) (err error) {
	var zeroTime time.Time
	if string(b) == "null" {
		t.Time = zeroTime
		return
	}
	t.Time, err = time.Parse("\"Mon, 2 Jan 2006 15:04:05 MST\"", string(b))
	return
}

// Result object returned by RunScan
type ScanObject struct {
	StartTime       HTTPObsTime       `json:"start_time"`
	EndTime         HTTPObsTime       `json:"end_time"`
	State           string            `json:"state"`
	Grade           string            `json:"grade"`
	Score           int               `json:"score"`
	ScanID          float64           `json:"scan_id"`
	TestsFailed     int               `json:"tests_failed"`
	TestsPassed     int               `json:"tests_passed"`
	TestsQuantity   int               `json:"tests_quantity"`
	ResponseHeaders map[string]string `json:"response_headers"`
	Error           string            `json:"error"`
}

// Run a new scan using the HTTP Observatory. Hostname is the host which you want
// to target.
//
// If hidden is true, the request will note the scan results should be
// hidden from public results returned by getRecentScans.
//
// If rescan is true, this informs the API to conduct a new scan of the host and
// not return any cached results from a recent scan of the same host.
func RunScan(hostname string, hidden bool, rescan bool) (ScanObject, error) {
	var results ScanObject
	u := APIUrl + fmt.Sprintf("analyze?host=%v", url.QueryEscape(hostname))

	form := url.Values{}
	form.Add("hidden", fmt.Sprintf("%v", hidden))
	form.Add("rescan", fmt.Sprintf("%v", rescan))

	start := time.Now()
	for {
		res, err := http.PostForm(u, form)
		if err != nil {
			return results, err
		}
		err = json.NewDecoder(res.Body).Decode(&results)
		if err != nil {
			return results, err
		}
		res.Body.Close()
		if results.Error != "" {
			err = fmt.Errorf("httpobsgo: %v", results.Error)
			return results, err
		}
		if results.State == "FINISHED" {
			break
		}
		if time.Now().Sub(start) > MaxWait {
			err = fmt.Errorf("httpobsgo: maximum scan duration exceeded")
			return results, err
		}
		time.Sleep(PollInterval)
	}
	return results, nil
}
