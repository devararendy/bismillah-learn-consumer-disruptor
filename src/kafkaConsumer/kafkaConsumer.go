package kafkaconsumer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func DeInit() {
	fmt.Println("This is DeInit of Kafka Consumer")

}

type Order struct {
	Token string
	NIK   string
	Topup int
}

func Init() {
	topics := os.Getenv("KAFKA_TOPICS")

	c, e := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        os.Getenv("KAFKA_BROKERS"),
		"go.events.channel.enable": true,
		"group.id":                 os.Getenv("KAFKA_GROUP_ID"),
		"auto.offset.reset":        "earliest",
		"enable.partition.eof":     true,
		// "partition.assignment.strategy": "range",
	})

	if e != nil {
		if ke, ok := e.(kafka.Error); ok == true {
			switch ec := ke.Code(); ec {
			case kafka.ErrInvalidArg:
				fmt.Printf("😢 Can't create the Consumer because you've configured it wrong (code: %d)!\n\t%v\n\nTo see the configuration options, refer to https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md", ec, e)
			default:
				fmt.Printf("😢 Can't create the Consumer (Kafka error code %d)\n\tError: %v\n", ec, e)
			}
		} else {
			// It's not a kafka.Error
			fmt.Printf("😢 Oh noes, there's a generic error creating the Consumer! %v", e.Error())
		}
	} else {

		// Subscribe to the topic
		if e := c.Subscribe(topics, nil); e != nil {
			fmt.Printf("☠️ Uh oh, there was an error subscribing to the topic :\n\t%v\n", e)
		} else {
			doTerm := false
			for !doTerm {
				if ev := c.Poll(1000); ev == nil {
					// the Poll timed out and we got nothin'
					// fmt.Printf("……\n")
					continue
				} else {
					// Look at the type of Event we've received
					switch ev.(type) {

					case *kafka.Message:
						// It's a message
						km := ev.(*kafka.Message)
						fmt.Printf("✅ Message received from topic '%v' (partition %d at offset %d)\n",
							string(*km.TopicPartition.Topic),
							km.TopicPartition.Partition,
							km.TopicPartition.Offset)
						// fmt.Printf("✅ Message '%v' received from topic '%v' (partition %d at offset %d)\n",
						// 	string(km.Value),
						// 	string(*km.TopicPartition.Topic),
						// 	km.TopicPartition.Partition,
						// 	km.TopicPartition.Offset)

						// decode into struct
						//decode back from byte array to struct Order
						networkIn := bytes.NewBuffer(km.Value)
						dec := gob.NewDecoder(networkIn)
						var recvOrder Order
						err := dec.Decode(&recvOrder)
						if err != nil {
							log.Fatal("decode error:", err)
						} else {
							fmt.Printf("Decode Struct : %+v\n", recvOrder)
						}
					case kafka.PartitionEOF:
						// We've finished reading messages on this partition so let's wrap up
						// n.b. this is a BIG assumption that we are only consuming from one partition
						pe := ev.(kafka.PartitionEOF)
						fmt.Printf("🌆 Got to the end of partition %v on topic %v at offset %v\n",
							pe.Partition,
							string(*pe.Topic),
							pe.Offset)
						// doTerm = true
					case kafka.OffsetsCommitted:
						continue
					case kafka.Error:
						// It's an error
						em := ev.(kafka.Error)
						fmt.Printf("☠️ Uh oh, caught an error:\n\t%v\n", em)
					default:
						// It's not anything we were expecting
						fmt.Printf("Got an event that's not a Message, Error, or PartitionEOF 👻\n\t%v\n", ev)

					}
				}

			}
			fmt.Printf("👋 … and we're done. Closing the consumer and exiting.\n")
			// Now we can exit
			c.Close()
		}

	}

	// defer DeInit()
	fmt.Println("This is Init of Kafka Consumer")

}
