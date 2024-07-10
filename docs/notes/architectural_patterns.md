<a name="readme-top"></a>

- [Architectural Patterns](#architectural-patterns)
- [Transitioning Away From MVC](#transitioning-away-from-mvc)
  - [Model-View-Controller (MVC)](#model-view-controller-mvc)
  - [Circular Dependency in MVC](#circular-dependency-in-mvc)
- [Clean Architecture](#clean-architecture)
  - [Key principles of Clean Architecture](#key-principles-of-clean-architecture)
  - [The structure typically involves](#the-structure-typically-involves)
- [Hexagonal Architecture](#hexagonal-architecture)
  - [Key concepts of Hexagonal Architecture](#key-concepts-of-hexagonal-architecture)
  - [Benefits of Hexagonal Architecture](#benefits-of-hexagonal-architecture)
- [CQRS (Command Query Responsibility Segregation)](#cqrs-command-query-responsibility-segregation)
  - [Key concepts of CQRS](#key-concepts-of-cqrs)
  - [Benefits of CQRS](#benefits-of-cqrs)
  - [Polyglot Persistence](#polyglot-persistence)
  - [Example Scenario](#example-scenario)
  - [Event-Sourcing](#event-sourcing)
    - [Benefits of Event Sourcing in CQRS](#benefits-of-event-sourcing-in-cqrs)
    - [Key Benefits of Event Sourcing for Audit Requirements](#key-benefits-of-event-sourcing-for-audit-requirements)
    - [How Event Sourcing Meets Audit Requirements in CQRS](#how-event-sourcing-meets-audit-requirements-in-cqrs)
    - [Example Scenario: Auditing a Financial System](#example-scenario-auditing-a-financial-system)
    - [Challenges and Considerations](#challenges-and-considerations)
    - [Handling a REST API POST request that returns the created entity](#handling-a-rest-api-post-request-that-returns-the-created-entity)
- [Event-Driven Architecture](#event-driven-architecture)
  - [Key components of EDA include](#key-components-of-eda-include)
  - [Advantages of EDA](#advantages-of-eda)
- [Layered Architecture](#layered-architecture)
- [Microservices Architecture](#microservices-architecture)
  - [Key characteristics of Microservices Architecture](#key-characteristics-of-microservices-architecture)
  - [Benefits of Microservices Architecture](#benefits-of-microservices-architecture)
- [References](#references)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Architectural Patterns

**Architectural Patterns** are specific, reusable solutions to common problems within a given context in software architecture.

- **Focus**: Concrete solutions, reusable designs, and best practices.
- **Purpose**: Provide templates for solving specific problems within a system's architecture.
- **Scope**: Often focus on a specific aspect or layer of the system.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Transitioning Away From MVC

## Model-View-Controller (MVC)

The Model-View-Controller (MVC) design pattern separates an application into three interconnected components, improving the organization and manageability of the code.

- **Model:** Manages the data and business logic of the application.
- **View:** Represents the user interface and displays the data.
- **Controller:** Acts as an intermediary between the Model and View, processing user input and updating the Model and View accordingly.

## Circular Dependency in MVC

- **Coupling**: The components (Model, View, Controller) are tightly coupled. The Controller often directly interacts with both the Model and the View, and sometimes, the View may need to update the Model directly, leading to potential circular dependencies.
- **Dependencies**: Because each layer often has knowledge of and interacts directly with the other layers, changes in one layer can necessitate changes in another, leading to a higher risk of creating circular dependencies.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

![interest_on_technical_debt](/assets/interest_on_technical_debt.svg)

# Clean Architecture

Clean Architecture is a software design philosophy introduced by Robert C. Martin (Uncle Bob). It aims to create a structure that is easy to maintain, scalable, and resilient to change. The primary goals are separation of concerns and independence of frameworks, user interfaces, databases, and any external agencies.

## Key principles of Clean Architecture

- **Independence of Frameworks**: The architecture does not depend on the existence of some library of feature-laden software. This allows you to use such frameworks as tools rather than requiring the software to fit within the framework.
- **Testability**: Business rules can be tested without the user interface, database, web server, or any other external element.
- **Independence of UI**: The UI can change easily, without changing the rest of the system.
- **Independence of Database**: You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else without changing the business rules.
- **Independence of any external agency**: In fact, your business rules simply don’t know anything at all about the outside world.

## The structure typically involves

- **Entities**: Business objects that have business rules.
- **Use Cases**: Application-specific business rules.
- **Interface Adapters**: Converters that interface the use cases and entities to external parts of the system.
- **Frameworks and Drivers**: The external agents like the web, UI, database, and frameworks.

![clean_architecture](/assets/clean_architecture.jpeg)

# Hexagonal Architecture

Hexagonal Architecture, also known as Ports and Adapters, is a design pattern introduced by Alistair Cockburn. It aims to create loosely coupled application components that can be easily connected to their software environment via ports and adapters. This architecture allows for a clear separation between the core logic and the external dependencies.

## Key concepts of Hexagonal Architecture

- **Domain Layer**:
  - Contains the core domain logic, which includes **business rules and entities**.
  - Focuses on the core business functionality and maintains the state of the business model.
- **Application Layer**:
  - Contains application-specific logic that orchestrates the domain logic.
  - Defines use cases and services that leverage the core domain logic.
  - **Ports** are part of this layer, providing an abstraction over the core domain logic.
- **Infrastructure Layer**:
  - Contains **adapters** that interact with external systems such as databases, web services, and user interfaces.
  - Implements the ports defined in the application layer to facilitate communication with the external world.
  - Provides technical capabilities required by the application layer.

## Benefits of Hexagonal Architecture

- **Separation of Concerns**: The core logic is isolated from external concerns, making it easier to develop and maintain.
- **Testability**: The core can be tested independently of external systems, improving test coverage and reliability.
- **Flexibility**: Adapters can be swapped or modified without affecting the core logic, allowing for easy integration with different technologies.

![https://docs.google.com/drawings/d/1E_hx5B4czRVFVhGJbrbPDlb_JFxJC8fYB86OMzZuAhg/edit](/assets/hexagonal_architecture_shodiq_muhammad.svg)

# CQRS (Command Query Responsibility Segregation)

CQRS is a pattern that separates read and write operations into different models. This pattern is particularly useful in systems where the read and write workloads are different and can be optimized separately.

## Key concepts of CQRS

- **Commands**: Operations that change the state of the system (writes). They are responsible for updating the data.
- **Queries**: Operations that do not change the state of the system (reads). They are responsible for fetching data.
- **Command Model**: The part of the system that handles commands and updates the state.
- **Query Model**: The part of the system that handles queries and returns data.

## Benefits of CQRS

- **Separation of Concerns**: The separation of read and write concerns can lead to more scalable and maintainable systems.
- **Optimized Models**: Each model can be optimized for its specific use case (e.g., write models can be optimized for transactional integrity, while read models can be optimized for query performance).
- **Flexibility**: Easier to evolve and maintain as changes to the read model do not affect the write model and vice versa.

## Polyglot Persistence

**Polyglot Persistence** refers to using different data storage technologies to handle different data storage needs within the same application. Instead of using a single type of database for all storage needs, different databases are used for different parts of the application based on their unique requirements.

CQRS with Polyglot Persistence combines these two concepts. In this approach, the read and write sides of the CQRS pattern can use different types of databases that are best suited for their specific needs.

## Example Scenario

- **Command Side (Write Model)**
  - **Database**: A relational database (e.g., PostgreSQL) might be used for the write model to ensure ACID (Atomicity, Consistency, Isolation, Durability) properties.
  - **Structure**: The data model is normalized to reduce redundancy and maintain data integrity.
- **Query Side (Read Model)**
  - **Database**: A NoSQL database (e.g., MongoDB, Elasticsearch) might be used for the read model to provide fast read performance and flexible querying capabilities.
  - **Structure**: The data model is denormalized to optimize read operations, making it easier and faster to retrieve the necessary information.

## Event-Sourcing

**Event Sourcing is a design pattern** in which state changes of an application are stored as a sequence of events. Rather than storing only the current state, event sourcing records each change to the state as an event, which can then be replayed to determine the current state.

### Benefits of Event Sourcing in CQRS

1. **Auditability**: Since all state changes are stored as events, it is easy to audit and understand how the system reached its current state. This can be particularly useful for debugging and compliance purposes.
2. **Temporal Queries**: Event sourcing allows querying the state of the system at any point in time by replaying the events up to that point. This can be valuable for historical analysis and reporting.
3. **Replayability**: The ability to replay events enables reconstruction of the state and facilitates debugging, as you can reproduce any past state of the system.
4. **Scalability and Performance**: Write operations are typically append-only, which can be more performant than updating records in a traditional database.
5. **Integration with CQRS**: In a CQRS architecture, event sourcing naturally **fits with the command side**. Commands result in events, which are stored in the event store. The read side (queries) can subscribe to these events and update the read model accordingly.

### Key Benefits of Event Sourcing for Audit Requirements

- **Complete History of Changes**
  - **Event Logging**: Every change to the system’s state is logged as an immutable event. This ensures a complete and accurate history of all actions taken on the data.
  - **Traceability**: You can trace back every state change to understand how and why the current state was reached.
- **Immutability and Tamper-Proof Records**
  - **Immutable Events**: Once an event is recorded, it cannot be altered. This immutability ensures that the audit trail is tamper-proof.
  - **Audit Integrity**: The integrity of the audit logs is maintained, as past events are preserved exactly as they occurred.
- **Temporal Queries**
  - **State at Any Point in Time**: Event Sourcing allows you to reconstruct the state of the system at any point in time by replaying events up to that point. This is invaluable for audits that require understanding the state of the system at specific historical moments.
  - **Time-Based Analysis**: Auditors can analyze how data evolved over time, providing insights into patterns and decision-making processes.
- **Compliance**
  - **Regulatory Requirements**: Many regulations require detailed logs of data changes (e.g., financial transactions, patient records). Event Sourcing inherently meets these requirements by capturing all changes as events.
  - **Non-Repudiation**: With a clear and immutable record of events, it’s easier to demonstrate compliance with legal and regulatory standards.

### How Event Sourcing Meets Audit Requirements in CQRS

- **Handling Commands**:
  - When a command (e.g., “UpdateCustomerAddress”) is processed, the corresponding aggregate (e.g., Customer) generates an event (e.g., “CustomerAddressUpdated”).
  - This event is stored in the event store, providing a record of the action taken.
- **Storing Events**:
  - Events are persisted in the event store in an ordered sequence. Each event includes metadata such as the timestamp, the user who performed the action, and any relevant contextual information.
  - This structured storage makes it straightforward to query and retrieve events for audit purposes.
- **Querying for Audit**:
  - Auditors can query the event store to retrieve all events related to a particular entity or a specific type of action.
  - They can filter events based on criteria such as time range, user, or event type, facilitating comprehensive audits.
- **Reconstructing State**:
  - To understand the state of an entity at a specific point in time, auditors can replay events up to that point.
  - This ability to reconstruct past states is crucial for audits that need to verify historical data or investigate specific incidents.

### Example Scenario: Auditing a Financial System

- **Command Handling**:
  - A command “TransferFunds” is issued to transfer money between accounts.
  - The aggregate responsible for the accounts processes this command and generates events like “FundsWithdrawn” and “FundsDeposited”.
- **Event Storage**:
  - These events are stored in the event store with metadata such as the amount transferred, source and destination accounts, and timestamps.
- **Audit Queries**:
  - Auditors can query the event store to retrieve all “TransferFunds” events within a specific period.
  - They can trace each transfer by looking at the “FundsWithdrawn” and “FundsDeposited” events to ensure that every transfer is accounted for and properly recorded.
- **State Reconstruction**:
  - If there is a dispute or an investigation, auditors can reconstruct the state of the accounts at the time of the transfer to verify the balance and transaction history.

### Challenges and Considerations

- **Event Schema Evolution**:
  - As the system evolves, the structure of events may change. Proper versioning and handling of different event schemas are necessary to maintain auditability.
- **Storage and Performance**:
  - Storing a large volume of events can become challenging. Efficient storage solutions and archiving strategies may be required to ensure performance and accessibility.
- **Privacy and Security**:
  - Sensitive information in events needs to be handled with care. Proper encryption and access control mechanisms are essential to protect the integrity and confidentiality of the audit logs.

### Handling a REST API POST request that returns the created entity

- **Command Handling**:
  - When a `POST` request is made to create a new entity, the command is processed by the write side of the system.
  - The command results in one or more events being generated and stored in the event store.
- **HTTP Response**:
  - Instead of returning the created entity immediately, the system responds with a `204 No Content` status code.
  - The `Content-Location` header is set to the URL of the newly created resource. This URL can be used by the client to fetch the created entity later.
- **Querying for the Entity**:
  - The client can then make a separate `GET` request to the URL provided in the `Content-Location` header to retrieve the created entity.
  - This `GET` request is handled by the read side of the system, which queries the read model to return the current state of the entity.

![cqrs_polyglot_persistence](/assets/cqrs_polyglot_persistence.webp)

# Event-Driven Architecture

Event-Driven Architecture (EDA) is a software architecture pattern that promotes the production, detection, consumption, and reaction to events. An event is a significant change in state. EDA is especially useful in distributed systems to decouple components and improve scalability and flexibility.

## Key components of EDA include

- **Event Producers**: Components that detect or sense events and then publish them.
- **Event Consumers**: Components that listen for events and process them.
- **Event Channels**: Pathways through which events travel between producers and consumers.
- **Event Processors**: Specialized processors that might filter, transform, or route events to appropriate consumers.

## Advantages of EDA

- **Loose Coupling**: Components are decoupled because they only communicate via events.
- **Scalability**: Easier to scale because components can be scaled independently.
- **Real-time Processing**: Supports real-time processing and responsiveness.

# Layered Architecture

An layered architecture separates each layer or tier of the application both physically and logically, which enhances flexibility and manageability.

- **Frontend (Presentation Layer):** This is the user interface of the application. In an layered architecture setup, you can update or replace the frontend without needing to alter the backend, as long as the interface contract (e.g., APIs) remains consistent.
- **Backend (Business Logic Layer):** This layer contains the business rules and logic. You can modify or replace this layer as needed without affecting the other tiers, provided the interfaces between the layers don't change.
- **Database (Data Access Layer):** This layer handles data storage and retrieval. You can switch out databases or change how data is accessed, and as long as you maintain the same data contracts or APIs, the other layers remain unaffected.

Benefits of layered architecture include:

- **Scalability:** Each layer can be scaled independently based on its resource requirements.- **Flexibility:** You can replace or upgrade one layer without significant rework of the others.- **Maintainability:** Smaller, well-defined codebases for each layer are easier to manage and understand.

However, layered architecture can also introduce complexity in terms of network latency, configuration, and management. Each layer likely needs to communicate over a network, which can affect performance and add complexity that wouldn't be as prominent in a monolithic architecture.

# Microservices Architecture

Microservices Architecture is an architectural style that structures an application as a collection of loosely coupled services. Each service is fine-grained and implements a single business capability. Microservices architecture enables continuous delivery and deployment of large, complex applications.

## Key characteristics of Microservices Architecture

- **Decentralization**: Each microservice is developed, deployed, and scaled independently. Teams can work on different services simultaneously without affecting each other.
- **Service Independence**: Each microservice can be written in different programming languages and use different data storage technologies.
- **Single Responsibility**: Each microservice focuses on a single business function, following the principle of single responsibility.
- **Inter-Service Communication**: Microservices communicate with each other via lightweight protocols such as HTTP/REST, messaging queues, or event streams.

## Benefits of Microservices Architecture

- **Scalability**: Services can be scaled independently, allowing for better resource utilization and improved performance.
- **Resilience**: Failure in one service does not necessarily impact the entire system. Services can be designed to handle failures gracefully.
- **Flexibility**: Teams can choose the best tools and technologies for each service, promoting innovation and improving productivity.
- **Continuous Deployment**: Microservices can be deployed independently, enabling faster and more frequent releases.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# References

- [Moving Towards Domain Driven Design in Go](https://www.calhoun.io/moving-towards-domain-driven-design-in-go/)
- [GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)
- [Golang UK Conference 2016 - Marcus Olsson - Building an enterprise service in Go](https://www.youtube.com/watch?v=twcDf_Y2gXY)
- [Go and a Package Focused Design](https://blog.gopheracademy.com/advent-2016/go-and-package-focused-design/)
- [Hexagonal Architecture](https://fideloper.com/hexagonal-architecture)
- [Hexagonal Architecture and Domain-Driven Design (DDD)](https://prbpedro.substack.com/p/hexagonal-architecture-and-domain)
- [Modern Business Software in Go](https://threedots.tech/series/modern-business-software-in-go)
- [WTF Dial is an example web application written in Go](https://github.com/benbjohnson/wtf)
- [go-kit/kit - A standard library for microservices](https://github.com/go-kit/kit)
- [Design Patterns Architecture](https://www.geeksforgeeks.org/design-patterns-architecture/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
<p align="right">(<a href="../SOFTWARE_DEV.MD">back to main</a>)</p>

---
