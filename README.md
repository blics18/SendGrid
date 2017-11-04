SendGrid
-------------

Informatics 117 : Team Optimize Prime

- Brenda La
- David Pham
- Eric Chou
- Jose Gomez
- Sheila Truong

## Installation
- go get -u github.com/willf/bloom --> Get Bloom Filter files
- go get github.com/stretchr/testify --> Get go's testing package in order to run tests on Bloom Filter
- go get -u github.com/go-sql-driver/mysql --> Get go's mySQL driver

## Endpoints

- To add data into the Bloom Filter: localhost:8082/populateBF
- To check if data is in the Bloom Filter: localhost:8082/checkBF
- To clear the Bloom Filter: localhost:8082/clearBF

## Files:

- populate.go: Generates random user ID and emails, packages as JSON, and sends to localhost:8082/populateBF.
- clear.go: Clears the Bloom Filter by accessing the endpoint localhost:8082/clearBF.
- check.go: Passes userid and list of emails to localhost:8082/checkBF to check if they are in the Bloom Filter. 
- server.go: Retreives JSON data, parses, and populates BF with it.
- server_test.go: Unit testing


## How to test:

  1. Make sure all installation files are present (see "Installation")
  2. Change into directory of server_test.go
  3. In the terminal, type "go test"
