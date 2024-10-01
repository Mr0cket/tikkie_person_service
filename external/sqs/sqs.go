package sqs

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Sqs struct {
	QueueName string
	client    sqs.Client
}

func NewClient(queueName string) Sqs {

	// 3. Publish a new event to the SQS queue
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := sqs.NewFromConfig(cfg)

	return Sqs{
		QueueName: queueName,
		client:    *svc,
	}
}

func (s *Sqs) SendMessage(attributes map[string]string, message interface{}) {
	result, err := s.client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(s.QueueName),
	})
	if err != nil {
		log.Fatalf("Unable to fetch queue URL: %v", err)
	}

	_, err = s.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		DelaySeconds:      *aws.Int32(10),
		MessageAttributes: s.createAttributes(attributes),
		MessageBody:       aws.String("New person created"), // TODO: use a proper JSON object
		QueueUrl:          result.QueueUrl,
	})

	if err != nil {
		log.Fatalf("failed to send message, %v", err)
	}
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
