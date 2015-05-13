# InfraRed
**Infrastructure configuration management**

InfraRed provides a configuration management endpoint which is intended to be consumed by services operating within a private cloud.
It offers a way to quickly and easily register and configure nodes on your network, as well as make those nodes discoverable for consumers, as well as providing a centralized configuration source for those nodes.

## Requirements
 - [MongoDB](https://www.mongodb.org/)

## Example Usage
```sh
irclient register http://localhost:8080 mongodb
irclient beam http://localhost:8080
```

```sh
irserver config.json
```

## Command Line Options
```
irclient command server [nodetype] [hostname] [port]

commands:
  register  - Registers this node with the specified server
  beam      - Submits a stream of heartbeats to the server until the process is terminated
  
server   - The Infrared server managing the nodes
nodetype - The type of node to register
hostname - The hostname or IP address to be used when accessing this node
port     - The port number of the service running on this node
```

```
irserver config [listenOn] [dbServers] [database]

config    - The configuration file specifying the server's various options
listenOn  - A listen specification like :8080
dbServers - The database servers to connect to, something like dbhost1,dbhost2:27016
database  - The name of the database to store information in
```

## Configuration

```json
{
  "listen": ":8080",
  "database": {
    "url": "localhost",
    "db": "infrared"
  }
}
```

## The API
InfraRed provides a REST API which allows you to register nodes under a specific `node_type` as well as assigning and retrieving configuration
information for each of those `node_type`s.

### Node API
The Node API allows you to easily manage registered nodes on the InfraRed server.

#### List of Online Nodes
**GET /api/v1/:node_type**

Gets the list of nodes of type `node_type` which have recently triggered a heartbeat.

```json
[{
  "id": "abcdef1234567890",
  "type": "mongodb"
  "hostname": "127.0.0.1",
  "port": 27016,
  "lastSeen": "1970-01-01T00:00:00.000Z"
}]
```

#### Node Heartbeat Endpoint
**GET /api/v1/:node_type/:id/heartbeat**

Submits a heartbeat request to InfraRed to inform it that this node is still online.

#### Retrieving a Node's Details
**GET /api/v1/:node_type/:id**

Requests the details of a specific node from the InfraRed server.

```json
{
  "address": "127.0.0.1",
  "port": 1234,
  "online": true	
}
```

#### Registering a Node
**POST /api/v1/:node_type**
```json
{
  "address": "127.0.0.1",
  "port": 1234	
}
```

Creates a new node on the InfraRed server, the server will respond with a unique identifier for this node.

```json
{
  "id": "abcdef123456789"	
}
```

#### Updating a Node's Details
**PUT /api/v1/:node_type/:id
```json
{
  "address": "127.0.0.1",
  "port": 1235	
}
```

Updates the details associated with a node in InfraRed - for example, if the node changes its IP address or listening port.

#### Removing a Node
**DELETE /api/v1/:node_type/:id**

Removes a node from InfraRed.

### Configuration API
The configuration API allows you to specify configuration options for specific node types and then subsequently
retrieve them for later use.

#### Get Configuration
**GET /api/v1/:node_type/config**

Gets the configuration associated with the `node_type` on the InfraRed server.

```json
{
  "myOption": true	
}
```

#### Set Configuration
**PUT /api/v1/:node_type/config**
```json
{
  "myOption": true 
}
```

Sets the configuration options for a specific node type.

### UDP Heartbeat Endpoint
In the interest of performance, you can send heartbeats to a UDP endpoint. These heartbeats are encoded using ProtocolBuffers and defined by the
following protocol specification.

```
package infrared;

message Heartbeat {
  required string id = 1;
  required string node_type = 2;
}
```