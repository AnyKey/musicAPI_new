package main

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

type Register struct {
	dbConn *sql.DB
	rConn  *redis.Client
	qConn  *amqp.Connection
}

func NewReg(
	dbConn *sql.DB,
	rConn *redis.Client,
	qConn *amqp.Connection,
) *Register {
	return &Register{
		dbConn: dbConn,
		rConn:  rConn,
		qConn:  qConn,
	}
}
