## Example of Go REST microservice

- REST API Server
- DB with ORM
- VueJS frontend
- Tests

#### To run:

    make deploy

Then point your browser to http://localhost:3000/

#### To test:

    make test

## API endpoints

**/api/v1/wallet/change-balance**

    METHOD: POST
    BODY: {
              "walletId": <uuid4>,
              "operationType": <str: deposit|withdraw>,
              "amount": <int>
          }
    RESP: 201
    ERR_RESP: 400, 404

**/api/v1/wallet/<Wallet ID: uuid4>**

    RESP: 200
    ERR_RESP: 400, 404

**/api/v1/wallet**

    METHOD: POST
    BODY: (empty)
    RESP: 201
    RESP_BODY: Wallet object (Id, Amount, CreatedAt, UpdatedAt)
    ERR_RESP: 400

**/api/v1/wallet**

    METHOD: GET
    RESP: 200
    RESP_BODY: Array of Wallet objects (max 10)
    ERR_RESP: 400

**/api/v1/wallet/<Wallet ID: uuid4>**

    METHOD: DELETE
    RESP: 204
    ERR_RESP: 400, 404

## Load testing

    wrk -t4 -c100 -d10s -R1000 http://0.0.0.0:8000/api/v1/wallet/e2266d5a-8804-4d0c-91de-3b7dd18c12c6

    wrk -t4 -c100 -d10s -R1000 -s wrk.lua http://0.0.0.0:8000/api/v1/wallet
