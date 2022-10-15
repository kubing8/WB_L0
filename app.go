package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"log"
	"time"
)

// Модель данных
type Data struct {
	OrderUid    string `json:"order_uid" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Entry       string `json:"entry" validate:"required"`
	Delivery    struct {
		Name    string `json:"name" validate:"required"`
		Phone   string `json:"phone" validate:"required"`
		Zip     string `json:"zip" validate:"required"`
		City    string `json:"city" validate:"required"`
		Address string `json:"address" validate:"required"`
		Region  string `json:"region" validate:"required"`
		Email   string `json:"email"`
	} `json:"delivery"`
	Payment struct {
		Transaction  string `json:"transaction" validate:"required"`
		RequestId    string `json:"request_id"`
		Currency     string `json:"currency" validate:"required"`
		Provider     string `json:"provider" validate:"required"`
		Amount       int    `json:"amount" validate:"required"`
		PaymentDt    int    `json:"payment_dt" validate:"required"`
		Bank         string `json:"bank" validate:"required"`
		DeliveryCost int    `json:"delivery_cost" validate:"required"`
		GoodsTotal   int    `json:"goods_total" validate:"required"`
		CustomFee    int    `json:"custom_fee" `
	} `json:"payment"`
	Items []struct {
		ChrtId      int    `json:"chrt_id" validate:"required"`
		TrackNumber string `json:"track_number" validate:"required"`
		Price       int    `json:"price" validate:"required"`
		Rid         string `json:"rid" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Sale        int    `json:"sale"`
		Size        string `json:"size"`
		TotalPrice  int    `json:"total_price" validate:"required"`
		NmId        int    `json:"nm_id" validate:"required"`
		Brand       string `json:"brand" validate:"required"`
		Status      int    `json:"status" validate:"required"`
	} `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey" validate:"required"`
	SmId              int       `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard"`
}

// Валидатор модели данных
var Valid = validator.New()

func (data *Data) Validate() error {
	if err := Valid.Struct(data); err != nil {
		// Проверка на правильность (полноту) полеченных данных
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		return err
	}
	return nil
}

// Струкутра кэша
type Cache struct {
	item map[string]Data
}

func (c *Cache) Get(id string) (Data, bool) {
	data, ok := c.item[id]
	return data, ok
}

func (c *Cache) Set(data Data) {
	c.item[data.OrderUid] = data
}

var refCache *Cache

func dataCache() *Cache {
	return refCache
}

func main() {
	// Данные для подклдчения к базе данных
	connectDatabaseName := "user=postgres password=admin dbname=OrderDB sslmode=disable"
	conBD, err := BDConnect(connectDatabaseName)
	if err != nil {
		fmt.Printf("Connection DB: %v\n", err)
	}
	defer conBD.Close()

	// Инициализация кэша
	iniData := make(map[string]Data)
	var temp = Cache{iniData}

	refCache = &temp // сохранение ссылки на кэш

	// Загружаем кэш из базы данных
	dbLoadCash(conBD, &temp)

	// Подписываемся на канал
	go subsChanel(conBD, &temp)
	// Запуск сервера
	port := 8080               // Устанавливаем порт для нашего сервера
	log.Fatal(runServer(port)) // Запуск сервера на порту: port

}
