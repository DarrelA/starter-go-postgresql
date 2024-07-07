<a name="readme-top"></a>

- [Architectural Patterns](#architectural-patterns)
  - [Model-View-Controller (MVC)](#model-view-controller-mvc)
  - [Layered (n-tier) Architecture](#layered-n-tier-architecture)
  - [Client-Server Architecture](#client-server-architecture)
  - [Event-Driven Architecture](#event-driven-architecture)
  - [Microkernel Architecture](#microkernel-architecture)
  - [Microservices Architecture](#microservices-architecture)
  - [Space-Based Architecture](#space-based-architecture)
  - [Master-Slave Architecture](#master-slave-architecture)
  - [Pipe-Filter Architecture](#pipe-filter-architecture)
  - [Broker Architecture](#broker-architecture)
  - [Peer-to-Peer Architecture](#peer-to-peer-architecture)
  - [Service-Oriented Architecture (SOA)](#service-oriented-architecture-soa)
- [Combining Microservices and N-tier Architectures](#combining-microservices-and-n-tier-architectures)
- [References](#references)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Architectural Patterns

**Architectural Patterns** are specific, reusable solutions to common problems within a given context in software architecture.

- **Focus**: Concrete solutions, reusable designs, and best practices.
- **Purpose**: Provide templates for solving specific problems within a system's architecture.
- **Scope**: Often focus on a specific aspect or layer of the system.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

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

## Layered (n-tier) Architecture

An Layered (n-tier) architecture separates each layer or tier of the application both physically and logically, which enhances flexibility and manageability.

1. **Frontend (Presentation Layer):** This is the user interface of the application. In an Layered (n-tier) setup, you can update or replace the frontend without needing to alter the backend, as long as the interface contract (e.g., APIs) remains consistent.

2. **Backend (Business Logic Layer):** This layer contains the business rules and logic. You can modify or replace this layer as needed without affecting the other tiers, provided the interfaces between the layers don't change.

3. **Database (Data Access Layer):** This layer handles data storage and retrieval. You can switch out databases or change how data is accessed, and as long as you maintain the same data contracts or APIs, the other layers remain unaffected.

Benefits of Layered (n-tier) architecture include:

1. **Scalability:** Each layer can be scaled independently based on its resource requirements.
2. **Flexibility:** You can replace or upgrade one layer without significant rework of the others.
3. **Maintainability:** Smaller, well-defined codebases for each layer are easier to manage and understand.

However, N-tier architecture can also introduce complexity in terms of network latency, configuration, and management. Each layer likely needs to communicate over a network, which can affect performance and add complexity that wouldn't be as prominent in a monolithic architecture.

## Client-Server Architecture

Divides the system into clients that request services and servers that provide them. Clients initiate communication while servers await client requests.

## Event-Driven Architecture

Components communicate through events, enabling asynchronous communication and decoupling of components.

## Microkernel Architecture

Building extensible systems with core functionality and plug-ins that add features or services, promoting modularity and flexibility.

## Microservices Architecture

Structuring applications as a collection of loosely coupled, independently deployable services, often organized around business capabilities, enhancing scalability and maintainability.

## Space-Based Architecture

Designed to address scalability and concurrency issues by distributing processing and storage across multiple nodes, reducing the risk of bottlenecks and improving system resilience.

## Master-Slave Architecture

Master component distributes tasks to slave components, often used in scenarios requiring parallel processing.

## Pipe-Filter Architecture

Data processing elements are arranged in a linear sequence, where each element (filter) processes data and passes it to the next.

## Broker Architecture

Components interact through a broker, often used in distributed systems to manage communication and data exchange.

## Peer-to-Peer Architecture

Components (peers) interact directly without a central server, sharing resources and data among themselves.

## Service-Oriented Architecture (SOA)

Designing the system as a collection of services, each providing a business function and communicating over a network, promoting reusability and integration.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Combining Microservices and N-tier Architectures

Combining N-tier architecture with microservices involves separating concerns within an application while decoupling it into independently deployable services. Each microservice can adhere to N-tier principles, enhancing modularity and scalability.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# References

- [Types of Software Architecture Patterns](https://www.geeksforgeeks.org/types-of-software-architecture-patterns/)
- [14 software architecture design patterns to know](https://www.redhat.com/architect/14-software-architecture-patterns)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
<p align="right">(<a href="../SOFTWARE_DEV.MD">back to main</a>)</p>

---
