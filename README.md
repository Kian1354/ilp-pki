# ilp-pki

Set up Redis 
`docker pull redis`  
`docker run --name redis-test-instance -p 6379:6379 -d redis`

Get go-redis package
`go get github.com/go-redis/redis`

## Setup

### To set up a server
Run `go run cmd/main/gen_ca/gen_ca.go` to generate the CA keys. This will generate a folder called `/ca_credentials/` inside of the `/main/gen_ca/` directory. 
Now, publish the `/ca_credentials/ca.crt` for any client to receive the public key with the filename `ca.crt` (for instance as a github gist). 

### To run the server
- Run `go build .` in the `./cmd/main` to build the code.
- Run `./main` in the `./cmd/main` directory
- The server should be running on `localhost:8080`
- Go to `localhost:8080` and verify that it says `404 page not found` to make sure it is running

### To set up a client
- Update the `trustedCertURL` variable to the link of the `ca.crt` (aka public key) of the CA server you wish to use in the `client.go` file.
- Run the `Bootstrap()` method to download the public key of the server; this should download the public in inside of the `/client/root_ca_cert/` directory (comment out all lines in main except for the Bootstrap one)). 
- If you want to sanity check, make sure it matches the public key on the link.
- Now, you can call any method (for example, `createIdentity()`) and pass in the URL of the server to interact with it, but do not call `Bootstrap()` again.

## Session Flow

### Users
- User1
- User2

### Steps
1. User1 wishes to communicate with User2 after verifying User2’s identity.
2. User2 offers User1 its ILP Address (and vic-versa).
3. Using User1’s ILP Address, User2 uses the `RetrieveIdentity()` api call to retrieve the corresponding certificate.
4. User2 now uses the `IsValidCertificate()` api call to verifiy that the retrieved certificate has not expired and that it is infact signed by the CA. If so, then User2 proceeds to communicate with User1.

## Running Unit Tests
Run `go test` in the `./cmd/main` directory
