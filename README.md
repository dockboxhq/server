# dockbox

Welcome! This is the backend for `dockbox`, an application that gets you started quickly with on-the-go coding environments. Ever wanted to try out a library quickly before trying it out in your project? Have you ever stressed about having a lot of unwanted resources tied up on your machine? `dockbox` aims to solve these problems.


## Server

The server is written in Go using the [Gin](https://github.com/gin-gonic/gin) framework. It supports and manages websocket connections from various client, and proxies websocket connections to the docker server

### Prerequisites
The environment variables are loaded from `utils/config.go`. These values need to be defined in the local file before starting the server; otherwise, an error will be thrown. 

Currently, the following values are required in a `.env` file in the root directory:

```
type configType struct {
	ENVIRONMENT        string
	PORT               string
	DOCKER_SERVER_HOST string
	DATABASE_NAME      string
	DATABASE_HOST      string
	DATABASE_PORT      string
	DATABASE_USER      string
	DATABASE_PASSWORD  string
	MOUNT_POINT        string
}
```
Note: the mount point is where the user data will be stored. 

To set up for production, the set up EC2 script can be used.

### Architecture

Currently, these are the components in the architecture of the server:
- Docker server(s)
- Go server(s)
- Mounted shared file system for each server
- Database


The specifics of the architecture is implemented with AWS resources.

- Servers: EC2 machines running Ubuntu 20.0.2
- File System: EFS (Elastic File System)
- Database: RDS (PostgreSQL) 

Along with the above, the resources are running in a single VPC and spans multiple AZ in `us-east-1` (North Virginia). An Application Load Balancer (ALB) is used for load balancing and terminates SSL. An ALB instead of a classic Load Balancer was used since it supports WebSockets.


### Functions
 
Currently, there is support for creating, retrieving, and connecting to a `dockbox`. A dockbox is an environment for your project, and contains details about the source repository; its ID can be used to connect a live websocket to it. 
From the UI, a dockbox would be a terminal window for the user.

To prevent the server from being overloaded, active containers need to managed properly. The ConnectionManager object takes care of the specifics, and is thread-safe. The server keeps track of the number of connections for each container, and once the number of connections for a specific container reaches 0, a new goroutine is spawned and the container is added to the deletion queue. Once the container is in the deletion queue for 1 minute, the container is stopped. 



## Credits
- https://github.com/koding/websocketproxy
- https://github.com/gin-gonic/gin



