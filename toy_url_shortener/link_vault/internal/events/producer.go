package events

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventProducer interface {
	PublishClick(ctx context.Context, code string)
	Close() error
}

type KafkaProducer struct {
	writer      *kafka.Writer
	queue       chan string
	workerCount int
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewKafkaProducer(brokers []string, topic string, queueSize int,
	workerCount int) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
	}

	ctx, cancel := context.WithCancel(context.Background())

	p := &KafkaProducer{
		writer:      writer,
		queue:       make(chan string, queueSize),
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
	}

	p.startWorkers()
	return p
}

func (k *KafkaProducer) startWorkers() {
	for i := 0; i < k.workerCount; i++ {
		k.wg.Add(1)
		go k.worker(i)
	}
}

func (k *KafkaProducer) worker(id int) {
	defer k.wg.Done()
	for {
		select {
		case code := <-k.queue:

			msg := kafka.Message{
				Key:   []byte(code),
				Value: []byte(code),
			}
			err := k.writer.WriteMessages(k.ctx, msg)
			if err != nil {
				log.Printf("worker %d kafka write error: %v\n", id, err)
			}
		case <-k.ctx.Done():
			return
		}
	}
}

func (k *KafkaProducer) PublishClick(ctx context.Context, code string) {
	select {
	case k.queue <- code:
		// event queued
	default:
		// queue full → drop event
		log.Println("kafka queue full, dropping click event")
	}
}

func (k *KafkaProducer) Close() error {
	k.cancel()
	close(k.queue)
	k.wg.Wait()
	return k.writer.Close()
}
