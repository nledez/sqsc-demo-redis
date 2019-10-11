# Why

This is a simple project to test Redis connectivity.


# Deploy on SquareScale

Use one of their methods:

  - `sqsc-demo-redis` [docker image](https://cloud.docker.com/repository/docker/nledez/sqsc-demo-redis/)
  - Build from [https://github.com/nledez/sqsc-demo-redis/](https://github.com/nledez/sqsc-demo-redis/)

You need to define the `REDIS_INSTANCE_NAME` environnment variable which
contains the name of the redis instance you provisionned in Squarescale.

All other environment variables as set in Squarescale will be automatically
populated.

This server listens to port 8081.
