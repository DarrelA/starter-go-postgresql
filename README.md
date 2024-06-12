- [N-tier Architecture](#n-tier-architecture)
- [Model-View-Controller (MVC)](#model-view-controller-mvc)
- [Setup](#setup)
  - [Handle Initial Files](#handle-initial-files)
  - [Generate the Private and Public Keys](#generate-the-private-and-public-keys)
- [Maintenance](#maintenance)
  - [`go.mod` file](#gomod-file)
  - [psql](#psql)
  - [golang-migrate](#golang-migrate)
- [References](#references)

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

# Model-View-Controller (MVC)

MVC is a design pattern that separates an application into three interconnected components:

1. **Model:** Represents the data and business logic. It directly manages the data and rules of the application.
2. **View:** Presents data to the user. It represents the UI components and is used to display the model's data.
3. **Controller:** Acts as an intermediary between the Model and View. It listens to the user input (through the View) and processes the requests (updating the Model, selecting the View to display).

# Setup

## Handle Initial Files

- Respective `.env` files in `configs` folder
- **`init.sql` script**: Create initial schemas
- **`servers.json` file**: Establish server connection from pgAdmin to Postgres

## Generate the Private and Public Keys

- [Online RSA Key Generator](https://travistidwell.com/jsencrypt/demo/)
  - Key Size: 2048 bit
- [BASE64 Decode and Encode](https://www.base64encode.org/)

# Maintenance

## `go.mod` file

```sh
# Updating `go.mod`
go get -u
go mod tidy

# pgAdmin
http://localhost:5050/browser/
```

## psql

```sh
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

# References

- [go-project-layout](https://appliedgo.com/blog/go-project-layout)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout/tree/master)
- [wpcodevo/golang-fiber-jwt-rs256](https://github.com/wpcodevo/golang-fiber-jwt-rs256)
