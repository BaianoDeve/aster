package kafka

import (
	"testing"

	"github.com/BaianoDeve/aster/adapter/presenter/transaction"
	"github.com/BaianoDeve/aster/domain/entity"
	"github.com/BaianoDeve/aster/usecase/process_transaction"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
)

func TestProducerPublish(t *testing.T) {

	expectedOutput := process_transaction.TransactionDtoOutput{
		ID:           "1",
		Status:       entity.REJECTED,
		ErrorMessage: "you exceeded the limit of this transaction",
	}

	// outputJson, _ := json.Marshal(expectedOutput)

	configMap := ckafka.ConfigMap{
		"test.mock.num.brokers": 3,
	}
	producer := NewKafkaProducer(&configMap, transaction.NewTransactionKafkaPresenter())
	err := producer.Publish(expectedOutput, []byte("1"), "test")

	assert.Nil(t, err)
}
