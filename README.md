# SkiNav-Server
##
SkiNav-Server is the backend server of SkiNav, which provides two APIs.
1. /maps: the server will respond with all available resorts names.
2. /maps/:resortName: the server will respond with the specified resort map data.
Currently, the services is served at AWS, so you can access the APIs with the following URLs.

http://ec2-18-188-212-45.us-east-2.compute.amazonaws.com:3000/api/v1/maps

To get the resort you want, use the resort name as the path parameter. For example, http://ec2-18-188-212-45.us-east-2.compute.amazonaws.com:3000/api/v1/maps/UCSD
## Architecture
![Architecture](./assets/architecture.png)
The backend server consists of three components.
1. Golang Gin server: handles the users' requests and responds with corresponding data.
2. Neo4j database: stores the maps data.
3. Memcached: serves as a memory cache to enhance the overall performance.

