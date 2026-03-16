package consumer

import (
	"context"
	"log"

	"github.com/linkanalytics/internal/store"
	"github.com/segmentio/kafka-go"
)

type ClickConsumer struct {
	reader *kafka.Reader
	store  *store.ClickStore
}

func NewClickConsumer(brokers []string, topic string, groupID string, store *store.ClickStore) *ClickConsumer {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		CommitInterval: 0,
	})

	return &ClickConsumer{
		reader: reader,
		store:  store,
	}
}

func (c *ClickConsumer) Start(ctx context.Context) {

	for {
		select {

		case <-ctx.Done():
			log.Println("consumer shutting down")
			return

		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				log.Println("kafka read error:", err)
				continue
			}

			log.Printf("received message: %s\n", string(msg.Value))

			// process message here
			err = c.store.RecordClick(ctx, string(msg.Value))
			if err != nil {
				log.Println("db error:", err)
			}

			// commit offset only after success
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Println("commit error:", err)
			}
		}
	}
}

func (c *ClickConsumer) Close() error {
	return c.reader.Close()
}
