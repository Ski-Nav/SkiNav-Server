# SkiNav-Server

## Introduction

SkiNav-Server is the backend server for SkiNav, a navigation application for ski resorts. It provides two APIs that allow users to retrieve information about available resorts and access specific resort map data.

## APIs

1. **/maps**: This API endpoint provides a list of all available resort names. Users can access it by making a GET request to the following URL: **http://ec2-18-188-212-45.us-east-2.compute.amazonaws.com:3000/api/v1/maps**.
2. **/maps/:resortName**: This API endpoint provides the map data for a specified resort. Users can retrieve the map data by appending the resort name as a path parameter to the URL. For example, to get the map data for the "UCSD" resort, make a GET request to **http://ec2-18-188-212-45.us-east-2.compute.amazonaws.com:3000/api/v1/maps/UCSD**.

## Architecture
The SkiNav-Server architecture consists of three key components, as shown in the diagram below:

![Architecture](./assets/architecture.png)
1. **Golang Gin Server**: This component handles user requests and generates responses based on the requested data. It serves as the main interface between the client applications and the backend server.
2. **Neo4j Database**: The Neo4j database stores the map data for ski resorts. It provides efficient storage and retrieval of the map data, allowing the server to quickly access the required information.
3. **Memcached**: Memcached is utilized as a memory cache to optimize the overall performance of the server. It stores frequently accessed data in memory, reducing the need to fetch it from the database on every request. This caching mechanism improves response times and enhances the user experience.
4. 
By combining these three components, SkiNav-Server provides a robust and efficient backend solution for serving ski resort map data to client applications.