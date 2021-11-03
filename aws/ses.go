package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	etc2 "github.com/msyhu/naekaracubae-scraping/etc"
	_struct2 "github.com/msyhu/naekaracubae-scraping/struct"
	"time"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "noreply@msyhu.com"

	// The character encoding for the email.
	CharSet = "UTF-8"
)

func SendMail(contents *string, subscribers []_struct2.Subscriber) string {
	// The subject line for the email.
	var today = time.Now().Format("2006-01-02")
	Subject := "[네,카라쿠배] " + today + " 개발자 채용 일보가 도착했습니다!👩‍💻"

	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2")},
	)
	etc2.CheckErr(err)

	// Create an SES session.
	svc := ses.New(sess)
	var result string
	var toAddresses []*string
	for i := range subscribers {
		toAddresses = nil
		toAddresses = append(toAddresses, &subscribers[i].Email)

		// Assemble the email.
		input := &ses.SendEmailInput{
			Destination: &ses.Destination{
				CcAddresses: []*string{},
				ToAddresses: toAddresses,
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Html: &ses.Content{
						Charset: aws.String(CharSet),
						Data:    aws.String(*contents),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(Subject),
				},
			},
			Source: aws.String(Sender),
			// Uncomment to use a configuration set
			//ConfigurationSetName: aws.String(ConfigurationSet),
		}

		// Attempt to send the email.
		output, err := svc.SendEmail(input)
		result += output.String() + "\n"

		// Display error messages if they occur.
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case ses.ErrCodeMessageRejected:
					fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
				case ses.ErrCodeMailFromDomainNotVerifiedException:
					fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
				case ses.ErrCodeConfigurationSetDoesNotExistException:
					fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}

			return ""
		}
	}

	return result

}
