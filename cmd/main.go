package main

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/BaianoDeve/aster/adapter/broker/kafka"
	"github.com/BaianoDeve/aster/adapter/factory"
	"github.com/BaianoDeve/aster/adapter/presenter/transaction"
	"github.com/BaianoDeve/aster/usecase/process_transaction"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect to db
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}

	// repository
	repositoryFactory := factory.NewRepositoryDatabaseFactory(db)
	repository := repositoryFactory.CreateTransactionRepository()

	// configMapProducer
	configMapProducer := &ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	}

	// producer
	kafkaPresenter := transaction.NewTransactionKafkaPresenter()
	producer := kafka.NewKafkaProducer(configMapProducer, kafkaPresenter)

	// configMapConsumer
	var msgChan = make(chan *ckafka.Message)
	configMapConsumer := &ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"client.id":         "asterpay",
		"group.id":          "asterpay",
	}

	// topic
	topics := []string{"transactions"}

	// consumer
	consumer := kafka.NewKafkaComsumer(configMapConsumer, topics)
	go consumer.Consumer(msgChan)

	// usecases
	usecase := process_transaction.NewProcessTransaction(
		repository,
		producer,
		"transactions_result",
	)

	for msg := range msgChan {
		var input process_transaction.TransactionDtoInput
		json.Unmarshal(msg.Value, &input)
		usecase.Execute(input)
	}
}
