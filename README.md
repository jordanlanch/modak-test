# 🔍 modak_TEST
This service control de notification and limit rate by MessageType

# Enpoint
```
  curl --location 'http://localhost:8080/api/notification' \
--header 'Content-Type: application/json' \
--data-raw '{
    "recipient": "test@test.com",
    "message_type": "Marketing",
    "content": "content"
}'
```


# 💻 Requirements
  - Port: [8080] - REST
  - make
  - docker version 20.10.21.
  - docker-compose version 1.29.2.

copy .env.example to .env
# 🚀 Run the app
```sh
$ ./run-dev.sh
```
## Architecture Layers of the project

- Router
- Controller
- Usecase
- Repository
- Domain


# Run Tests (unit test and integration test)

```
  make test
```
