# WASA (Web And Software Architecture Project).

# How to run (in development mode)
You can launch the backend only using:
```sh
go run ./cmd/webapi/
```

If you want to launch the WebUI, open a new tab and launch:
```sh
./open-node.sh
# (here you're inside the container)
yarn run dev
```

# Docker
To launch the entire application using docker:
```sh
docker-compose up
```
