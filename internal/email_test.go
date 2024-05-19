package internal

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
)

var _ = Describe("SES", func() {
	Context("sending an email", func() {
		var (
			svc *ses.Client
		)

		BeforeSuite(func() {
			// Start the local server
			go func() {
				StartServer()
			}()

			server := smtpmock.New(smtpmock.ConfigurationAttr{
			LogToStdout:       true,
			LogServerActivity: true,
			})
		
			// To start server use Start() method
			if err := server.Start(); err != nil {
				fmt.Println(err)
			}
			Config.SmtpHost="127.0.0.1"
			Config.SmtpPort=server.PortNumber()

			// Set up a new SES session
			sesConfig, err := constructAWSConfiguration("http://localhost:8081")
			Expect(err).NotTo(HaveOccurred())

			// Create a new SES client
			svc = ses.NewFromConfig(sesConfig)
		})

		It("should send an email successfully", func() {
			// Set up the email message
			input := &ses.SendEmailInput{
				Destination: &types.Destination{
					ToAddresses: []string{
						"recipient@example.com",
					},
				},
				Message: &types.Message{
					Body: &types.Body{
						Html: &types.Content{
							Data: aws.String("<p>Hello, world!</p>"),
						},
					},
					Subject: &types.Content{
						Data: aws.String("Test email"),
					},
				},
				Source: aws.String("sender@example.com"),
			}

			// Send the email
			result, err := svc.SendEmail(context.TODO(), input)
			Expect(err).NotTo(HaveOccurred())

			// Check the response
			Expect(result).NotTo(BeNil())
			Expect(result.MessageId).NotTo(BeNil())
			Expect(*result.MessageId).NotTo(BeEmpty())

			fmt.Println("Email sent successfully!")
		})

		It("should fail to send an email if the recipient address is invalid", func() {
			// Set up the email message with an invalid recipient address
			input := &ses.SendEmailInput{
				Destination: &types.Destination{
					ToAddresses: []string{
						"invalid_email.com",
					},
				},
				Message: &types.Message{
					Body: &types.Body{
						Html: &types.Content{
							Data: aws.String("<p>Hello, world!</p>"),
						},
					},
					Subject: &types.Content{
						Data: aws.String("Test email"),
					},
				},
				Source: aws.String("sender@example.com"),
			}

			// Send the email
			_, err := svc.SendEmail(context.TODO(), input)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("400"))
		})
	})
})

func TestEmailSendingCompliance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

func constructAWSConfiguration(endpoint string) (cfg aws.Config, err error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if endpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           endpoint,
				SigningRegion: "eu-central-1",
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	// Create an AWS configuration
	configArray := []func(options *config.LoadOptions) error{
		config.WithRegion("eu-central-1"),
		config.WithDefaultRegion("eu-central-1"),
		config.WithEndpointResolverWithOptions(customResolver),
	}

	// For local development we need to use anonymous credentials
	configArray = append(configArray, config.WithCredentialsProvider(aws.AnonymousCredentials{}))

	return config.LoadDefaultConfig(
		context.Background(),
		configArray...,
	)
}
