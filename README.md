<a name="readme-top"></a>

- [Setup](#setup)
  - [Handle Initial Files](#handle-initial-files)
  - [Generate the Private and Public Keys](#generate-the-private-and-public-keys)
    - [Shell Script](#shell-script)
    - [Browser Method](#browser-method)
- [Shell](#shell)
  - [testing](#testing)
  - [psql](#psql)
  - [redis](#redis)

# Setup

## Handle Initial Files

1. **Respective `.env` files in `configs` folder**
2. **Respective env server `.json` file:** Establishes server connection from pgAdmin to Postgres
3. **Run `chmod` command for the shell script(s)**

## Generate the Private and Public Keys

### Shell Script

```sh
# generate keys in base64
# alternatively use the browser method
chmod +x build/scripts/refresh_token_keygen.sh
cd build && ./refresh_token_keygen.sh && cd ..
```

### Browser Method

1. [Online RSA Key Generator](https://travistidwell.com/jsencrypt/demo/): Key Size: 2048 bit
2. [BASE64 Decode and Encode](https://www.base64encode.org/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

# Shell

## testing

```sh
# makefile defaults to dev env
make t
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
