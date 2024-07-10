<a name="readme-top"></a>

- [Setup](#setup)
- [JWT](#jwt)
  - [Implementation of Refresh Token with Redis](#implementation-of-refresh-token-with-redis)
  - [Redis and JWT Statelessness](#redis-and-jwt-statelessness)
  - [What are Access and Refresh Tokens?](#what-are-access-and-refresh-tokens)
  - [How Does It Work?](#how-does-it-work)
  - [Managing Security Threats: CSRF and XSS Attacks](#managing-security-threats-csrf-and-xss-attacks)
  - [Security Implementation Details](#security-implementation-details)
  - [Workflow Overview](#workflow-overview)
  - [Diagram: Token Workflow](#diagram-token-workflow)
- [Concepts And Methodologies](#concepts-and-methodologies)
  - [1. Architectural Styles](#1-architectural-styles)
  - [2. Architectural Patterns](#2-architectural-patterns)
  - [3. Design Patterns](#3-design-patterns)
  - [4. Programming Paradigms](#4-programming-paradigms)
  - [5. Methodologies](#5-methodologies)
  - [6. Development Practices](#6-development-practices)
  - [7. Software Development Frameworks](#7-software-development-frameworks)
  - [8. Software Design Principles](#8-software-design-principles)
  - [9. Software Architectural Principles](#9-software-architectural-principles)
  - [10. Software Quality Attributes](#10-software-quality-attributes)
- [References](#references)
  - [Golang](#golang)
  - [Refresh JWTs \& Redis](#refresh-jwts--redis)
  - [Others](#others)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Setup

| No. |                     Includes                     |                                             Description                                              |
| :-: | :----------------------------------------------: | :--------------------------------------------------------------------------------------------------: |
| 1.  | Easy setup using Env files, Docker, and Makefile |                       Simplifies initial setup and environment configuration.                        |
| 2.  | Ease of maintenance of Go packages with testing  |                      Ensures code reliability and simplifies testing processes.                      |
| 3.  |           Transitioning Away From MVC            |                 [🔗](./notes/architectural_patterns.md#transitioning-away-from-mvc)                  |
| 4.  |    Implementation of Refresh Token with Redis    |                                              [🔗](#jwt)                                              |
| 5.  |     Fiber as a web framework for performance     |                    Optimized for minimal memory allocation and high performance.                     |
| 6.  |                   Middlewares                    |                         [🔗](./notes/software_dev_framework.md#middlewares)                          |
| 7.  |   Internal logging for security using zerolog    | Provides secure and efficient logging mechanisms,</br>including structured logging with JSON output. |

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# JWT

## Implementation of Refresh Token with Redis

JWTs are stateless and eliminate the need for frequent database queries. Redis enhances token validation efficiency by reducing database queries and providing built-in expiration times for stored keys, making it ideal for managing tokens with limited lifespans. Redis handles token revocation effectively by invalidating tokens and preventing unauthorized access without additional database queries. Leveraging Redis achieves a seamless and secure authentication process.

## Redis and JWT Statelessness

Storing token details in Redis does not violate the statelessness principle of JWTs (JSON Web Tokens). Here's why:

1. **JWT Statelessness**: JWTs are considered stateless because they contain all the necessary information within the token itself, including claims about the user and any other metadata needed for verification. Statelessness refers to the server's ability to verify the token without needing to store session information or additional user data on the server side.

2. **Redis as a Cache**: Using Redis in this context is typically for performance optimization, not for maintaining the state of a session. Token details are stored in Redis to quickly access them for operations like token revocation or checking expiration, rather than maintaining user session state.

3. **Optional Use of Redis**: Redis is an optional layer for efficiency. JWTs can still be verified independently of Redis, as they contain all the necessary information. Redis enhances performance (due to its in-memory capabilities) and manages token lifecycles (like handling expiration more efficiently), but it does not fundamentally change the stateless nature of JWT authentication.

4. **Token Revocation**: When using Redis for storing tokens (or token blacklists), it is typically for the purpose of token revocation. This is an added security measure that doesn't conflict with JWT statelessness. The tokens remain self-contained, but there is an additional mechanism to invalidate them if necessary.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## What are Access and Refresh Tokens?

**Access Tokens** are short-term keys that grant users access to specific resources or data, typically valid for a short period (e.g., 15 minutes). When users log in, they receive an access token, which is included in their requests to access protected resources.

**Refresh Tokens** are long-term keys that allow users to request new access tokens without logging in again. These tokens have a longer lifespan (e.g., 60 minutes) and keep users logged in by providing new access tokens when the old ones expire.

## How Does It Work?

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

## Managing Security Threats: CSRF and XSS Attacks

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

## Security Implementation Details

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

## Workflow Overview

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

## Diagram: Token Workflow

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

# Concepts And Methodologies

## [1. Architectural Styles](./notes/architectural_styles.md#architectural-styles)

- **Layered Architecture**:
  - Separates functions into distinct layers, each with specific roles and responsibilities.
  - Layers communicate only with adjacent layers.
  - Common in enterprise applications.
- **Object-Oriented Architecture**:
  - Structures a system as a collection of interacting objects.
  - Each object represents an instance of a class.
  - Emphasizes reusability and modularity.
- **Data Centered Architecture**:
  - Centralized data repository accessed by multiple components.
  - Ensures consistency and availability of data.
  - Example: Database management systems.
- **Event-Based Architecture**:
  - Components interact through the production and consumption of events.
  - Decouples components for flexibility and scalability.
  - Common in real-time systems.

## [2. Architectural Patterns](./notes/architectural_patterns.md#architectural-patterns)

1. **Hexagonal Architecture**: Also known as Ports and Adapters, this architecture isolates the core logic from external dependencies.
2. [**Layered Architecture**: Divides the system into layers.](./notes/architectural_patterns.md#layered-architecture)
3. [**Microservices Architecture**: Structuring applications as a collection of loosely coupled services.](./notes/architectural_patterns.md#microservices-architecture)

## 3. Design Patterns

- **Creational Patterns**:
  - **Singleton**: Ensuring a type has only one instance.
  - **Factory**: Creating objects without specifying the exact type.
  - **Builder**: Constructing complex objects step by step.
- **Structural Patterns**:
  - **Adapter**: Allowing incompatible interfaces to work together.
  - **Decorator**: Adding behavior to objects dynamically.
  - **Facade**: Providing a simplified interface to a complex subsystem.
- **Behavioral Patterns**:
  - **Strategy**: Encapsulating algorithms within a family and making them interchangeable.
  - **Observer**: Notifying dependent objects about state changes.
  - **Command**: Encapsulating a request as an object.

## 4. Programming Paradigms

- **Concurrent Programming**: Using goroutines and channels for concurrency.
- **Procedural Programming**: Organizing code in procedures/functions.
- **Functional Programming**: Using higher-order functions and immutability.

## 5. Methodologies

- **Agile**: Iterative development with a focus on collaboration.
- **Scrum**: Framework for managing work with Sprints.
- **Kanban**: Visual workflow management.
- **Lean**: Reducing waste and improving efficiency.
- **Waterfall**: Sequential design process.

## [6. Development Practices](./notes/development_practices.md)

- **Test-Driven Development (TDD)**: Writing tests using the `testing` package before code.
- **Behavior-Driven Development (BDD)**: An extension of TDD that emphasizes collaboration between developers, testers, and business stakeholders. It involves writing high-level test scenarios in natural language using tools like `godog`, which are then automated to ensure the system behaves as expected. BDD focuses on the behavior of the system from the user's perspective and ensures all features add business value.
- **Domain-Driven Design (DDD)**: Structuring and modeling software around the business domain, emphasizing collaboration with domain experts and clear bounded contexts.
- **Continuous Integration (CI)**: Using tools like GitHub Actions, Travis CI, or CircleCI.
- **Continuous Deployment (CD)**: Automating deployment pipelines with GoCD or Jenkins.

## [7. Software Development Frameworks](./notes/software_dev_framework.md)

- **[Fiber: Express-inspired web framework for building REST APIs.](./notes/software_dev_framework.md)**
- **Gin**: Lightweight web framework.
- **Beego**: Full-featured web framework.
- **Echo**: High-performance, minimalist web framework.
- **Buffalo**: Rapid development web framework.
- **GORM**: ORM library for database interactions.

## [8. Software Design Principles](./notes/software_design_principles.md)

- [**SOLID**: The five fundamental principles of object-oriented programming and design help developers create understandable, flexible, and maintainable software systems.](./notes/software_design_principles.md#solid)
- **DRY (Don't Repeat Yourself)**: Avoiding code duplication.
- **KISS (Keep It Simple, Stupid)**: Prioritizing simplicity.
- **YAGNI (You Aren't Gonna Need It)**: Avoiding unnecessary features.

## 9. Software Architectural Principles

- **Separation of Concerns**: Modularizing code.
- **Single Responsibility Principle**: Each function/type should have one responsibility.
- **Open/Closed Principle**: Types should be open for extension but closed for modification.

## 10. Software Quality Attributes

- **Scalability**: Using concurrency patterns and distributed systems.
- **Reliability**: Ensuring robust error handling and resilience.
- **Maintainability**: Writing clean, readable code and documentation.
- **Usability**: Building intuitive APIs and user interfaces.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# References

## Golang

- [Standard Package Layout](https://www.gobeyond.dev/standard-package-layout/)
- [go-project-layout](https://appliedgo.com/blog/go-project-layout)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- [A Makefile/Dockerfile example for Go projects](https://github.com/thockin/go-build-template)
- [Ashley McNamara + Brian Ketelsen. Go best practices.](https://www.youtube.com/watch?v=MzTcsI6tn-0)
- [Go best practices, six years in](https://peter.bourgon.org/go-best-practices-2016/)
- [Common CRUD Design in Go](https://www.gobeyond.dev/crud/)
- [Difference Between Architectural Style, Architectural Patterns and Design Patterns](https://www.geeksforgeeks.org/difference-between-architectural-style-architectural-patterns-and-design-patterns/)

## Refresh JWTs & Redis

- [How to Properly Refresh JWTs for Authentication in Golang](https://codevoweb.com/how-to-properly-use-jwt-for-authentication-in-golang/)

## Others

- [Go postgres driver](https://github.com/lib/pq) recommends [pgx - PostgreSQL Driver and Toolkit](https://github.com/jackc/pgx)
- [Should I use UUID as well as ID](https://dba.stackexchange.com/questions/115766/should-i-use-uuid-as-well-as-id)
- [Why your software should use UUIDs](https://devforth.io/blog/why-your-software-should-use-uuids/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
