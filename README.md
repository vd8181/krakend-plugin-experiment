This is a POC on the rate limiting aspect of krakend.  

I have built a basic backend of a journal app using springboot and REST Principles.  

You can run ./gradlew bootRun to run the spring app.  

Along with that , inorder to get the krakend server up and running , you must build the container first using "docker build krakend ." .  

After that you can run the container using "docker run -p 3000:3000 krakend" .   

Now , you must also set the redis container up and running using docker run -p 6379:6379  -v ${PWD}/src/redis/data:/data redis .  

The krakend server is running in port 3000 of host machine , redis is running on port 6379 and the spring boot app is running on port 8000.  

You can now use postman to send requests to krakend first which is running on port 3000 and then krakend would forward the requests to the respective endpoints. 



redis.conf: Configures Redis server settings, including persistence, directory paths, and other operational parameters.
dump.rdb: A binary file representing a snapshot of the Redis database, used for fast loading of data upon server start.


1. Dockerfile
Purpose: This file is used to build Docker images. It outlines the steps to compile Go plugins and integrate them with the KrakenD API Gateway.
Stages:
Builder Stage: Compiles Go plugins (validator.so, rate_limiting.so) using the Go language and dependencies.
Final Stage: Sets up KrakenD with the compiled plugins and configuration, exposing the default KrakenD port for use.
2. go.mod File
Purpose: Manages dependencies for the Go project.
Content: Lists the required versions of Go and dependencies. It specifies indirect dependencies used by the project.
3. go.sum File
Purpose: Ensures integrity and authenticity of the dependencies listed in go.mod.
Content: Contains checksums for the specific versions of the modules, ensuring they haven’t been tampered with.
4. KrakenD Configuration File (krakend.json)
Purpose: Configures the KrakenD API Gateway.
Content: Defines endpoints, rate limiting strategies, and plugins. Specifies how requests are routed and handled, integrating with backend services.
5. Go Plugin Code (rate_limiting.go)
Purpose: Implements a rate limiting plugin for KrakenD.
Content: Intercepts HTTP requests, checks API keys against Redis, and applies headers for rate limiting based on tiers.
6. Go Plugin Code (validator.go)
Purpose: Implements a validator plugin for KrakenD.
Content: Validates signatures using RSA public keys. Intercepts requests to verify signatures against a stored public key, ensuring requests are authorized before proceeding. This file mainly aims to implement the certificate
validation functionality that provisioning service does. This ensures that the certificate validation runs in the krakend service itself.

Context in a Spring Boot Project:
Integration: These Go plugins likely interact with the Spring Boot application’s REST endpoints via the KrakenD API Gateway. The Spring Boot app could be part of the backend that KrakenD routes requests to after applying validation and rate limiting.
Microservices Architecture: This setup suggests a microservices architecture where different components (Spring Boot app, Go plugins, API Gateway) work together to provide a cohesive service.


