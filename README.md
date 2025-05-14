
# Go Pack Optimizer

Go Pack Optimizer is an application built in Go that allows optimizing the usage of packs to accommodate requested items. The application exposes a RESTful API that interacts with a simple frontend to visualize the pack size options and calculate the number of packs required.

This project includes:
- **Go Backend**: API that handles the pack optimization logic.
- **Go Frontend**: Simple interface to interact with the API and display the results.
- **Database (LevelDB)**: Stores the available pack sizes.

## Architecture

The project is divided into two main parts:

1. **Backend**: Exposes a RESTful API to manage packs and calculate the number of packs required.
2. **Frontend**: User interface to view available packs, add/remove packs, and calculate the necessary packs for a given number of items.

## Live Demo

Access the complete demonstration of the application at: [https://ns3365216.ip-37-187-75.eu/](https://ns3365216.ip-37-187-75.eu/)

You're welcome to explore its features.


## Installation

### Prerequisites

- Go 1.23 or higher.
- Docker (optional, if you want to run the project in containers).
- `docker-compose` (if using containers).

### Steps to run the project locally:

1. Clone the repository:
    ```bash
    git clone https://github.com/jmsilvadev/go-pack-optimizer.git
    cd go-pack-optimizer
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. To run the backend locally:
    ```bash
    cd cmd/backend
    go run main.go
    ```

4. To run the frontend locally:
    ```bash
    cd cmd/frontend
    go run main.go
    ```

### Docker

If you prefer running the project in Docker containers, you can use the provided `docker-compose.yml` file.

1. **Build and Run**:
    - To build the Docker images:
      ```bash
      make up-build
      ```
      This command will build Docker images for the frontend and backend and start the containers.

    - To start the containers (without rebuilding):
      ```bash
      make up
      ```

    - To stop the containers:
      ```bash
      make down
      ```

2. **Access the Frontend Locally**:
   - The frontend interface will be available at: [http://localhost:3000](http://localhost:3000)

3. **Access the Backend Locally**:
   - The backend API will be available at: [http://localhost:8080](http://localhost:8080)

---

## How to Use

### Backend API

The API provides endpoints to interact with the pack system. Here are the main endpoints:

- **GET /v1/packs**: Returns all available packs.
- **POST /v1/packs**: Adds a new pack size.
- **DELETE /v1/packs/{size}**: Removes an existing pack size.
- **POST /v1/order**: Calculates the best combination of packs to use

### Frontend

The frontend is a simple interface that allows you to:

- **View the list of available packs**.
- **Add new packs**.
- **Remove packs**.
- **Calculate the number of packs needed for a given number of items**.

---

## OAS (OpenAPI Specification)

The OpenAPI Specification (OAS) for the API is available [here](oas.yaml). This document outlines all the available endpoints, request parameters, and response formats.

---

## Generating Docker Images

To build Docker images for the frontend and backend:

1. **Frontend Image**: The frontend Dockerfile (`Dockerfile_frontend`) creates an image for the frontend, which you can build using the following command:
   ```bash
   make build-image-frontend
   ```

2. **Backend Image**: The backend Dockerfile (`Dockerfile_backend`) builds the backend image:
   ```bash
   make build-image-backend
   ```

### Build the Complete Docker Image

To build both frontend and backend images and run them using Docker Compose:
```bash
make up-build
```

### Tests
To run the tests, use the following command:

```sh
make tests
```

Showing coverage:
```sh
make tests-cover
```


### Commands
```
$ make

build-backend                  Build backend component
build-frontend                 Build frontend component
build-image-backend            Build docker image in daemon mode
build-image-frontend           Build Docker image for frontend using Dockerfile_frontend
clean                          Clean all builts
clean-tests                    Clean tests
down                           Stop docker container
logs                           Watch docker log files
tests-cover                    Run tests with coverage
tests                          Run unit tests
up-build                       Start docker container and rebuild the image
up                             Start docker container

```


## Improvements

Here are some ideas for future enhancements to make this system more scalable, observable, and production-ready:

### Distributed Database
Replace the embedded LevelDB with a distributed database (e.g., PostgreSQL, CockroachDB, or Redis) to allow multi-instance deployments and data persistence across replicas.

### Observability
Integrate observability tools:
- **Metrics** with Prometheus
- **Tracing** with OpenTelemetry
- **Logging** aggregation using tools like Fluent Bit or Loki
This will help monitor the health and performance of the services.

### Infrastructure as Code
Introduce **Terraform** to define cloud infrastructure as code, making deployments consistent and reproducible across environments (Dev, QA, Prod).

### Circuit Breakers
Implement **circuit breakers** (e.g., with [goresilience](https://github.com/slok/goresilience)) in case this service depends on external APIs or databases, ensuring graceful degradation in case of failures.

### Memoization Improvements
Enhance memoization strategy with:
- TTL-based in-memory cache (e.g., using [ristretto](https://github.com/dgraph-io/ristretto))
- Distributed cache support (e.g., Redis) for shared environments
- Auto-invalidation after pack updates or deletions

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
