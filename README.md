Design and implement “Word of Wisdom” tcp server

### Requirements

* TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
* The choice of the POW algorithm should be explained.
* After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
* Docker file should be provided both for the server and for the client that solves the POW challenge


### How to

Run server:

```sh
make run_server
```

Run client:

```sh
make run_client SERVER_HOST=$(ipconfig getifaddr en0):8081
```



### Why sha256?

1. Its well known and widely used (in Bitcoin for example)
2. Considered to be highly secure (not a signel collision found yet)
3. Due its efficiency which can also be implemented in hardware
