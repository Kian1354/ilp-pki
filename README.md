# ilp-pki

Set up Redis 
`docker pull redis`  
`docker run --name redis-test-instance -p 6379:6379 -d redis`

Get go-redis package
`go get github.com/go-redis/redis`

## Session Flow

### Users
- User1
- User2

### Steps
1. User1 wishes to communicate with User2 after verifying User2’s identity.
2. User2 offers User1 its ILP Address (and vic-versa).
3. Using User1’s ILP Address, User2 uses the `RetrieveIdentity()` api call to retrieve the corresponding certificate.
4. User2 now uses the `IsValidCertificate()` api call to verifiy that the retrieved certificate has not expired and that it is infact signed by the CA. If so, then User2 proceeds to communicate with User1.

## Running Sample Program Flow
Run `go run pki.go` in the `./cmd/main` directory

## Running Unit Tests
Run `go test` in the `./cmd/main` directory

## Generate CA keys
Perform this command first to set up CA credentials on CA server:  
Run `go gen_ca.go` in the `./cmd/gen_ca` directory
