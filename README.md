# API Gateway and Microservices

This project is built on a microservices architecture, enabling user and post management through an API Gateway.

## Project Structure

- `api-gateway/` - Handles API Gateway requests and routes them to `user-service` and `post-service`.
- `user-service/` - A microservice handling CRUD operations related to users.
- `post-service/` - A microservice handling CRUD operations related to posts.

---

## API Endpoints

### 1. API Gateway

#### User Service Endpoints
| Method | Endpoint | Description |
|--------|---------|-------------|
| POST | `/users/register` | User registration |
| POST | `/users/login` | User login |
| GET  | `/users/:id` | Retrieve user details |
| PUT  | `/users/:id` | Update user information |
| DELETE | `/users/:id` | Delete a user |

#### Post Service Endpoints
| Method | Endpoint | Description |
|--------|---------|-------------|
| POST | `/posts` | Create a new post |
| GET  | `/posts` | Retrieve all posts |
| GET  | `/posts/:id` | Retrieve a single post |
| PUT  | `/posts/:id` | Update a post |
| DELETE | `/posts/:id` | Delete a post |

---

## Running the Project

The project can be run using Docker Compose. Use the following command to start all services:

```sh
docker-compose up --build
```

To run a specific service, use:

```sh
docker-compose up api-gateway
```

---

## Technologies Used
- **Go (Golang)** - Primary programming language for backend services
- **Gin** - Web framework used for HTTP API Gateway and services
- **gRPC** - Communication between microservices
- **Docker & Docker Compose** - Service orchestration
- **MongoDB & PostgreSQL** - Databases used for data storage

---

## Author
This project was developed by [Your Name].
