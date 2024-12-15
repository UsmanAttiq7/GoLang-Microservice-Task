# GO Microservices
This repository contains 3 microservices:
 * **User Service**: This service is responsible for user management. 
 * **Booking Service**: This service is responsible for booking management.
 * **Ride Service**: This service is responsible for ride management.

## Getting Started
These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

 * **POSTGRES** - The root of the repo contains a docker folder with a `docker-compose` file that will spin up a PG database with initial seed data for the services.
 * **GO** - The services are written in Go. You will need to have Go installed on your machine.

Run `go mod tidy` to install all the dependencies for each service which are defined in the `go.mod` file of each the respective service folder.

Run `go run <service>/cmd/main.go` to start the services.

Each service has `.env` pre-configured to run with the docker-compose PG database on port `5432`.

### User Service
After starting the user service, you can access the user service on `http://localhost:50051`. 
Use the following grpcurl commands to interact with the user service:
* Get a User by user_id
```shell
grpcurl -plaintext -d '{"user_id": 1}' localhost:50051 user.v1.UserService/GetUser
```

* Create a User
```shell
grpcurl -plaintext -d '{"name": "Bilal"}' localhost:50051 user.v1.UserService/CreateUser
```

* Delete a User
```shell
grpcurl -plaintext -d '{"user_id": 4}' localhost:50051 user.v1.UserService/DeleteUser
```

### Booking Service
After starting the booking service, you can access the booking service on `http://localhost:50052`.
Use the following grpcurl commands to interact with the booking service:

* Get a Booking by booking_id
```shell
grpcurl -plaintext -d '{
  "booking_id": 1  
}' localhost:50052 booking.v1.BookingService/GetBooking
```

* Create a Booking
```shell
grpcurl -plaintext -d '{
  "user_id": 1,
  "ride": {
    "source": "Downtown",
    "destination": "Airport",
    "distance": 15,
    "cost": 250
  }
}' localhost:50052 booking.v1.BookingService/CreateBooking
```

### Ride Service
After starting the ride service, you can access the ride service on `http://localhost:50053`.
Use the following grpcurl commands to interact with the ride service:

* Update a Ride by ride_id
```shell
grpcurl -plaintext -d '{
  "ride_id": 1,
  "ride": {
    "source": "Downtown",
    "destination": "Mall",
    "distance": 10,
    "cost": 200
  }
}' localhost:50053 ride.v1.RideService/UpdateRide
```