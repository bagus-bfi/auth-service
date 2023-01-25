# README

1. If incoming request to the server are from more-than-one transport (for example, gRPC, RabbitMQ, Kafka, etc) 
AND use same validation process for both, we can add a service package here which include validation 
and calling data store repository(ies) as necessary.

2. Otherwise we can just call data store from handler to reduce application layer. 
If we do this, all request validation should be done in handler level.

Note: 
- If we do #1, please make sure any errors returned from service package are identifiable from the callers.
- Use service layer if we need to encapsulate "a feature" that is gather/process data from multiple data store or external services.