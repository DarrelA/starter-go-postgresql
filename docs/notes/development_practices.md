<a name="readme-top"></a>

- [Testing](#testing)
  - [Unit Testing, Integration Testing, and Acceptance Testing](#unit-testing-integration-testing-and-acceptance-testing)
    - [Black Box Testing \& White Box Testing](#black-box-testing--white-box-testing)
- [Dependency Injection (DI)](#dependency-injection-di)
  - [Methods of DI](#methods-of-di)
- [Error Handling](#error-handling)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Testing

## Unit Testing, Integration Testing, and Acceptance Testing

| Aspect              | Unit Testing                        | Integration Testing                                | Acceptance Testing                          |
| ------------------- | ----------------------------------- | -------------------------------------------------- | ------------------------------------------- |
| **Purpose**         | Verify individual units/components  | Verify interaction between components              | Verify the entire system meets requirements |
| **Scope**           | Small, isolated pieces of code      | Multiple components or subsystems                  | Entire application or major features        |
| **Nature**          | Typically automated, quick to run   | Typically automated, moderate to run               | Can be automated or manual, longer to run   |
| **Examples**        | Testing a single function or method | Testing database interaction, API calls            | Testing end-to-end user scenarios           |
| **Tools**           | Go `testing` package                | Go `testing` package                               | `godog` package                             |
| **Characteristics** | White-box testing, high frequency   | Mix of white-box and black-box, moderate frequency | Black-box testing, less frequent            |

[User Stories] become [Acceptance Tests] which is [Behavior Driven Development] "Doing the RIGHT thing."

[Code Functionality] becomes [Unit Testing] which is [Test Driven Development] "Doing the THING right."

### Black Box Testing & White Box Testing

| Aspect                 | Black Box Testing                                     | White Box Testing                                           |
| ---------------------- | ----------------------------------------------------- | ----------------------------------------------------------- |
| **Knowledge Required** | None about internal code                              | In-depth knowledge of internal code                         |
| **Focus**              | Functional behavior                                   | Internal logic and structure                                |
| **Test Basis**         | Requirements, specifications, use cases               | Source code, architecture, internal documentation           |
| **Techniques**         | Equivalence partitioning, boundary value analysis     | Path testing, loop testing, code coverage analysis          |
| **Advantages**         | User perspective, identifies discrepancies            | Thorough coverage, optimizes code, identifies hidden errors |
| **Disadvantages**      | Limited coverage, difficult to identify all scenarios | Time-consuming, requires in-depth knowledge, expensive      |

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Dependency Injection (DI)

Dependency Injection (DI) is a design pattern that achieves Inversion of Control (IoC) by allowing components to receive their dependencies from an external source rather than creating them internally. In Go, DI is typically implemented by passing dependencies as parameters to functions or struct constructors. This approach decouples application components, making them more modular, easier to test, and maintain. DI allows replacing real implementations of dependencies with mock versions during testing, ensuring components can be tested in isolation.

- **Promotes Loose Coupling**: DI promotes loose coupling and modularity.
- **Defines Clear Contracts**: Interfaces define clear contracts for components.
- **Encapsulates Related Data**: Structs encapsulate related data and behaviors.
- **Enables Effective Mocking**: Together, these concepts enable effective mocking and testing in Go.

## Methods of DI

1. **Constructor Injection**: Best for mandatory dependencies that do not change after object creation.
2. **Setter Injection**: Useful for optional dependencies or when dependencies need to be changed after object creation.
3. **Interface Injection**: Suitable for cases where you want the dependency itself to control the injection.
4. **Field Injection**: Handy in frameworks that automate dependency injection and where you want minimal boilerplate code.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Error Handling

- [Appendix A. PostgreSQL Error Codes](https://www.postgresql.org/docs/current/errcodes-appendix.html)
- [What is Test Driven Development (TDD)?](https://www.geeksforgeeks.org/test-driven-development-tdd/)
- [What is Behavior-Driven Development (BDD)?](https://www.geeksforgeeks.org/behavioral-driven-development-bdd-in-software-engineering/)
- [Domain-Driven Design (DDD)](https://www.geeksforgeeks.org/domain-driven-design-ddd/)
- [Dependency Injection(DI) Design Pattern](https://www.geeksforgeeks.org/dependency-injectiondi-design-pattern/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
<p align="right">(<a href="../SOFTWARE_DEV.MD">back to main</a>)</p>

---
