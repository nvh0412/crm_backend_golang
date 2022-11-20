<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# CRM Backend

### Status

![main workflow](https://github.com/nvh0412/crm_backend_golang/actions/workflows/go.yml/badge.svg)

## About the project

The project represent the backend of a Customer Relationship management (CRM) web application. As users interact with app via some UI, our server will support all of the funtionalities via API:

- Getting a list of all customers
- Getting data of a single customer
- Adding a customer
- Updating a customer's information
- Removing a customer

### API docs

- Getting a single customer through a /customers/{id} path
```
$ curl --location --request GET 'localhost:3000/customers/1'
```

- Getting all customers through a the /customers path
```
$ curl --location --request GET 'localhost:3000/customers'
```

- Creating a customer through a /customers path
```
$ curl --location --request POST 'localhost:3000/customers' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Hoa Nguyen",
    "role": "Software Engineer",
    "email": "nvh0412@gmail.com",
    "phone": "123123123",
    "contacted": false
}'
```

- Updating a customer through a /customers/{id} path
```
$ curl --location --request PUT 'localhost:3000/customers/1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Hoa Nguyen",
    "role": "Manager"
}'
```

- Deleting a customer through a /customers/{id} path
```
$ curl --location --request DELETE 'localhost:3000/customers/1'
```

## Getting started

### 1. Setup ENV

Make sure that we've configured these system variables on our dev machine
```
    APP_DB_USERNAME: ...
    APP_DB_PASSWORD: ...
    APP_DB_NAME: ...
```

### 2. Get dependencies

```
    $ cd <project_path>
    $ go get
```

### 3. Build project

```
    $ go build .
```

### 4. Run project

```
    $ ./crm_backend_golang
```

### 5. Run test

```
    $ go test .
```


