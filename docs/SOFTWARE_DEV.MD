<a name="readme-top"></a>

- [Setup](#setup)
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

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Setup

| No. |                     Includes                     |                                             Description                                              |
| :-: | :----------------------------------------------: | :--------------------------------------------------------------------------------------------------: |
| 1.  | Easy setup using Env files, Docker, and Makefile |                       Simplifies initial setup and environment configuration.                        |
| 2.  | Ease of maintenance of Go packages with testing  |                      Ensures code reliability and simplifies testing processes.                      |
| 3.  |           Model-View-Controller (MVC)            |                  [🔗](./notes/architectural_patterns.md#model-view-controller-mvc)                   |
| 4.  |          Layered (N-tier) Architecture           |                 [🔗](./notes/architectural_patterns.md#layered-n-tier-architecture)                  |
| 5.  |    Implementation of Refresh Token with Redis    |          [🔗](./notes/software_dev_framework.md#implementation-of-refresh-token-with-redis)          |
| 6.  |                   Middlewares                    |                         [🔗](./notes/software_dev_framework.md#middlewares)                          |
| 7.  |     Fiber as a web framework for performance     |                    Optimized for minimal memory allocation and high performance.                     |
| 8.  |   Internal logging for security using zerolog    | Provides secure and efficient logging mechanisms,</br>including structured logging with JSON output. |
| 9.  |   A common folder structure for project layout   |                      Facilitates scalability into a microservice architecture.                       |

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

1. [**MVC (Model-View-Controller)**: Used in web frameworks like Fiber.](./notes/architectural_patterns.md#model-view-controller-mvc)
2. [**Layered (n-tier)**: Divides the system into layers.](./notes/architectural_patterns.md#layered-n-tier-architecture)
3. **Client-Server Architecture**: Divides the system into clients that request services and servers that provide them.
4. **Event-Driven Architecture**: Components communicate through events, enabling asynchronous communication.
5. **Microkernel Architecture**: Building extensible systems with core functionality and plug-ins.
6. [**Microservices Architecture**: Structuring applications as a collection of loosely coupled services.](./notes/architectural_patterns.md#combining-microservices-and-n-tier-architectures)
7. **Space-Based Architecture**: Designed to address scalability and concurrency issues by distributing processing and storage across multiple nodes.
8. **Master-Slave Architecture**: Master component distributes tasks to slave components.
9. **Pipe-Filter Architecture**: Data processing elements are arranged in a linear sequence.
10. **Broker Architecture**: Components interact through a broker, often used in distributed systems.
11. **Peer-to-Peer Architecture**: Components (peers) interact directly without a central server.
12. **Service-Oriented Architecture (SOA)**: Designing the system as a collection of services.

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
- **Pair Programming**: Collaborative coding to improve quality and share knowledge.

## [7. Software Development Frameworks](./notes/software_dev_framework.md)

- **Gin**: Lightweight web framework.
- **Beego**: Full-featured web framework.
- **Echo**: High-performance, minimalist web framework.
- **Buffalo**: Rapid development web framework.
- **[Fiber: Express-inspired web framework for building REST APIs.](./notes/software_dev_framework.md)**
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

- [Difference Between Architectural Style, Architectural Patterns and Design Patterns](https://www.geeksforgeeks.org/difference-between-architectural-style-architectural-patterns-and-design-patterns/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
