package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mitchellh/mapstructure"
)

// Подключение к БД
func BDConnect(connectName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Загрузка данных из БД
func dbLoadCash(db *sql.DB, c *Cache) {
	tArrData := reqDBCache(db)
	for _, tData := range *tArrData {
		if err := tData.Validate(); err != nil {
			fmt.Printf("Error: expected validation data from DB on %v, got %v\n", tData, err)
			return
		}
		c.Set(tData)
	}
}

// Запрос к базе данных
func reqDBCache(db *sql.DB) *[]Data {
	rows, err := db.Query(ReqGetData)
	if err != nil {
		fmt.Printf("Error for request to BD: %v\n", err)
	}
	defer rows.Close()

	tArrData := make([]Data, 0)
	for rows.Next() {
		var tByteData []byte
		if err := rows.Scan(&tByteData); err != nil {
			fmt.Printf("Error read: %v\n", err)
		}
		var tData Data
		if err := json.Unmarshal(tByteData, &tData); err != nil {
			fmt.Printf("Error unmarshalling: %v\n", err)
		}
		tArrData = append(tArrData, tData)
	}

	return &tArrData
}

// Добавление данных в БД
func dbAdd(db *sql.DB, id string, dataItem *[]byte) {
	_, err := db.Exec(ReqAddData, id, dataItem)
	if err != nil {
		fmt.Printf("Error add new data in DB: %v\n", err)
	}
}
