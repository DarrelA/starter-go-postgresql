- [Design Considerations](#design-considerations)
  - [Requirements](#requirements)
  - [Model-View-Controller (MVC)](#model-view-controller-mvc)
  - [N-tier Architecture](#n-tier-architecture)
  - [Combining Microservices and N-tier Architectures](#combining-microservices-and-n-tier-architectures)
  - [Middlewares](#middlewares)
    - [Logging and Monitoring](#logging-and-monitoring)
      - [`X-Correlation-ID` \& `Request-ID`](#x-correlation-id--request-id)
  - [JWT](#jwt)
    - [Implementation of Refresh Token with Redis](#implementation-of-refresh-token-with-redis)
- [Setup](#setup)
  - [Handle Initial Files](#handle-initial-files)
  - [Generate the Private and Public Keys](#generate-the-private-and-public-keys)
- [Maintenance](#maintenance)
  - [`go.mod` File](#gomod-file)
  - [psql](#psql)
  - [golang-migrate](#golang-migrate)
  - [redis](#redis)
- [References](#references)

<a name="readme-top"></a>

# Design Considerations

## Requirements

| No. |                     Feature                      |                                                Reason                                                |
| :-: | :----------------------------------------------: | :--------------------------------------------------------------------------------------------------: |
| 1.  | Easy setup using Env files, Docker, and Makefile |                       Simplifies initial setup and environment configuration.                        |
| 2.  | Ease of maintenance of Go packages using Go Test |                      Ensures code reliability and simplifies testing processes.                      |
| 3.  |           Model-View-Controller (MVC)            |                                   [🔗](#model-view-controller-mvc)                                   |
| 4.  |               N-tier Architecture                |                                      [🔗](#n-tier-architecture)                                      |
| 5.  |    Implementation of Refresh Token with Redis    |                          [🔗](#implementation-of-refresh-token-with-redis)                           |
| 6.  |                   Middlewares                    |                                          [🔗](#middlewares)                                          |
| 7.  |     Fiber as a web framework for performance     |                    Optimized for minimal memory allocation and high performance.                     |
| 8.  |   Internal logging for security using zerolog    | Provides secure and efficient logging mechanisms,</br>including structured logging with JSON output. |
| 9.  |   A common folder structure for project layout   |                      Facilitates scalability into a microservice architecture.                       |

## Model-View-Controller (MVC)

The Model-View-Controller (MVC) design pattern separates an application into three interconnected components, improving the organization and manageability of the code.

1. **Model:** Manages the data and business logic of the application.
   - DTO files are located in the `internal/domains` folder.
2. **View:** Represents the user interface and displays the data.
   - Refer to the frontend repository.
3. **Controller:** Acts as an intermediary between the Model and View, processing user input and updating the Model and View accordingly.
   - Handlers process HTTP requests and are located in `internal/handlers`.
   - Services contain business logic and are located in `internal/services`.
   - DAO files are in the `internal/domains` folder for database interactions.

## N-tier Architecture

An N-tier architecture separates each layer or tier of the application both physically and logically, which enhances flexibility and manageability.

1. **Frontend (Presentation Layer):** This is the user interface of the application. In an N-tier setup, you can update or replace the frontend without needing to alter the backend, as long as the interface contract (e.g., APIs) remains consistent.

2. **Backend (Business Logic Layer):** This layer contains the business rules and logic. You can modify or replace this layer as needed without affecting the other tiers, provided the interfaces between the layers don't change.

3. **Database (Data Access Layer):** This layer handles data storage and retrieval. You can switch out databases or change how data is accessed, and as long as you maintain the same data contracts or APIs, the other layers remain unaffected.

Benefits of N-tier architecture include:

1. **Scalability:** Each layer can be scaled independently based on its resource requirements.
2. **Flexibility:** You can replace or upgrade one layer without significant rework of the others.
3. **Maintainability:** Smaller, well-defined codebases for each layer are easier to manage and understand.

However, N-tier architecture can also introduce complexity in terms of network latency, configuration, and management. Each layer likely needs to communicate over a network, which can affect performance and add complexity that wouldn't be as prominent in a monolithic architecture.

## Combining Microservices and N-tier Architectures

Combining N-tier architecture with microservices involves separating concerns within an application while decoupling it into independently deployable services. Each microservice can adhere to N-tier principles, enhancing modularity and scalability.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Middlewares

| **Middleware**                                        | **Use Case(s)**                                                                                                                                                    |
| ----------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Handling CORS<br/>(Cross-Origin Resource Sharing)** | Allowing or restricting resources on a web server to be requested from another domain.<br/>Configuring CORS policies to secure and control access to resources.    |
| **Authentication and Authorization**                  | Ensuring users are authenticated before accessing certain resources.<br/>Validating user permissions and roles to control access to specific endpoints or actions. |
| **Logging and Monitoring**                            | Tracking and recording requests and responses for auditing and debugging.<br/>Monitoring application performance and detecting issues or anomalies.                |
| **Data Validation and Sanitization**                  | Ensuring incoming request data adheres to expected formats and rules.<br/>Removing or escaping potentially harmful data to prevent injection attacks.              |
| **Session Management**                                | Maintaining user session state across multiple requests.<br/>Handling session creation, expiration, and validation.                                                |

### Logging and Monitoring

| **Information**             | **Use Case (Software Engineers)**                                                                                                                   | **Use Case (DevOps)**                                                                                                                                             |
| --------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Host/</br>Server Info**   | Identifying issues related to specific servers or instances in a distributed environment.                                                           | Monitoring server health and performance to identify potential issues and balance load across instances.                                                          |
| **Request Method**          | Identifying the type of HTTP request (GET, POST, etc.) to understand user actions and flow through the application.                                 | Monitoring the distribution of different request types to identify unusual patterns or spikes that may indicate issues.                                           |
| **HTTP Referrer**           | Understanding the sources of traffic and user navigation patterns to optimize user experience and marketing efforts.                                | Tracking referral sources to identify potential security threats, such as phishing or spam attacks, and to optimize resource allocation based on traffic sources. |
| **Request URL**             | Debugging issues related to specific endpoints and understanding which resources are being accessed.                                                | Tracking usage patterns across different endpoints to optimize performance and resource allocation.                                                               |
| **Status Code**             | Identifying success or failure of requests to pinpoint potential issues in the application logic.                                                   | Monitoring overall application health and identifying trends in error rates or unusual spikes in specific status codes.                                           |
| **Response Time**           | Measuring the performance of individual requests to identify slow responses and potential bottlenecks in the application.                           | Ensuring the application meets performance SLAs and identifying areas for optimization to improve response times.                                                 |
| **Latency**                 | Measuring the time taken for different parts of the request processing to identify and optimize slow components.                                    | Monitoring end-to-end latency to ensure the application meets performance SLAs and identifying areas for optimization.                                            |
| **IP Address**              | Tracking user activity and identifying potential malicious activities or misuse of the application.                                                 | Monitoring the geographic distribution of traffic and identifying potential security threats such as DDoS attacks.                                                |
| **User Agent**              | Understanding the types of clients (browsers, mobile devices, etc.) accessing the application to ensure compatibility and optimize user experience. | Tracking the mix of client types to plan for capacity and performance testing, as well as identifying potentially malicious clients.                              |
| **Correlation ID**          | Tracing requests through the entire application stack to diagnose issues and ensure consistency across services.                                    | Linking logs from different services and components to get a holistic view of the application's behavior and identify root causes of issues.                      |
| **Request ID**              | Tracking and debugging specific requests within a service to identify issues and optimize request handling.                                         | Monitoring individual request logs to identify patterns, performance issues, and potential security threats at a granular level.                                  |
| **User ID/</br>Session ID** | Tracking user activity and debugging issues related to specific users or sessions.                                                                  | Monitoring user behavior and session patterns to identify potential misuse or performance issues affecting specific user groups.                                  |

<p align="right">(<a href="#readme-top">back to top</a>)</p>

#### `X-Correlation-ID` & `Request-ID`

- **`X-Correlation-ID`**
  - **Purpose:** Used to trace and correlate a series of related requests across multiple services. It helps in understanding the journey of a particular transaction or user action through the system.
  - **Scope:** Typically spans multiple services and systems. It's used to connect different parts of a transaction that might pass through various microservices.
  - **Usage:** The `X-Correlation-ID` header is included in HTTP requests and responses to maintain the correlation ID across different services.
- **`Request-ID`**
  - **Purpose:** Used to uniquely identify a single request within a service. It helps in logging and debugging individual requests.
  - **Scope:** Usually limited to a single request-response cycle within a service. It can be useful for tracking the processing of a request through different components of the same service.
  - **Usage:** The `Request-ID` header is included in HTTP requests and responses to identify the specific request.
- **Key Differences**
  - **Scope:** `X-Correlation-ID` has a broader scope, typically spanning multiple services, while `Request-ID` is more narrowly focused on a single request within a service.
  - **Usage Scenario:** Use `X-Correlation-ID` when you need to trace an end-to-end transaction through multiple microservices. Use `Request-ID` for tracking and debugging individual requests within a specific service.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## JWT

### Implementation of Refresh Token with Redis

JWTs are stateless and should not require database round trips. Redis improves token validation efficiency by reducing database queries and offers built-in expiration times for stored keys, which is ideal for managing tokens with limited lifespans. It effectively handles token revocation by invalidating tokens and preventing unauthorized access without needing additional database queries. By leveraging Redis, we achieve a seamless and secure authentication process.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Setup

## Handle Initial Files

1. **Respective `.env` files in `configs` folder**
2. **`init.sql` script:** Creates initial schemas
3. **`servers.json` file:** Establishes server connection from pgAdmin to Postgres

## Generate the Private and Public Keys

1. [Online RSA Key Generator](https://travistidwell.com/jsencrypt/demo/): Key Size: 2048 bit
2. [BASE64 Decode and Encode](https://www.base64encode.org/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Maintenance

## `go.mod` File

```sh
# Updating `go.mod`
go get -u
go mod tidy
```

## psql

```sh
# pgAdmin
http://localhost:5050/browser/

# If not using pgAdmin
docker exec -it postgres bash
su - postgres
psql -U user1 -d pgstarter
SELECT * FROM users;
```

## golang-migrate

```sh
# To work with the scripts in db/migration folder
brew install golang-migrate
migrate create -ext sql -dir db/migration/ -seq init_schema

# Format
# {version}_{title}.down.sql
# {version}_{title}.up.sql
```

## redis

```sh
docker exec -it redis /bin/sh
redis-cli
INFO

KEYS *
# Get the value of the key (user_uuid)
GET key
# Check the remaining time to live (TTL) of a key
TTL key
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# References

- [go-project-layout](https://appliedgo.com/blog/go-project-layout)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout/tree/master)
- [Go postgres driver](https://github.com/lib/pq) recommends [pgx - PostgreSQL Driver and Toolkit](https://github.com/jackc/pgx)
- [wpcodevo/golang-fiber-jwt-rs256](https://github.com/wpcodevo/golang-fiber-jwt-rs256)
- [Should I use UUID as well as ID](https://dba.stackexchange.com/questions/115766/should-i-use-uuid-as-well-as-id)
- [Why your software should use UUIDs](https://devforth.io/blog/why-your-software-should-use-uuids/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---
