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

- ```healthBF```    check the status of the server
- ```populateBF```  add data into the Bloom Filter
- ```checkBF```     check if data exists in the Bloom Filter
- ```clearBF```     clear the Bloom Filter

## Curl command examples

- healthBF   ```curl localhost:1234/healthBF```
- populateBF ```curl localhost:1234/populateBF -d "[{\"UserID\": 123, \"Email\": [\"test@gmail.com\"]}]"```
- checkBF    ```curl localhost:1234/checkBF -d "{\"UserID\": 123, \"Email\": [\"test@gmail.com\"]}"```
- clearBF    ```curl localhost:1234/clearBF```

## Environment Variables

- ```BLOOM_SIZE``` size of the Bloom Filter
- ```BLOOM_PORT``` port that the server will run on
- ```BLOOM_NUM_USERS``` number of users that will be generated
- ```BLOOM_NUM_EMAILS``` up to this number of emails that each user will have
- ```BLOOM_NUM_HASH_FUNCTIONS``` number of hash functions
- ```BLOOM_NUM_TABLES``` is NOT configurable, because createTables.sql is hard coded to create only 5 tables. 

## Files & Description:

- ```client.go``` 
  - Populate() - grabs data from MySQL, packages it as a JSON, and sends it to populateBF endpoint
  - Clear() - calls clearBF endpoint to clear the bloom filter
  - Check() - passes userid and list of emails to checkBF endpoint; calls checkBF to check if an input exists in the bloom   filter
  
- ```runClient.go``` the client "main", and what we used to demo. Calls on all relevant endpoints (health, populate, check) 
  
- ```server.go``` 
  - healthBF() - checks to see if the server is up and running
  - populateBF() - retreives the JSON data, parses it, and populates the bloom filter with it
  - checkBF() - retreives the JSON data, parses it, and checks the data against the bloom filter. Also cross checks against the MySQL database if the data exists within the bloom filter
  - clearBF() - clears the bloom filter
  
- ```createTables.sql``` SQL file that creates the database and tables

- ```database.go``` creates and inserts randomly-generated data into the MySQL database

- ```server_test.go``` unit testing
  
- ```benchmark_test.go``` benchmarking for each endpoint except healthBF
  
- ```generateEmail.go``` generates random userIDs and emails

- ```data.txt``` file with randomly generated data (filled with suppression data and NON-suppresion data)

## How to do unit testing:

  1. Change into directory of ```server_test.go```
  2. In the terminal, type ```go test```

## How to do benchmarking:

  1. Change into directory of ```benchmark_test.go```
  2. In the terminal, type ```go test -run=XXX -bench=name_of_function()```
  
## Graphite Visualization

  1. Make sure that Docker is downloaded and running
  2. In the terminal, type ```docker run -it -d -e MYSQL_ALLOW_EMPTY_PASSWORD=1 -p 3306:3306 percona```
  3. Set up Graphite in docker:
         ```docker run -d\
         --name graphite\
         --restart=always\
         -p 80:80\
         -p 2003-2004:2003-2004\
         -p 2023-2024:2023-2024\
         -p 8125:8125/udp\```
  7. Run server.go, and test using curl commands
  8. Go to ```localhost``` in any browser to check and configure the metrics
  
## Setting up MySQL

1. Create the docker container: ```docker run -v $(pwd)/createTables.sql:/root/createTables.sql --name bfmysql -it -d -e MYSQL_ALLOW_EMPTY_PASSWORD=1 -p 3306:3306 percona```
2. Run this command on the container: ```docker exec -it bfmysql /bin/bash```
3. Change directories to where the .sql file is stored: ```cd /root```
4. Run the file that creates the schema with ```mysql < createTables.sql```
