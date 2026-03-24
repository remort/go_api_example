## Example of Go REST microservice

- REST API Server
- DB with ORM
- Tests

#### To run:

    make deploy

#### To test:

    make test

## API endpoints

**/api/v1/wallet**

    METHOD: POST
    BODY: {
              "valletId": <uuid4>,
              "operationType": <str: deposit|withdraw>,
              "amount": <int>
          }
    RESP: 201
    ERR_RESP: 400, 404

**/api/v1/wallet/<Wallet ID: uuid4>**

    RESP: 200
    ERR_RESP: 400, 404


## Load testing

    wrk -t4 -c100 -d10s -R1000 http://0.0.0.0:8000/api/v1/wallet/e2266d5a-8804-4d0c-91de-3b7dd18c12c6

    wrk -t4 -c100 -d10s -R1000 -s wrk.lua http://0.0.0.0:8000/api/v1/wallet
