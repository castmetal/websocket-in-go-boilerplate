# Websocket Go Boilerplate

## Case 1: Publish Message and Receive a Result

![how it Works](https://raw.githubusercontent.com/castmetal/websocket-in-go-boilerplate/main/WebsocketFlow1.png)

> Easy Websocket Boilerplate in Go

- This version was adapted from [eranyanay](https://github.com/eranyanay/1m-go-websockets/)

## How it Works

- Publish messages on server and receive a result of your use case and business rules back with a simple full-duplex message at TCP connections
- Send messages with a default `AUTH_HEADER` environment and receive again a result for this publish

## Case 2: Publish to All Clients

![how it Works](https://raw.githubusercontent.com/castmetal/websocket-in-go-boilerplate/main/WebsocketFlow2.png)

- Publish a message to all TCP active connections. This example is common when you need to send to all clients a state or a new notify message about a new content and you need update your app with this

## Case 3: Publish to A Specific User

![how it Works](https://raw.githubusercontent.com/castmetal/websocket-in-go-boilerplate/main/WebsocketFlow3.png)

- Publish a message to a specific user with TCP active connections. This example is common when you need to send to a user about changes in your contents, process, states or rules.

## Subject

- This code was developed to handle many requests and simultaneous connections

### Running at local environment

#### Change Your .env.example and copy to another file named .env

- Install Postgresql and create your database with the same DB_DATABASE_NAME value in your .env file

#### Execute Migrations. Example in PostgreSQL

- Run:

```sh
    go run src/infra/migration/setup.go
```

#### Run your server

- Run:

```sh
    go run server.go
```

#### Executing Client 1 case

- Run:

```sh
    go run src/examples/execute-use-case/client.go YourUserId
```

#### Executing Client 2 case

- Estabilish a Single Connection:

```sh
    go run src/examples/connect/client.go YourUserId
```

- Run:

```sh
    go run src/examples/write-to-all-clients/client.go YourUserId
```

#### Executing Client 3 case

- Estabilish a Single Connection:

```sh
    go run src/examples/connect/client.go YourUserId
```

- Run:

```sh
    go run src/examples/write-to-an-user/client.go YourUserId
```

### To do

- For further security rules, add auth header as JWT and valid iss and exp with low timestamp tls
