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

- client.go: 
  - Populate() - populates the bloom filter with data from MySQL and package it as a JSON, and sends to localhost:8082/populateBF.
  - Clear() - clears the bloom filter by accessing the endpoint localhost:8082/clearBF.
  - Check() - passes userid and list of emails to localhost:8082/checkBF to check if they are in the bloom filter. 
- database.go: creates and inserts randomly generated data into MySQL database. 
- generateEmail.go: generates random userIDs and emails.
- server.go: 
  - populateBF() - retreives JSON data, parses, and populates bloom filter with it.
  - checkBF() - retreives JSON data, parses, and checks the data against the bloom filter and cross checks against MySQL if data is in the bloom filter. 
  - clearBF() - clears the bloom filter. 
- server_test.go: unit testing


## How to test:

  1. Make sure all installation files are present (see "Installation")
  2. Change into directory of server_test.go
  3. In the terminal, type "go test"
