package main

var (
	// Получение всего массива данных из БД
	ReqGetData = `
SELECT data
FROM order_table
`
	// Вставка новых данных в БД
	ReqAddData = `
INSERT INTO order_table (id, data)
VALUES ($1, $2)
`
)
