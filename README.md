# go-3um
anonymous bbs with almost no features

## setup
get dependencies (upper.io db functions, gorilla mux + securecookie):
```
go get -u github.com/realytcracker/go-3um
```

setup your ssl bullshit:

```
openssl genrsa -out server.key 4096
openssl ecparam -genkey -name secp384r1 -out server.key
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
echo "" >> server.key
cat server.crt >> server.key
echo "" >> server.key
```
setup your mysql bullshit:
```
mysql -uroot -p
[enter password]
CREATE DATABASE 3um;
[control-D]
cat 3um.sql | mysql -uroot -p 3um
```

rename `config.defaults.json` to `config.json` and edit the values within properly.

`go build` and run the resulting binary.

visit `https://host:8443/api/setup` and receive your admin credentials.

## todo

`[ ] add additional user apis`

`[ ] frontend`

`[ ] rate limiting`

`[ ] dockerize and shit`

## remember eternal

```
sd1
smurda
brand0n
adrian
lxuke
justincredible
toyo4321
olaf
rj2
ib
ackflags
christophermichael
maru
didi
anti
```


