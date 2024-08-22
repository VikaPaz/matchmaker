# Matchmaker

A system for matchmaking players in multiplayer games based on their skills and network latency.

## Discription
Matchmaker is a player matching system for multiplayer games that automatically assigns players to matches based on their skills and network latency (ping). 

### Key Features

- **Dynamic Player Addition:** The system supports the dynamic addition of new players during the game, allowing for real-time matchmaking.
- **Redis Integration:** Players are stored in Redis, which enables fast access and management of player data.
- **PostgreSQL Storage:** Player details, such as profiles and historical data, are stored in PostgreSQL for persistent storage and querying.

### Matching Algorithms

- **K-Means Clustering:**
  - **Description:** K-means clustering is an algorithm used to partition players into clusters based on their skill levels and latency. Each player is assigned to the nearest cluster center, which is recalculated iteratively to minimize the distance between players and their assigned cluster centers.
  - **Benefits:** This approach helps group players with similar attributes together, improving the balance and fairness of matches.

- **Euclidean Distance:**
  - **Description:** Euclidean distance is used to measure the similarity between players based on their skill and latency. It calculates the straight-line distance between two points in a multidimensional space, where each dimension represents a different attribute (e.g., skill and latency).
  - **Benefits:** This metric helps in determining how close players are to each other in terms of their attributes, which is crucial for effective clustering and matchmaking.

The system combines these algorithms to ensure players are matched with others of similar skill levels, enhancing the overall gaming experience. The integration of Redis and PostgreSQL ensures efficient data handling and reliable storage.

## Installation

1. **Install dependencies:**
Make sure you have the following components installed:
    - Go (version go1.22.2 or higher)
    - Redis
    - PostgreSQL

2. **Clone the repository:**
    ```bash
    git clone https://github.com/yourusername/matchmaker.git
    cd matchmaker
    ```

3. **Set Go Dependencies:**
    ```bash
    go mod tidy
    ```

4. **Start the containers:**
Use the `make` command to start the PostgreSQL and Redis containers. Make sure you have the `docker-compose.yaml` file in the `build` folder configured to start the required services:
    ```bash
    make docker_compose_build
    ```

5. **Set up a PostgreSQL database:**
- Create the database:
    Open the PostgreSQL command prompt and create a database for the project:
    ```sql
    CREATE DATABASE users;
    ```
- Perform the migrations:
    The ``migrations`` folder contains the SQL scripts needed to create the database schema. Execute these scripts to configure tables and other database objects.

## Usage.

1. **Starting the server:**
After successfully completing all installation and configuration steps, start the server:
   ```bash
   go run main.go
   ```

2. **API Documentation:**
API documentation is available through Swagger. You can access the Swagger UI at the following URL: http://localhost:8900/swagger/index.html

3. **Example Requests:**
Adding a player: To add a player to the system, use the following HTTP request:
    ```http
    POST /matchmaker/users
    Host: localhost:8900
    Accept: application/json
    Content-Type: application/json

    {
    “latency": 30,
    “name": "player1",
    “skill": 50
    }
    ```

## Configuration
дисбалансу в матчах
To customize your project, you can use environment variables. Here are the configuration settings you can customize:

1. **Environment Variables:**
   - `GROUP_SIZE` - The size of the player group. The default value is `3`.
   - `GROUP_WAIT` - Waiting time in milliseconds before creating a group. Default value: `5000`.
   - `HOST` - The host on which the server is started. Default value: `localhost`.
   - `PORT` - The port on which the server is running. Default value: `8900`.
   - `REDIS_PORT` - The Redis port. Default value: `6379`.
   - `PASSWORD` - Password for accessing Redis. The default value is `password1`.
   - `LOGS_LEVEL` - Logging level. Default value: `Info`.
   - `MAX_SKILL` - Maximum level of player's skill. Default value: `10000`.
   - `MAX_LATENCY` - Maximum delay (latency) of the player. Default value: `300`.
   - `POSTGRES_PORT` - PostgreSQL port. Default value: `5432`.
   - `USER` - PostgreSQL user. Default value: `user`.
   - `POSTGRES_PASSWORD` - PostgreSQL user password. Default value: `password`.
   - `DB_NAME` - The name of the PostgreSQL database. Default value: `users`.
   - `WRITE_TO_POSTGRES` - Flag to write data to PostgreSQL. Default value: `true`.

### Optimization of the clustering algorithm

**Current Problem:**.
When the number of users is small, the K-means algorithm may inefficiently assign players to clusters, resulting in imbalanced matches and increased waiting time.

**Planned improvements:**
1. **Adaptation for small number of users:**
   - Implement a method to transfer players between clusters to improve balancing. This will help to distribute players more efficiently and create more balanced groups.
   - **Objective:** Reduce imbalance in matches and reduce waiting time by improving the distribution of players across clusters.

2. **Optimize clustering:**
   - Improve clustering algorithms to more accurately distribute players, especially with limited data. Consider using adaptive K-means that changes the number of clusters depending on the number of users.
   - **Objective:** Make the clustering process more flexible and suitable for different amounts of data, thus improving the quality of matches and reducing latency.
