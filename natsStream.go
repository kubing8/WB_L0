package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"sync"
)

// Данные для подключения к каналу nats-streaming
const clusterID = "test-cluster"
const clientID = "ser1"

// Подключение к каналу
func subsChanel(db *sql.DB, temp *Cache) {
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		fmt.Printf("Error connect to nuts-streaming: %v\n", err)
	}
	defer sc.Close()

	sc.Subscribe("nats_test_stream", func(m *stan.Msg) {
		var unmarData Data
		if err := json.Unmarshal(m.Data, &unmarData); err != nil {
			fmt.Printf("Error, can't unmarshal data from Nats: %v\n", err)
			m.Ack() // Подтверждаем, что данные получены
			return
		}

		if _, ok := temp.Get(unmarData.OrderUid); ok == true {
			m.Ack() // Подтверждаем, что данные получены
			fmt.Printf("Duplicate data, uid: %v\n", unmarData.OrderUid)
			return
		}

		if err := unmarData.Validate(); err != nil {
			fmt.Printf("Error: expected validation data from Nats on %v, got %v\"", unmarData, err)
			m.Ack() // Подтверждаем, что данные получены
			return
		}

		dbAdd(db, unmarData.OrderUid, &m.Data)
		temp.Set(unmarData)

		m.Ack() // Подтверждаем, что данные получены
		fmt.Printf("Received a message: %s\n", unmarData.OrderUid)
	}, stan.DurableName("my-durable"), stan.SetManualAckMode())

	Block()
}

func Block() {
	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}
