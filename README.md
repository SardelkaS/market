# Market

## Introduction
Simple implementation of an online store.
It has a flexible configuration through the database, 
so it can be used to sell any goods. This version does 
not support online payment, however, it is quite simple to add it.

## Config
```
service:
  host: ''
  port: ''
postgres:
  host: ''
  port: ''
  user: ''
  password: ''
  db_name: ''
redis:
  host: ''
  port: ''
  password: ''
auth:
  refresh_life_time: 0
  access_life_time: 0
  secret: ''
tg_bot:
  token: ""
  chat_id: 0
```

## Run
### Simple
```
go run cmd/api/main.go
```
### Docker-compose
```
docker-compose up --build -d --no-deps
```

## API Documentation
https://documenter.getpostman.com/view/25451223/2s93m32NXq
