package slow_query

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	log "github.com/wfxiang08/cyutils/utils/rolling_log"
	"strings"
)

const (
	kDefaultCharacterSet = "utf-8"
)

func SendEmail(title, alarmContent string, sender string, receivers []string, conf *DatabaseConfig) {

	sess := session.Must(session.NewSession())

	svc := ses.New(sess, &aws.Config{Region: aws.String(conf.AwsRegion)})

	var receiverList []*string
	for _, receiver := range receivers {
		receiverList = append(receiverList, aws.String(receiver))
	}
	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: receiverList,
		},
		Message: &ses.Message{ // Required
			Body: &ses.Body{ // Required
				Html: &ses.Content{
					Data:    aws.String(alarmContent), // Required
					Charset: aws.String(kDefaultCharacterSet),
				},
				Text: &ses.Content{
					Data:    aws.String(alarmContent), // Required
					Charset: aws.String(kDefaultCharacterSet),
				},
			},
			Subject: &ses.Content{ // Required
				Data:    aws.String(title), // Required
				Charset: aws.String(kDefaultCharacterSet),
			},
		},
		Source: aws.String(sender),
		ReplyToAddresses: []*string{
			aws.String(sender),
		},
		ReturnPath: aws.String(sender),
		Tags: []*ses.MessageTag{
			{ // Required
				Name:  aws.String("MessageTagName"),  // Required
				Value: aws.String("MessageTagValue"), // Required
			},
		},
	}
	_, err := svc.SendEmail(params)

	if err != nil {
		log.ErrorErrorf(err, "SendEmail failed")
		return
	} else {
		log.Printf("Email Send Succeed")
	}
}

func FormatMail(lines []*Summary, maxLine int, verbose bool) string {
	var results []string
	for _, line := range lines {
		if line.Count > 10 || float64(line.Total)/float64(line.Count) > 0.22 {
			results = append(results, line.HTML())
			if len(results) > maxLine {
				break
			}
			if verbose {
				fmt.Printf("%s\n", line.String())
			}
		}
	}
	return strings.Join(results, "<br/>")
}
