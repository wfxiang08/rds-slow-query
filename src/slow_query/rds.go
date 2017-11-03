package slow_query

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	log "github.com/wfxiang08/cyutils/utils/rolling_log"
	"strings"
)

func GetSession(conf *DatabaseConfig) *session.Session {
	awsCreditial := credentials.NewStaticCredentials(conf.AwsKey, conf.AwsSecret, "")
	config := aws.NewConfig().WithRegion(conf.AwsRegion).WithCredentials(awsCreditial)
	r, _ := session.NewSession(config)
	return r
}

func DownloadHourlyLogFile(db, logfile string, conf *DatabaseConfig) (string, error) {
	client := rds.New(GetSession(conf))

	var results []string

	input := &rds.DownloadDBLogFilePortionInput{
		DBInstanceIdentifier: aws.String(db),
		Marker:               aws.String("0"),
		LogFileName:          aws.String(logfile),
		NumberOfLines:        aws.Int64(10000), // 每次下载一万条
	}

	// 遍历所有的Pages, 下载完整的log文件
	err := client.DownloadDBLogFilePortionPages(input,
		func(page *rds.DownloadDBLogFilePortionOutput, lastPage bool) bool {
			results = append(results, *page.LogFileData)
			return true
		})

	if err != nil {
		log.ErrorErrorf(err, "DownloadDBLogFilePortionPages failed")
		return "", err
	} else {
		log.Printf("Download %s %s succeed", db, logfile)
		return strings.Join(results, "\n"), nil
	}
}
