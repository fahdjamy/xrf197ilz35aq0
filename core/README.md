## The core package.

This package contains the core business domains and logic for this application.

The project is inspired by two Architectures thus Dependency Inversion ().
1. [Hexagonal Architecture (Onion Architecture)](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)).
   1. Read more here ([link-1](https://medium.com/ssense-tech/hexagonal-architecture-there-are-always-two-sides-to-every-story-bc0780ed7d9c))
2. [DDD](https://en.wikipedia.org/wiki/Domain-driven_design) (Domain Driven Development)

### Models

At very core lies models.

### Services

Services orchestrate the interaction between the domain and the outside world. Services implement
Keep in mind we have two types of services, we have core services that orchestrate the interaction of the outside world
with the domain models. they facilitate services to the domain models before they are exposed to the outside.
All interaction from adapters (like)

Services in the core define interfaces that adapters (external) like the http, gRPC, Kafka services would call to interact with
out business logic but in turn they depend on interfaces defined by other adapters (internal) like databases.

For-example, a user service defines an interface on how to create a user