<a name="readme-top"></a>

- [Design Considerations](#design-considerations)
  - [Requirements](#requirements)
  - [Model-View-Controller (MVC)](#model-view-controller-mvc)
  - [N-tier Architecture](#n-tier-architecture)
  - [Combining Microservices and N-tier Architectures](#combining-microservices-and-n-tier-architectures)
  - [Middlewares](#middlewares)
    - [Logging and Monitoring](#logging-and-monitoring)
      - [Custom Headers: `Correlation-ID` and `Request-ID`](#custom-headers-correlation-id-and-request-id)
  - [JWT](#jwt)
    - [Implementation of Refresh Token with Redis](#implementation-of-refresh-token-with-redis)
      - [Redis and JWT Statelessness](#redis-and-jwt-statelessness)
      - [What are Access and Refresh Tokens?](#what-are-access-and-refresh-tokens)
      - [How Does It Work?](#how-does-it-work)
      - [Managing Security Threats: CSRF and XSS Attacks](#managing-security-threats-csrf-and-xss-attacks)
      - [Security Implementation Details](#security-implementation-details)
      - [Workflow Overview](#workflow-overview)
      - [Diagram: Token Workflow](#diagram-token-workflow)
- [Setup](#setup)
  - [Handle Initial Files](#handle-initial-files)
  - [Generate the Private and Public Keys](#generate-the-private-and-public-keys)
- [Maintenance](#maintenance)
  - [`go.mod` File](#gomod-file)
  - [psql](#psql)
  - [golang-migrate](#golang-migrate)
  - [redis](#redis)
- [Testing](#testing)
  - [Unit Testing and Acceptance Testing](#unit-testing-and-acceptance-testing)
  - [Dependency Injection](#dependency-injection)
- [References](#references)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

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

<p align="right">(<a href="#readme-top">back to top</a>)</p>

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

#### Custom Headers: `Correlation-ID` and `Request-ID`

- **`Correlation-ID`**
  - **Purpose:** Used to trace and correlate a series of related requests across multiple services. It helps in understanding the journey of a particular transaction or user action through the system.
  - **Scope:** Typically spans multiple services and systems. It's used to connect different parts of a transaction that might pass through various microservices.
  - **Usage:** The `Correlation-ID` header is included in HTTP requests and responses to maintain the correlation ID across different services.
- **`Request-ID`**
  - **Purpose:** Used to uniquely identify a single request within a service. It helps in logging and debugging individual requests.
  - **Scope:** Usually limited to a single request-response cycle within a service. It can be useful for tracking the processing of a request through different components of the same service.
  - **Usage:** The `Request-ID` header is included in HTTP requests and responses to identify the specific request.
- **Key Differences**
  - **Scope:** `Correlation-ID` has a broader scope, typically spanning multiple services, while `Request-ID` is more narrowly focused on a single request within a service.
  - **Usage Scenario:** Use `Correlation-ID` when you need to trace an end-to-end transaction through multiple microservices. Use `Request-ID` for tracking and debugging individual requests within a specific service.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## JWT

### Implementation of Refresh Token with Redis

JWTs are stateless and eliminate the need for frequent database queries. Redis enhances token validation efficiency by reducing database queries and providing built-in expiration times for stored keys, making it ideal for managing tokens with limited lifespans. Redis handles token revocation effectively by invalidating tokens and preventing unauthorized access without additional database queries. Leveraging Redis achieves a seamless and secure authentication process.

#### Redis and JWT Statelessness

Storing token details in Redis does not violate the statelessness principle of JWTs (JSON Web Tokens). Here's why:

1. **JWT Statelessness**: JWTs are considered stateless because they contain all the necessary information within the token itself, including claims about the user and any other metadata needed for verification. Statelessness refers to the server's ability to verify the token without needing to store session information or additional user data on the server side.

2. **Redis as a Cache**: Using Redis in this context is typically for performance optimization, not for maintaining the state of a session. Token details are stored in Redis to quickly access them for operations like token revocation or checking expiration, rather than maintaining user session state.

3. **Optional Use of Redis**: Redis is an optional layer for efficiency. JWTs can still be verified independently of Redis, as they contain all the necessary information. Redis enhances performance (due to its in-memory capabilities) and manages token lifecycles (like handling expiration more efficiently), but it does not fundamentally change the stateless nature of JWT authentication.

4. **Token Revocation**: When using Redis for storing tokens (or token blacklists), it is typically for the purpose of token revocation. This is an added security measure that doesn't conflict with JWT statelessness. The tokens remain self-contained, but there is an additional mechanism to invalidate them if necessary.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

#### What are Access and Refresh Tokens?

**Access Tokens** are short-term keys that grant users access to specific resources or data, typically valid for a short period (e.g., 15 minutes). When users log in, they receive an access token, which is included in their requests to access protected resources.

**Refresh Tokens** are long-term keys that allow users to request new access tokens without logging in again. These tokens have a longer lifespan (e.g., 60 minutes) and keep users logged in by providing new access tokens when the old ones expire.

#### How Does It Work?

1. **User Login:**
   During login, credentials (username and password) are verified. If valid, the system generates both an access token and a refresh token. The access token grants immediate access, while the refresh token is stored securely and used to obtain new access tokens when needed.

2. **Accessing Resources:**
   Every time users want to access protected resources, a request is sent with the access token. The server verifies the token to ensure it is valid and not expired. If valid, access to the resource is granted.

3. **Refreshing Tokens:**
   When an access token expires, there is no need to log in again. Instead, the refresh token is used to request a new access token. The server verifies the refresh token and, if valid, issues a new access token. This process keeps the user experience smooth and uninterrupted.

4. **Storing Tokens in Redis:**
   To manage sessions efficiently, token metadata (like expiration times) is stored in Redis. Redis, an in-memory data store known for its speed, saves token metadata, allowing the server to quickly verify tokens and check their status. This setup also allows easy invalidation of tokens if needed (e.g., upon user logout).

5. **Using Cookies for Tokens:**
   To keep things secure and simple, access and refresh tokens are stored in cookies. Cookies are small pieces of data stored on the user's browser. Marking these cookies as HTTP-only ensures they are not accessible via JavaScript, adding an extra layer of security.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

#### Managing Security Threats: CSRF and XSS Attacks

**Cross-Site Request Forgery (CSRF):**

CSRF attacks occur when a malicious website tricks a user's browser into performing an unwanted action on a different site where the user is authenticated. Protection against CSRF attacks involves a combination of techniques:

- **SameSite Cookies:** Setting the `SameSite` attribute of cookies to `Strict` or `Lax` prevents the browser from sending cookies along with cross-site requests. This helps mitigate CSRF attacks by ensuring that cookies are only sent in a first-party context.
- **CSRF Tokens:** In addition to SameSite cookies, CSRF tokens are generated and included in forms and API requests (using fetch or Axios). These tokens are validated on the server side to ensure the request is legitimate.

**Cross-Site Scripting (XSS):**

XSS attacks involve injecting malicious scripts into web pages viewed by other users. These scripts can steal cookies, session tokens, or other sensitive information. Protection against XSS attacks involves the following best practices:

- **Escaping User Input:** Properly escape any user input rendered in the HTML to prevent the execution of malicious scripts.
- **Content Security Policy (CSP):** Implementing a CSP restricts the sources from which scripts, styles, and other resources can be loaded, preventing unauthorized script execution.
- **HTTP-Only Cookies:** Storing tokens in HTTP-only cookies ensures they are not accessible via JavaScript, preventing malicious scripts from stealing tokens.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

#### Security Implementation Details

**Frontend:**

- Use fetch API or Axios for making API requests.
- Include CSRF tokens in headers for API requests.
- Set `SameSite` and `HTTPOnly` attributes for cookies.

**Backend:**

- Generate and validate CSRF tokens for forms and API requests.
- Escape all user inputs rendered in responses.
- Implement a Content Security Policy (CSP) to restrict script sources.
- Store tokens in HTTP-only cookies to prevent client-side access.

**Deployment on Server:**

- Ensure the web server is configured to set secure headers, including CSP and `Strict-Transport-Security`.
- Use HTTPS to encrypt all communications between the client and server.
- Regularly update dependencies and apply security patches to the server and application.

#### Workflow Overview

1. **User Logs In:**

   - Credentials are verified.
   - Access and refresh tokens are generated.
   - The access token is used for immediate access; the refresh token is stored securely.
   - Tokens are stored in HTTP-only cookies on the user's browser.

2. **Accessing Resources:**

   - The user sends a request with the access token (automatically included by the browser from the cookie).
   - The server verifies the token using Redis.
   - If valid, access is granted.

3. **Token Refresh:**

   - The access token expires.
   - The user requests a new access token using the refresh token (automatically included by the browser from the cookie).
   - The server verifies the refresh token using Redis.
   - If valid, a new access token is issued and stored in a cookie.

4. **Managing Tokens with Redis:**

   - Token metadata is stored in Redis.
   - Redis allows fast verification and easy invalidation of tokens.

5. **Mitigating CSRF and XSS Attacks:**
   - Use SameSite and HTTP-only cookies.
   - Implement CSRF tokens for form submissions and API requests using fetch or Axios.
   - Escape user input and enforce a Content Security Policy (CSP).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

#### Diagram: Token Workflow

Below is a simplified diagram to illustrate the token workflow:

```plaintext
       +-----------------+         +--------------------+
       | User Logs In    |         | Server Verifies    |
       |                 |<------->| Credentials        |
       +-----------------+         +--------------------+
                |                              |
                v                              v
      +-----------------+         +----------------------+
      | Access Token    |         | Refresh Token        |
      | (Short-lived)   |<------->| (Long-lived)         |
      +-----------------+         +----------------------+
                |                              |
                v                              v
      +-----------------+         +----------------------+
      | Access Resource |         | Request New Access   |
      | (Using Access   |<------->| Token Using Refresh  |
      | Token)          |         | Token                |
      +-----------------+         +----------------------+
                |                              |
                v                              v
      +-----------------+         +----------------------+
      | Verify Token    |<------->| Verify Refresh Token |
      | Using Redis     |         | Using Redis          |
      +-----------------+         +----------------------+
```

This flow ensures a seamless user experience while maintaining high security standards. Access tokens provide quick access to resources, while refresh tokens ensure continuous access without frequent logins. Redis enhances the process by managing tokens efficiently and allowing quick verifications and invalidations. Storing tokens in HTTP-only and SameSite cookies keeps them secure from client-side attacks.

In summary, using access and refresh tokens with Redis, PostgreSQL, and cookies makes the authentication process both secure and user-friendly. Implementing measures to mitigate CSRF and XSS attacks on the frontend, backend, and during deployment ensures that user sessions are not only efficient but also highly secure.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Setup

## Handle Initial Files

1. **Respective `.env` files in `configs` folder**
2. **`init.sql` script:** Creates initial schemas
3. **Respective env server `.json` file:** Establishes server connection from pgAdmin to Postgres

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

# Remove all unused containers, networks, images (both dangling and unreferenced), and optionally, volumes.
# If you want to skip the confirmation prompt, you can add the -f flag:
docker system prune -a --volumes -f
```

## psql

```sh
# pgAdmin
http://localhost:5050/browser/

# If not using pgAdmin
docker exec -it postgres bash

# List all databases
psql -U user1 -d postgres
\l

# View table
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

# Testing

```sh
make APP_ENV=test

go test -cover
go test ./test
```

## Unit Testing and Acceptance Testing

| Aspect              | Unit Testing                        | Acceptance Testing                          |
| ------------------- | ----------------------------------- | ------------------------------------------- |
| **Purpose**         | Verify individual units/components  | Verify the entire system meets requirements |
| **Scope**           | Small, isolated pieces of code      | Entire application or major features        |
| **Nature**          | Typically automated, quick to run   | Can be automated or manual, longer to run   |
| **Examples**        | Testing a single function or method | Testing end-to-end user scenarios           |
| **Tools**           | Go `testing` package                | godog                                       |
| **Characteristics** | White-box testing, high frequency   | Black-box testing, less frequent            |

[User Stories] become [Acceptance Tests] which is [Behavior Driven Development] "Doing the RIGHT thing."

[Code Functionality] becomes [Unit Testing] which is [Test Driven Development] "Doing the THING right."

## Dependency Injection

Dependency Injection (DI) is a design pattern that achieves Inversion of Control (IoC) by allowing components to receive their dependencies from an external source rather than creating them internally. In Go, DI is typically implemented by passing dependencies as parameters to functions or struct constructors. This approach decouples application components, making them more modular, easier to test, and maintain. DI allows replacing real implementations of dependencies with mock versions during testing, ensuring components can be tested in isolation.

- **Promotes Loose Coupling**: DI promotes loose coupling and modularity.
- **Defines Clear Contracts**: Interfaces define clear contracts for components.
- **Encapsulates Related Data**: Structs encapsulate related data and behaviors.
- **Enables Effective Mocking**: Together, these concepts enable effective mocking and testing in Go.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# References

- [go-project-layout](https://appliedgo.com/blog/go-project-layout)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout/tree/master)
- [Go postgres driver](https://github.com/lib/pq) recommends [pgx - PostgreSQL Driver and Toolkit](https://github.com/jackc/pgx)
- [wpcodevo/golang-fiber-jwt-rs256](https://github.com/wpcodevo/golang-fiber-jwt-rs256)
- [Should I use UUID as well as ID](https://dba.stackexchange.com/questions/115766/should-i-use-uuid-as-well-as-id)
- [Why your software should use UUIDs](https://devforth.io/blog/why-your-software-should-use-uuids/)
- [Hypertext Transfer Protocol (HTTP) Field Name Registry](https://www.iana.org/assignments/http-fields/http-fields.xhtml)
- [learn-go-with-tests](https://quii.gitbook.io/learn-go-with-tests)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---
