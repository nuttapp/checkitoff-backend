Backend for Checkitoff
======================

Architecture
------------

### Controllers
Are the entry point of the application
Should handle data scrubbing and authentication
Should enqueue/dequeue messages to/from NSQ
Import models

### DAL
Should handle interaction with the database
Should use models to serialize/deserialize data
Import models

### Models
Represent high level logical resources used by the application. 
Ex: List, ListItem, User etc...
Should handle, serialization, deserialization, and basic validation

### Binaries

#### Persistor
Dequeues messages from NSQ and stores in them Cassandra
Enqueues responses on NSQ with the reply from Cassandra
Import DAL & Models

