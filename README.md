# Example Microservices

# Stack
+ REST API
+ PostgreSQL
+ sqlx
+ slog
+ Docker
+ NATS

# Getting Started
1. `git clone https://github.com/evgkhm/go-nats-pub-sub-restapi-postgresql`
2. `cd go-nats-pub-sub-restapi-postgresql`
3. `docker-compose up --build`

# API
| Endpoint              | Method |          Description |
|-----------------------|:------:|---------------------:|
| get_balance_user/:id  |  GET   |     Get balance user |
| create_user           |  POST  |          Create user |
| accrual_balance_user  |  POST  | Accrual balance user |

# TODO
Add tests, linter, metrics