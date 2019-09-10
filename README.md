# Why

This is a simple project to test Redis connectivity.


# Deploy on SquareScale

Use one of their methods:

  - `sqsc-demo-redis` [docker image](https://cloud.docker.com/repository/docker/nledez/sqsc-demo-redis/)
  - Build from [https://github.com/nledez/sqsc-demo-redis/](https://github.com/nledez/sqsc-demo-redis/)

Available environment variables to add:
  - `REDIS_NAME`: The name of the Redis database in Squarescale
  - `REDIS_${REDIS_NAME}_ADDRESS`: The address of the Redis database
  - `REDIS_${REDIS_NAME}_PASSWORD`: The password of the Redis database
  - `REDIS_${REDIS_NAME}_PORT`: The port of the Redis database
  - `REDIS_${REDIS_NAME}_DB`: The DB index of the Redis database
