--- 

# NoSQL Repository Service

The NoSQL Repository Service is a part of the ZmeyGorynych Project. 
It integrates with MongoDB and Redis to store and index data, and it uses Kafka for message brokering.

## Components

### NoSQL Repository (`zg_nosql_repo`)
This component manages the storage and indexing of data using MongoDB and Redis.

#### Docker Compose Configuration
```yaml
version: '3'

networks:
  local-net:
    external: true

services:
  zg_mongodb:
    image: mongo:latest
    container_name: zg_mongodb
    ports:
      - '27017:27017'
    volumes:
      - mongodbdata:/data/db
    networks:
      - local-net

  zg_mongodb_2:
    image: mongo:latest
    container_name: zg_mongodb_2
    ports:
      - '27018:27017'
    volumes:
      - mongodbdata:/data/db2
    networks:
      - local-net

  zg_nosql_repo:
    build:
      context: .
      dockerfile: ./Dockerfile
    container_name: zg_nosql_repo
    env_file:
      - .env-docker
    networks:
      - local-net
    volumes:
      - .:/app
    depends_on:
      - zg_mongodb
      - zg_mongodb_2

  zg_nosql_redis_index:
    image: redis:latest
    container_name: zg_nosql_redis_index
    ports:
      - '6379:6379'
    networks:
      - local-net

volumes:
  mongodbdata:
```

#### Environment Variables (`.env-docker`)
```env
ZG_KAFKA_ADDRESS=kafka:29092
ZG_MONGODB_URL_1=mongodb://zg_mongodb:27017/zg_mongodb
ZG_MONGODB_URL_2=mongodb://zg_mongodb_2:27017/zg_mongodb_2
ZG_REDIS_URL=zg_nosql_redis_index:6379
ZG_REDIS_DB=0
LOGSTASH_URL=http://logstash:5000
```

#### Configuration File (`config.yaml`)
```yaml
kafka:
  address: ${ZG_KAFKA_ADDRESS}
  group_id: nosql_repo
  user: guest
  password: guest
  topic: processing_1

dbs:
  mongodbs:
    - ${ZG_MONGODB_URL_1}
    - ${ZG_MONGODB_URL_2}

redis:
  address: ${ZG_REDIS_URL}
  db: ${ZG_REDIS_DB}

logstash:
  url: ${LOGSTASH_URL}
```

## Getting Started

### Prerequisites
- Docker
- Docker Compose

### Running the NoSQL Repository Service
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/message-generator.git
   cd message-generator/nosql-repo
   ```
2. Build and run the Docker containers:
   ```bash
   docker-compose up --build
   ```

### Environment Variables
Ensure to set the following environment variables in the `.env-docker` file:
- `ZG_KAFKA_ADDRESS`: Address of the Kafka server (e.g., `kafka:29092`).
- `ZG_MONGODB_URL_1`: URL of the first MongoDB instance (e.g., `mongodb://zg_mongodb:27017/zg_mongodb`).
- `ZG_MONGODB_URL_2`: URL of the second MongoDB instance (e.g., `mongodb://zg_mongodb_2:27017/zg_mongodb_2`).
- `ZG_REDIS_URL`: URL of the Redis server (e.g., `zg_nosql_redis_index:6379`).
- `ZG_REDIS_DB`: Redis database number (e.g., `0`).
- `LOGSTASH_URL`: URL of the Logstash server (e.g., `http://logstash:5000`).

## License
This project is licensed under the MIT License.

---

