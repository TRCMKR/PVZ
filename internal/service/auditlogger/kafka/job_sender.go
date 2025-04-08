package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/IBM/sarama"
)

func jobSender(ctx context.Context, cfg config.Config, jobs <-chan models.Log, failed chan<- models.Log) error {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{cfg.KafkaHost() + ":" + cfg.KafkaPort()}, kafkaConfig)
	if err != nil {
		return err
	}
	defer producer.Close()

	log.Print("Starting jobSender")
	for job := range jobs {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		jobBytes, err := json.Marshal(job)
		if err != nil {
			log.Printf("Error marshaling JSON: %v\n", err)

			continue
		}

		message := &sarama.ProducerMessage{
			Topic: "logs",
			Key:   sarama.StringEncoder(strconv.Itoa(job.ID)),
			Value: sarama.ByteEncoder(jobBytes),
		}

		_, _, err = producer.SendMessage(message)
		if err != nil {
			failed <- job
		}
	}

	return nil
}
