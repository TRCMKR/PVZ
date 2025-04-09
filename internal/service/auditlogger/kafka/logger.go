package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

var (
	errClosedChannel = errors.New("channel closed")
)

func logger(ctx context.Context, cfg config.Config, done chan<- models.Log) error {
	partConsumer, err := initConsumer(ctx, cfg)
	if err != nil {
		return err
	}

	log.Print("Starting partition logger")
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				return errClosedChannel
			}

			var receivedLog models.Log
			err = json.Unmarshal(msg.Value, &receivedLog)
			if err != nil {
				log.Print("Error unmarshalling log:", err)
				continue
			}

			done <- receivedLog

			log.Print(receivedLog.String())
		}
	}
}
