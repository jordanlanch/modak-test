# 🔍 modak_TEST
This service has signup and login with JWT, also you can get the comsuption by periods

# POSTMAN
./modak-test.postman_collection.json

# 💻 Requirements
  - Port: [8080] - REST
  - make
  - docker version 20.10.21.
  - docker-compose version 1.29.2.
  - Cofee ☕

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

### How to generate the mock code?

```
# Generate mock code for domain
make mocks

```

# Run Tests (unit test and integration test)

```
  make test
```
