# ritalobot ![Alt text](https://travis-ci.org/blackdev1l/ritalobot.svg?branch=master)

telegram bot written in golang which uses Markov Chain stored in redis

Installation
------------
you need *golang* >= *1.3* and *redis* installed on your machine.  
`go get github.com/blackdev1l/ritalobot`

Usage
------------

#### flags
```
flag | default | help
-c="./config.yml": path for ritalobot config
-conn="tcp": type of connection and/or ip of redis database
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
- [ ] increase / decrease chance from command
- [ ] better Markov chain
- [ ] command to start or stop bot


![works on my machine](http://www.edsquared.com/content/binary/Windows-Live-Writer/dbb6c39a79dc_68DE/WorksOnMyMachine_3.png)
