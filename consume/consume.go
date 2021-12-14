package consume

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/chmoon93/dabbi-sqs-consumer/log"
)

type SQSGetQueueUrlAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)
}

type SQSGetLPMsgAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

func GetQueueURL(c context.Context, api SQSGetQueueUrlAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

// GetLPMessages gets the messages from an Amazon SQS long polling queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a ReceiveMessageOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ReceiveMessage.
func GetLPMessages(c context.Context, api SQSGetLPMsgAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}

// SQSDeleteMessageAPI defines the interface for the GetQueueUrl and DeleteMessage functions.
// We use this interface to test the functions using a mocked service.
type SQSDeleteMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	DeleteMessage(ctx context.Context,
		params *sqs.DeleteMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

// RemoveMessage deletes a message from an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a DeleteMessageOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to DeleteMessage.
func RemoveMessage(c context.Context, api SQSDeleteMessageAPI, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return api.DeleteMessage(c, input)
}

func ConsumeMessages(ctx context.Context, cfg aws.Config) {
	// for {
	// 	log.Info("consume messages")
	// }

	client := sqs.NewFromConfig(cfg)
	input := &sqs.GetQueueUrlInput{
		QueueName: aws.String("dabbi-sqs"),
	}

	result, err := GetQueueURL(context.TODO(), client, input)
	if err != nil {
		log.Debug("Got an error getting the queue URL:")
		log.Debug(err)
		return
	}

	log.Info("# QueueURL: ", *result.QueueUrl)
	waitTime := 5
	queueURL := result.QueueUrl

	for {
		log.Info("111")
		mInput := &sqs.ReceiveMessageInput{
			QueueUrl: queueURL,
			AttributeNames: []types.QueueAttributeName{
				"SentTimestamp",
			},
			MaxNumberOfMessages: 1,
			MessageAttributeNames: []string{
				"All",
			},
			WaitTimeSeconds: int32(waitTime),
		}

		resp, err := GetLPMessages(context.TODO(), client, mInput)
		if err != nil {
			log.Debug("Got an error receiving messages:")
			log.Debug(err)
			return
		}

		// do process
		for _, msg := range resp.Messages {
			log.Debugf("Message [%s] [%s] [%s]", *msg.MessageId, *msg.Body, *msg.MD5OfBody)
		}

		// do delete
		for _, dMsg := range resp.Messages {
			dMInput := &sqs.DeleteMessageInput{
				QueueUrl:      queueURL,
				ReceiptHandle: dMsg.ReceiptHandle,
			}

			_, err = RemoveMessage(context.TODO(), client, dMInput)
			if err != nil {
				log.Debug("Got an error deleting the message: ")
				log.Debug(err)
				return
			}
		}

	}
}
