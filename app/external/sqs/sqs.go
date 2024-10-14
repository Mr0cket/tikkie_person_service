package sqs

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Sqs struct {
	client   sqs.Client
	QueueUrl *string
	ctx      context.Context
}

func NewClient(ctx context.Context, cfg aws.Config, queueName, region string) *Sqs {

	// 3. Publish a new event to the SQS queue

	svc := sqs.NewFromConfig(cfg)
	result, err := svc.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		log.Fatalf("Unable to fetch queue URL: %v", err)
	}

	return &Sqs{
		QueueUrl: result.QueueUrl,
		client:   *svc,
	}
}

func (s *Sqs) SendMessage(attributes map[string]string, message string) string {
	messageOutput, err := s.client.SendMessage(s.ctx, &sqs.SendMessageInput{
		DelaySeconds:      *aws.Int32(10),
		MessageAttributes: s.createAttributes(attributes),
		MessageBody:       aws.String(message),
		QueueUrl:          s.QueueUrl,
	})
	if err != nil {
		log.Fatalf("failed to send message, %v", err)
	}
	return *messageOutput.MessageId
}

// Only support string attributes for now
func (s *Sqs) createAttributes(attrs map[string]string) map[string]types.MessageAttributeValue {
	attributeMap := make(map[string]types.MessageAttributeValue, len(attrs))

	for k, v := range attrs {
		attributeMap[k] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(v),
		}

	}
	return attributeMap
}
