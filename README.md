# Disel - A simple HTTP Framework Written in Go
Disel, is a simple HTTP framework which can handle serving HTTP requests, routing and logging out of the box named after Shaq the most dominant center ever to play Ball. 

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Data Structures Used](#data-structures-used)
- [Usage](#usage)


## Overview
This project is an extension of a few challenges  I am doing from codecrafters. To know more about these challenges please visit these repos:
 - Build your own HTTP Challenge - https://github.com/harish876/harish876-codecrafters-http-server-go  [COMPLETED]
 - Build your own REDIS Challenge - https://github.com/harish876/harish876-codecrafters-redis-go       [IN_PROGRESS]
 - Build your own DNS Challenge -   https://github.com/harish876/harish876-codecrafters-dns-server-go  [IN_PROGRESS]

## Features
1. Establishes a simple TCP Connection with a client and parses out the packets based on the HTTP Protocol.
2. Spins up a go routine for each request or provides a pooling interface to handle requests.
3. Provides routing out of the box using a nice Radix Tree implementation.
4. Provides structured logging out of the box to log requests, responses and bind the logger object to each request handler interface.

## Dependecies Used
1. No Third party dependencies for the core application is used. Used a simple terminal based ASCII art generator to provide neat welcome messages.

## Data Structures Used
1. Router Implemented using a nifty optimization on the trie data structure, which is the **Radix Tree**. Very fun to implement. 
2. Each HTTP Method like GET or POST has a radix tree to maintain the routes and their handlers. This way routes with long prefixes are stored more efficiently. 
3. Nify little **Thread pool** with a capacity implemented in case spinning go routines for each 

## Future Plans 
1. Adding support for route groups, and some middleware features like basic auth and rate limiting.
2. I want to be able to create a framework which can be used to not only implement HTTP requests but any protocol based on TCP.
3. Like the RESP protocol, for redis, so that just the interface for the serialisation and deserialisation needs to be implemented.

## Usage
1. Have golang installed and run go run main.go
2. For development and HMR air command can be used.

