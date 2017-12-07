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
$ git clone https://github.com/blics18/SendGrid.git  # grab all the code from github
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
  
## How to get Graphite Visualization
  1. Make sure that docker is downloaded and running. (Docker for windows if you are a Windows user)
  2. Use command docker run -it -d -e MYSQL_ALLOW_EMPTY_PASSWORD=1 -p 3306:3306 percona
  3. Create the metrics registry inside of CreatebloomFilter func in Server (Default code on rcrowley, keep the default port (2003)
  4. Set up Graphite in docker:
        docker run -d\
         --name graphite\
         --restart=always\
         -p 80:80\
         -p 2003-2004:2003-2004\
         -p 2023-2024:2023-2024\
         -p 8125:8125/udp\
         Visit https://github.com/hopsoft/docker-graphite-statsd for more info
  5. docker ps -a shows what you have to make sure you have Graphite up and running
  6. Add code to whereever you want to store metrics.
  7. Start server and test using curl commands. Ex: curl localhost:8082/populateBF -d "[{\"UserID\": 123, \"Email\": [\"test@gmail.com\"]}]" populates the Bloom Filter with user JSON.
  8. Go to localhost in browser to check and configure your metrics.



