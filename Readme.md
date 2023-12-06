# Walmart API: With Golang

<img src="./showcase/ai.png" width="100%" />

Welcome to the Walmart Ecommerce API built with Go. In this document, I document the project features, considerations made and things that could be improved.



## Introduction

This is an Ecommerce API built with Go (Golang). I chose go because of it's shear speed, performance and simplicity. Furthermore, my primary tech stack is mostly composed of Python, TypeScript, JavaScript and Go. Having used express, nestjs, flask and fastAPI, I decided to build a project with Go Fiber which boasts itself as the fastest go Rest Framework. For data storage, I decided to use MongoDB which is a NoSQL database suitable for datasets that may vary over time and it's vast search ability.

## Key Features

1. User Management

2. User Authentication with JWT and cookies

3. User Authorization & http middleware

4. Product Management & Rating

5. Saved Carts to allow synchronization.

6. Saved Products (Favorite Products)

7. HTTP Request data validation



## Tech Used

1. Golang: The Ecommerce backend is powered by Go. Go is leading in the development of cloud products and services.

2. Fiber: Fiber is a Go web framework built on top of Fasthttp, the fastest HTTP engine for Go. It's designed to ease things up for fast development with zero memory allocation and performance in mind.

3. MongoDB: Mongo is a modern NoSQL database suitable for datasets that may vary over time and it's vast search ability


# Running
To run this server, you need to clone the repo to your local environment. 

```sh
git clone https://github.com/efenstakes/walmart-api-go api
```

Navigate to the folder.

```sh
cd ./api
```

Install dependencies
```sh
go get .
```

Create a `.env` and set the below values:

1. `PORT` the port the API runs on

2. `JWT_SIGNING_KEY` The JWT signing key used to sign JWT tokens

3. `DB_URI` Mongo db instance URI (from docker, local mongo compass or cloud)


Run the Go Fiber yoda server with:
```sh
go main.go
```


Happy building :(.

## Extras
I build a similar API to this in Golang and Node.js. You can find it here in my github.


## Contact
If you wish to contact me, use my email efenstakes101@gmail.com.
