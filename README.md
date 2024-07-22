<a name="readme-top"></a>

- [Intro](#intro)
  - [Overview](#overview)
  - [Documentation](#documentation)
- [Setup](#setup)
  - [Handle Initial Files](#handle-initial-files)
  - [Generate the Private and Public Keys](#generate-the-private-and-public-keys)
    - [Shell Script](#shell-script)
    - [Browser Method](#browser-method)
- [Shell](#shell)
  - [directory](#directory)
  - [testing](#testing)
  - [psql](#psql)
  - [redis](#redis)

# Intro

This is a Golang-based web application leveraging PostgreSQL for data storage, Redis for refresh token management, and an architectural foundation rooted in Hexagonal Architecture (Hex Arch) and Domain-Driven Design (DDD).

## Overview

- **Golang and PostgreSQL**: The application uses Go for its backend services, with PostgreSQL as the primary database for data storage and retrieval. This combination ensures high performance and reliability.
- **Refresh Token Management with Redis**: User authentication and session management are handled securely using Redis for storing refresh tokens, ensuring quick access and improved security.
- **Hexagonal Architecture and Domain-Driven Design**: Hex Arch principles and DDD create a modular, maintainable codebase, facilitating clear boundaries between the core domain logic and peripheral components.
- **Fiber Web Framework and Zerolog**: The project tightly couples the Fiber web framework and Zerolog to leverage the high performance and minimalistic design of Fiber along with Zerolog's efficient structured logging. This integration streamlines development, enhances debugging and monitoring capabilities, and reduces boilerplate code, though it introduces some tradeoffs in flexibility and complexity.

## Documentation

For more detailed information on the setup, architecture, and various components of this project, refer to the [`Docs`](./docs/README.md) folder.

# Setup

## Handle Initial Files

1. **Respective `.env` files in `internal/infrastructure/config` folder**
2. **Respective env server `.json` file in `internal/infrastructure/db/postgres/json` folder:** Establishes server connection from pgAdmin to Postgres
3. **Run `chmod` command for the shell script(s)**

## Generate the Private and Public Keys

### Shell Script

```sh
# generate keys in base64
# alternatively use the browser method
chmod +x deployment/build/scripts/refresh_token_keygen.sh
cd build/scripts && ./refresh_token_keygen.sh && cd ../..

# Format app.log
chmod +x deployment/build/scripts/format_app_log.sh
```

### Browser Method

1. [Online RSA Key Generator](https://travistidwell.com/jsencrypt/demo/): Key Size: 2048 bit
2. [BASE64 Decode and Encode](https://www.base64encode.org/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Shell

## directory

```sh
brew install tree # MacOS
tree -d
```

## testing

```sh
# makefile defaults to dev env
make wat
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
