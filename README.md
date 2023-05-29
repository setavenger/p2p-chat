
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

### Client application
To run the client you have to copy the `docker-compose-example.yaml` to a `docker-compose.yaml`.
The client will be started when running `docker-compose up` in the base directory. 
This will create the backend which communicates with the host. Furthermore, a webpage which will be the frontend 
of the application. Accessible on `http://localhost:3000`.
If you need to run it on a different port this can be changed in the docker-compose file in the base directory.

### Server
Change directory to the server subdirectory.
Copy the `docker-compose-example.yaml` to a `docker-compose.yaml`.
To run the server you have to define the domain on which the server runs in the `docker-compose.yaml`. 
Change the password (`POSTGRES_PASSWORD`) and then run the docker-compose file `docker-compose up -d`.
