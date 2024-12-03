# "Word of Wisdom" TCP server and client with protection from DDOS based on Proof of Work

## 1. Description
This project is an implementation of “Word of Wisdom” tcp server and client.
TCP server is protected from DDOS attacks with the [Poof of Work](https://en.wikipedia.org/wiki/Proof_of_work) approach,
the challenge-response protocol should be used.  

## 2. Getting started
### 2.1 Requirements
+ [Go 1.22+](https://go.dev/dl/) installed (to run tests, start local server or client)
+ [Task](https://taskfile.dev/installation/) installed (to run any task)
+ [Docker](https://docs.docker.com/engine/install/) installed (to start server and client with docker-compose)

### 2.2 Start server and client with docker-compose:
```
task start
```

### 2.3 Start local server:
```
task start-server
```

### 2.4 Start local client:
```
task start-client
```

### 2.5 Launch tests:
```
task test
```

## 3. Problem description
Design and implement “Word of Wisdom” tcp server. 
TCP server should be protected from DDOS attacks with the [Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work), 
the challenge-response protocol should be used.  
The choice of the PoW algorithm should be explained.  
After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
Docker file should be provided both for the server and for the client that solves the PoW challenge.

## 4. Protocol description
This solution uses TCP-based binary protocol with [protobuf](https://protobuf.dev) as serializing mechanisms.
Proto file is located in `./server/pkg/proto/proto.proto`.

The shared Go implementation of the interaction protocol is stored in `./server/pkg/proto` and is available as a library for integration into any Go project, enabling the creation of custom independent clients.

## 5. Proof of Work
Idea of Proof of Work for DDOS protection is that client, which wants to get some resource from server, 
should firstly solve some challenge from server. 
This challenge should require more computational work on client side and verification of challenge's solution - much less on the server side.

### 5.1 Selection of an algorithm
There is some different algorithms of Proof Work. 
I compared next three algorithms as more understandable and having most extensive documentation:
+ [Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree)
+ [Hashcash](https://en.wikipedia.org/wiki/Hashcash)
+ [Guided tour puzzle](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol)

After comparison, I chose Hashcash. Other algorithms have next disadvantages:
+ In Merkle tree server should do too much work to validate client's solution. For tree consists of 4 leaves and 3 depth server will spend 3 hash calculations.
+ In guided tour puzzle client should regularly request server about next parts of guide, that complicates logic of protocol.

Hashcash, instead has next advantages:
+ simplicity of implementation
+ lots of documentation and articles with description
+ simplicity of validation on server side
+ possibility to dynamically manage complexity for client by changing required leading zeros count

Of course Hashcash also has disadvantages like:

1. Compute time depends on power of client's machine. 
For example, very weak clients possibly could not solve challenge, or too powerful computers could implement DDOS-attackls.
But complexity of challenge could be dynamically solved by changing of required zeros could from server.
2. Pre-computing challenges in advance before DDOS-attack. 
Some clients could parse protocol and compute many challenges to apply all of it in one moment.
It could be solved by additional validation of hashcash's params on server. 
For example, on creating challenge server could save **rand** value to Redis cache and check it's existence on verify step.

But all of those disadvantages could be solved in real production environment. 

## 6. Project Structure

The project is divided into two separate Go projects located in `./server` and `./client`.

Each Go project 
is implemented in accordance with the [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) principles 
and
adheres to the  [Go-layout](https://github.com/golang-standards/project-layout) pattern.


Server directories:
- ./server/cmd/main - server main.go
- ./server/configs - server config files
- ./server/internal/config - server config model and logic
- ./server/internal/entity - server entities
- ./server/internal/inner - server inner layer (business logic aka use-cases, manager, services)
- ./server/internal/outer - server outer layer (tcp-server, repository, etc.)
- ./server/pkg/hashcash - implementation of Hashcash PoW approach 
- ./server/pkg/proto - implementation of the interaction protocol
- ./server/pkg/utils - some auxiliary utilities

Client directories:
- ./client/cmd/main - client main.go
- ./client/configs - client config files
- ./client/internal/config - client config model and logic
- ./client/internal/outer - client outer layer (tcp-client)


## 7. Ways to improve
Of course, every project could be improved. This project also has some ways to improve:
+ add dynamic management of Hashcash complexity based on server's overload 
(to improve DDOS protection)
+ move the array of quotes to an external data storage, such as **SQLite** or **PostgreSQL**.
+ add integration tests for simulate DDOS attack - spawn more client instances