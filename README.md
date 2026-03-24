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
