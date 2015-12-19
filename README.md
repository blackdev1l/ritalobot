# ritalobot ![Alt text](https://travis-ci.org/blackdev1l/ritalobot.svg)

telegram bot written in golang which uses Markov Chain stored in redis

Installation
------------
`go get github.com/blackdev1l/ritalobot`

Usage
------------

#### flags
```
flag | default | help
-c="./config.yml": path for ritalobot config
-conn="tcp": type of connection and/or ip of redis database
-id=0: Chat id of the group chat
-p=6379: port number of redis database
-token="": authentication token for the telegram bot
```

#### config file
create a `config.yml` or rename `example.yml` editing the variables.
make sure to have redis-server started.

TODO
------------

- [x] Flags
- [x] yaml config
- [ ] custom temporization
- [ ] better Markov chain
- [ ] command to start or stop bot
