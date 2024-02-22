# hive-docker

start your containers from a simple API

- start your containers from a simple API
- mount a binding 

## use case

build the Dockerfile in 'images' folder

```
docker build . t test:latest
```

start the container from api

```
curl  -X POST \
  'localhost:8080/docker/start' \
  --data-raw '{
  "folderLocation":"C:\\Users\\YOUR_USER",
  "imageName":"test:latest"
}'
```