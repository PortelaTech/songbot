# Songbot ![Alt text](https://travis-ci.org/blackdev1l/ritalobot.svg?branch=master)

Telegram bot written in golang which uses Markov Chain stored in redis

Follow this [project layout](https://github.com/golang-standards/project-layout).

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

    - go build
    - go test

TODO
------------

- [x] Flags
- [x] yaml config
- [ ] increase / decrease chance from command
- [ ] better Markov chago run botin
- [ ] command to start or stop bot



# Adaption of Hazel For Juan

## Original docs

- Part 1: [https://aranair.github.io/posts/2016/12/25/how-to-set-up-golang-telegram-bot-with-webhooks/][1]
- Part 2: [https://aranair.github.io/posts/2017/01/21/how-i-deployed-golang-bot-on-digital-ocean/][2]
- Part 3: [https://aranair.github.io/posts/2017/08/20/golang-telegram-bot-migrations-cronjobs-and-refactors/][3]

### Commands
- Hazel, anything: `replies with korean Hello :P`
- remind buy dinner
- remind me to do this and this
- remind me to sleep :9jul 10pm
- remind me to buy chocolate :today 10pm
- remind me to buy a gift :tomorrow 10pm
- clear 2
- clearall
- list

## Modifications for Juan

### Build

Running commands in Dockerfile locally
```
go get github.com/tools/godep
go run github.com/tools/godep restore
go install ./...
```
Actually we should use offical [dep](https://github.com/golang/dep) for dependencies.

### Create a self signed TLS certificate
```
# PKCS#10 certificate request and certificate generating utility.
openssl req -new > cert.csr 

Trace of output - passphrase and challenge password are both "password"
Country Name (2 letter code) [AU]:BR
State or Province Name (full name) [Some-State]:Bahia
Locality Name (eg, city) []:Salvador
Organization Name (eg, company) [Internet Widgits Pty Ltd]:Portela
Organizational Unit Name (eg, section) []:Juan
Common Name (e.g. server FQDN or YOUR name) []:t.me/JuanTestBot
Email Address []:juan.rgarcia@hotmail.com

# Write RSA key
openssl rsa -in privkey.pem -out key.pem 

# -signkey key.pem == self sign using supplied private key
openssl x509 -in cert.csr -out cert.pem -req -signkey key.pem -days 1001

# Add the private key to the certificate
cat key.pem>>cert.pem
```

### Create a bot

Send ```/newbot``` to the [Botfather](https://telegram.me/botfather)
```
{botId}:{apiKey} => 660697313:AAHmZ43BVU0ZjLI1JoRIlI6ySsABsUIoo70
```

[Run the getMe method on this bot](https://api.telegram.org/bot660697313:AAHmZ43BVU0ZjLI1JoRIlI6ySsABsUIoo70/getMe)

Set up a webhook by sening this curl request 
```
curl -F "url=https://t.me/JuanTestBot" -F "certificate=@cert.pem" https://api.telegram.org/bot660697313:AAHmZ43BVU0ZjLI1JoRIlI6ySsABsUIoo70/setWebhook
```

### Run

- Create a `configs.toml` file
- `docker-compose up` or `go run cmd/webapp/main.go`

### Deploy

- Set up git hooks in production
- `git push production master`

Sample post-receive hook

```bash
#!/bin/sh

git --work-tree=/var/app/remindbot --git-dir=/var/repo/site.git checkout -f
cd /var/app/remindbot
docker-compose build
docker-compose down
docker-compose -d
```

### Debug

```
go build -gcflags "all=-N -l" github.com/PortelaTech/remindbot/cmd/webapp
dlv --listen=:2345 --headless=true --api-version=2 exec ./webapp
```

### License
MIT
