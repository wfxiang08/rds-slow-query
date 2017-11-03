package slow_query

import (
	"io/ioutil"
	"os"
	"testing"
)

// go test slow_query -v -run "TestDownloadHourlyLogFile$"
func TestDownloadHourlyLogFile(t *testing.T) {

	content, err := DownloadHourlyLogFile("starmaker-sharding-00-readonly", "slowquery/mysql-slowquery.log.11")
	if err != nil {
		return
	}

	ioutil.WriteFile("001.txt", []byte(content), os.ModePerm)

}
