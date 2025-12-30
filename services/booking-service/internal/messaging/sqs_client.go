package messaging

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSClient struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSClient(ctx context.Context, region, queueURL string) (*SQSClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return &SQSClient{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}

func (c *SQSClient) Send(ctx context.Context, body string, messageGroupID string, messageDeduplicationID string) (string, error) {
	in := &sqs.SendMessageInput{
		QueueUrl:    aws.String(c.queueURL),
		MessageBody: aws.String(body),
	}

	if strings.HasSuffix(c.queueURL, ".fifo") {
		if messageGroupID != "" {
			in.MessageGroupId = aws.String(messageGroupID)
		}
		if messageDeduplicationID != "" {
			in.MessageDeduplicationId = aws.String(messageDeduplicationID)
		}
	}

	out, err := c.client.SendMessage(ctx, in)
	if err != nil {
		return "", err
	}
	return aws.ToString(out.MessageId), nil
}
