- [References](#references)
- [N-tier Architecture](#n-tier-architecture)
- [MVC \& DAO Design Pattern](#mvc--dao-design-pattern)
  - [Model-View-Controller (MVC)](#model-view-controller-mvc)
  - [Data Access Object (DAO)](#data-access-object-dao)
  - [Combining MVC and DAO](#combining-mvc-and-dao)
  - [Workflow](#workflow)
  - [Advantages of Combining MVC and DAO](#advantages-of-combining-mvc-and-dao)
  - [Considerations](#considerations)
- [Maintenance](#maintenance)

# References

- [go-project-layout](https://appliedgo.com/blog/go-project-layout)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout/tree/master)

# N-tier Architecture

In an n-tier architecture, each layer or tier of the application is physically and logically separate. This separation allows for greater flexibility and manageability.

1. **Frontend (Presentation Layer):** This is the user interface of the application. In an n-tier setup, you can update or completely replace the frontend without needing to alter the backend, as long as the interface contract (like APIs) between them remains consistent.

2. **Backend (Business Logic Layer):** This layer contains the business rules and logic. You can modify or replace this layer as needed without affecting the other tiers, again, as long as the interfaces between the layers don't change.

3. **Database (Data Access Layer):** This layer handles data storage and retrieval. You can switch out databases or change how data is accessed, and as long as you maintain the same data contracts or APIs, the other layers remain unaffected.

The n-tier architecture's decoupling allows for:

- **Scalability:** You can scale each layer independently based on its resource requirements.
- **Flexibility:** You can replace or upgrade one layer without significant rework of the others.
- **Maintainability:** Smaller, well-defined codebases for each layer are easier to manage and understand.

However, it's worth noting that while n-tier architecture offers these benefits, it also introduces complexity in terms of network latency, configuration, and management. Each layer will likely need to communicate over a network, which can introduce performance and complexity considerations that wouldn't be as prominent in a monolithic architecture.

# MVC & DAO Design Pattern

The combination of Model-View-Controller (MVC) and Data Access Object (DAO) design patterns is a powerful architectural approach in software development, particularly useful in web applications and services. Here's an overview of how these two design patterns can work together:

## Model-View-Controller (MVC)

MVC is a design pattern that separates an application into three interconnected components:

1. **Model:** Represents the data and business logic. It directly manages the data and rules of the application.
2. **View:** Presents data to the user. It represents the UI components and is used to display the model's data.
3. **Controller:** Acts as an intermediary between the Model and View. It listens to the user input (through the View) and processes the requests (updating the Model, selecting the View to display).

## Data Access Object (DAO)

DAO is a pattern used to abstract and encapsulate all access to the data source. It manages the connection with the data source to obtain and store data. The DAO implements the access mechanism required to work with the data source.

## Combining MVC and DAO

In a combined MVC and DAO setup, the patterns interact primarily through the Model component of MVC:

- **Model with DAO:** The Model component of MVC uses DAO to interact with the data source. DAO manages the persistence and retrieval of the Model data, hiding the details of the data source and retrieval mechanisms.

## Workflow

1. **User Interaction:** The user interacts with the View, which sends the user's input to the Controller.
2. **Controller Processing:** The Controller interprets the user input (sent from the View), commanding the Model to change its state (e.g., updating data).
3. **Model and DAO Interaction:** The Model uses DAO to retrieve or update data in the data source. DAO provides a clean API to the Model so that the Model does not need to know about the underlying database specifics.
4. **Update View:** Once the Model has changed state (data updated or retrieved), the View is updated to reflect the new data. The View gets the latest data from the Model.

## Advantages of Combining MVC and DAO

- **Separation of Concerns:** MVC separates the application logic, UI, and data access. DAO further isolates data access specifics from the business logic.
- **Reusability and Maintainability:** Both patterns promote reusability and maintainability. DAOs can be reused across different parts of the application. MVC facilitates changing the UI without affecting the underlying business logic.
- **Scalability and Flexibility:** Separating concerns makes it easier to scale and modify applications. You can change the data source or business logic without affecting other parts of the application.

## Considerations

- **Complexity:** Introducing multiple layers and abstractions can increase complexity. It's essential to balance the benefits with the complexity added.
- **Performance Overhead:** Each layer adds a bit of overhead. In performance-critical applications, the design should be carefully evaluated.

# Maintenance

```sh
# Updating `go.mod`
go get -u
go mod tidy
```
