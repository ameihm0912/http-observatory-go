// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor:
// - Aaron Meihm ameihm@mozilla.com

package httpobsgo

import (
	"testing"
)

func TestRunScan(t *testing.T) {
	res, err := RunScan("www.mozilla.org", false, false)
	if err != nil {
		t.Fatal(err)
	}
	if res.Grade == "" {
		t.Fatal("successful scan result returned no grade")
	}
}

func TestRunScanBadHost(t *testing.T) {
	_, err := RunScan("127.0.0.1", false, false)
	if err == nil {
		t.Fatal("scan of bad host did not return error")
	}
}
