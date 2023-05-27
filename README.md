
# P2P Chat
### An encrypted alternative to exchange messages

## Setup

Clone this repository.

### Generate Keys
A simple key generation function for mac (silicon and intel) is provided in ./cli/build.
Use the respective binary to create your keys.
```bash
$ build/p2pcli-intel
or 
$ build/p2pcli-silicon
```

### frontend
The frontend will be started when running `docker-compose up` in the base directory. 
This will create the backend which communicates with the host. Furthermore, a webpage which will be the frontend 
of the application. Accessible on `http://localhost:3000`.
This can be changed in the docker-compose file in the base directory