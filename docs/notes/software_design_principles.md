<a name="readme-top"></a>

- [SOLID](#solid)
  - [1. **Single Responsibility Principle (SRP)**](#1-single-responsibility-principle-srp)
  - [2. **Open/Closed Principle (OCP)**](#2-openclosed-principle-ocp)
  - [3. **Liskov Substitution Principle (LSP)**](#3-liskov-substitution-principle-lsp)
  - [4. **Interface Segregation Principle (ISP)**](#4-interface-segregation-principle-isp)
  - [5. **Dependency Inversion Principle (DIP)**](#5-dependency-inversion-principle-dip)
  - [Benefits of Applying SOLID Principles](#benefits-of-applying-solid-principles)
- [References](#references)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# SOLID

SOLID is an acronym representing five fundamental principles of object-oriented programming and design. These principles help developers create more understandable, flexible, and maintainable software systems. Although SOLID principles were originally intended for object-oriented programming, they are broadly applicable across various programming paradigms, including functional programming and procedural programming. Here's a brief overview of each principle:

## 1. **Single Responsibility Principle (SRP)**

- **Definition**: A class should have only one reason to change, meaning it should have only one job or responsibility.
- **Purpose**: To ensure that a class is focused on a single task, making it easier to understand, test, and maintain.
- **Example**: In a user management system, the `User` class should only handle user-related data and operations, while logging or email notifications should be handled by separate classes.

## 2. **Open/Closed Principle (OCP)**

- **Definition**: Software entities (classes, modules, functions, etc.) should be open for extension but closed for modification.
- **Purpose**: To allow the behavior of a module to be extended without modifying its source code, thus reducing the risk of introducing bugs.
- **Example**: Instead of modifying a payment processing class to support new payment methods, you can extend it by adding new subclasses that implement a common interface.

## 3. **Liskov Substitution Principle (LSP)**

- **Definition**: Subtypes must be substitutable for their base types without altering the correctness of the program.
- **Purpose**: To ensure that derived classes extend base classes without changing their behavior, making the system more predictable and robust.
- **Example**: If you have a base class `Bird` with a method `fly()`, a subclass `Penguin` should not override `fly()` in a way that violates the expected behavior (since penguins can't fly).

## 4. **Interface Segregation Principle (ISP)**

- **Definition**: Clients should not be forced to depend on interfaces they do not use.
- **Purpose**: To create more specific and granular interfaces, reducing the implementation burden on classes and improving system modularity.
- **Example**: Instead of having a single `Animal` interface with methods `fly()`, `swim()`, and `walk()`, create smaller interfaces like `Flyable`, `Swimmable`, and `Walkable` to ensure classes implement only what they need.

## 5. **Dependency Inversion Principle (DIP)**

- **Definition**: High-level modules should not depend on low-level modules. Both should depend on abstractions (e.g., interfaces). Abstractions should not depend on details; details should depend on abstractions.
- **Purpose**: To reduce coupling between components, making the system more flexible and easier to modify.
- **Example**: In a messaging system, instead of a `Notification` class depending on a concrete `EmailSender` class, it should depend on an `INotificationSender` interface. Different implementations of `INotificationSender` (e.g., `EmailSender`, `SMSSender`) can be used interchangeably.

## Benefits of Applying SOLID Principles

1. **Maintainability**: Code is easier to understand and maintain, reducing the likelihood of bugs and making it easier to fix issues when they arise.
2. **Scalability**: Systems can be more easily extended with new features without modifying existing code, which minimizes the risk of introducing new bugs.
3. **Testability**: Modular and well-defined code components make unit testing more straightforward and effective.
4. **Reusability**: By adhering to SOLID principles, code components become more reusable across different parts of the application or even in different projects.
5. **Flexibility**: Systems become more adaptable to changing requirements and can accommodate new features with minimal changes to the existing codebase.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# References

[Software Design Patterns Tutorial](https://www.geeksforgeeks.org/software-design-patterns/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>
<p align="right">(<a href="../SOFTWARE_DEV.MD">back to main</a>)</p>

---
