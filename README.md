SendGrid
-------------

Informatics 117 : Team Optimize Prime

- Brenda La
- David Pham
- Eric Chou
- Jose Gomez
- Sheila Truong

## Installation
```
$ go get -u github.com/govend/govend                 # make sure govend is installed
$ git clone https://github.com/blics18/SendGrid.git  # grab all the code Sarah pushed
$ govend -v                                          # download all the dependencies in the vendor.yml file
```

## Endpoints

- To add data into the Bloom Filter: localhost:8082/populateBF
- To check if data is in the Bloom Filter: localhost:8082/checkBF
- To clear the Bloom Filter: localhost:8082/clearBF

## Files:

- client.go: 
  - Populate() - grabs data from MySQL, packages it as a JSON, and sends it to localhost:8082/populateBF.
  - Clear() - calls clearBF endpoint (localhost:8082/clearBF) to clear the bloom filter. 
  - Check() - passes userid and list of emails to localhost:8082/checkBF; calls checkBF to check if an input exists in the bloom filter.
- database.go: creates and inserts randomly-generated data into the MySQL database.
- generateEmail.go (helper file): generates random userIDs and emails.
- server.go: 
  - populateBF() - retreives the JSON data, parses it, and populates the bloom filter with it.
  - checkBF() - retreives the JSON data, parses it, and checks the data against the bloom filter. Also cross checks against the MySQL database if the data exists within the bloom filter. 
  - clearBF() - clears the bloom filter. 
- server_test.go: unit testing

## How to test:

  1. Change into directory of server_test.go
  2. In the terminal, type "go test"
