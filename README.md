# Market

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
