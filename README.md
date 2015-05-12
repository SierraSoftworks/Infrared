# InfraRed
**Infrastructure configuration management**

InfraRed provides a configuration management endpoint which is intended to be consumed by services operating within a private cloud.
It offers a way to quickly and easily register and configure nodes on your network, as well as make those nodes discoverable for consumers.

## The API
InfraRed provides a REST API which allows you to register nodes under a specific `node_type` as well as assigning and retrieving configuration
information for each of those `node_type`s.

### Node API
The Node API allows you to easily manage registered nodes on the InfraRed server.

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