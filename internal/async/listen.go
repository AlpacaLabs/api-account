package async

import (
	"context"
	"fmt"

	"github.com/AlpacaLabs/api-account/internal/configuration"
	"github.com/AlpacaLabs/api-account/internal/service"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"

	goKafka "github.com/AlpacaLabs/go-kafka"
	log "github.com/sirupsen/logrus"
)

const (
	TopicForConfirmEmailAddressRequest = "confirm-email-address-request"
	TopicForConfirmPhoneNumberRequest  = "confirm-phone-number-request"
)

func HandleConfirmEmailAddressRequest(config configuration.Config, s service.Service) {
	handle(TopicForConfirmEmailAddressRequest, config, handleConfirmEmailAddressRequest(s))
}

func HandleConfirmPhoneNumberRequest(config configuration.Config, s service.Service) {
	handle(TopicForConfirmPhoneNumberRequest, config, handleConfirmPhoneNumberRequest(s))
}

func handle(topic string, config configuration.Config, fn goKafka.ProcessFunc) {
	ctx := context.TODO()

	groupID := config.AppName
	brokers := []string{
		fmt.Sprintf("%s:%d", config.KafkaConfig.Host, config.KafkaConfig.Port),
	}

	err := goKafka.ProcessKafkaMessages(ctx, goKafka.ProcessKafkaMessagesInput{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		ProcessFunc: fn,
	})
	if err != nil {
		log.Errorf("%v", err)
	}
}

func handleConfirmEmailAddressRequest(s service.Service) goKafka.ProcessFunc {
	return func(ctx context.Context, message goKafka.Message) {
		// Convert kafka.Message to Protocol Buffer
		pb := &accountV1.ConfirmEmailAddressRequest{}
		if err := message.Unmarshal(pb); err != nil {
			log.Errorf("failed to unmarshal protobuf from kafka message: %v", err)
		}

		if err := s.ConfirmEmailAddress(ctx, pb); err != nil {
			log.Errorf("failed to process kafka message in transaction: %v", err)
		}
	}
}

func handleConfirmPhoneNumberRequest(s service.Service) goKafka.ProcessFunc {
	return func(ctx context.Context, message goKafka.Message) {
		// Convert kafka.Message to Protocol Buffer
		pb := &accountV1.ConfirmPhoneNumberRequest{}
		if err := message.Unmarshal(pb); err != nil {
			log.Errorf("failed to unmarshal protobuf from kafka message: %v", err)
		}

		if err := s.ConfirmPhoneNumber(ctx, pb); err != nil {
			log.Errorf("failed to process kafka message in transaction: %v", err)
		}
	}
}
