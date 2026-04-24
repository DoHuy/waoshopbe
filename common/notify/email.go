package notify

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SESClient struct {
	client *ses.Client
	sender string
}

func NewSESClient(ctx context.Context, region, accessKey, secretKey, senderEmail string) (*SESClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("Error loading AWS config: %v", err)
	}

	return &SESClient{
		client: ses.NewFromConfig(cfg),
		sender: senderEmail,
	}, nil
}

func (s *SESClient) SendEmail(ctx context.Context, toEmail, subject, htmlBody string) error {
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{toEmail},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(s.sender),
	}

	_, err := s.client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("Error sending email via AWS SES: %v", err)
	}
	return nil
}
