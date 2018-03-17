# fpeService
a service that encrypts your PII using format preserving encryption

# TLS
The server requires a certificate and private key in the same folder it runs from. 

for testing purposes run 

```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
```

making sure you put `127.0.0.1:8080` in the "Common Name".

# API

* `GET /encrypt` attach a body to your GET request and it'll come back encrytped
