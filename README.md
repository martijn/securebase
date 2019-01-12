Securebase
==========

Securebase is an encrypted key-value store accessible through a HTTP API. Each
stored value is encrypted with a server-side secret (in the file `keyfile`) and
a client-supplied secret.

The use case is taking sensitive data out of the SQL database supporting your
web app. The client secret can be stored in that database and be used to
retrieve data from Securebase when required. This way the most sensitive data
in your application is moved out of the app server / database server attack
surface; sensitive data can no longer be leaked through database-related
attacks or user errors. Securebase should be hosted on an isolated server that
is not accessible from the outside world.

Values are stored encrypted with

    AES-GCM(HKDF(server-key | client-key))

The server key is used for every value and every value is encrypted with a
unique nonce plus the client-key which is unknown to an attacker attacking
the Securebase server.

Setting up
==========

First, generate a server-secret in the keyfile. For example, fill it with some
random data:

    openssl rand -base64 32

Then run the binary. It will create an empty datastore in the `datastore` file
and listen for HTTP requests on port 5800.

API usage
=========

The HTTP API supports GET, POST and DELETE. The path is the key name. The
client secret is supplied in the Client-Secret header.

The key-value store can store anything from small strings to large binaries.

Example
=======

Retrieving empty key 'testkey':

    $ curl -i -H "Client-secret: testsecret" -X GET http://localhost:5800/testkey
    HTTP/1.1 404 Not Found

    Key not found

Setting a value with secret 'testsecret':
 
    $ curl -i -H "Client-secret: testsecret" -X POST --data-binary "test" http://localhost:5800/testkey
    HTTP/1.1 200 OK

Retrieving the value:

    $ curl -i -H "Client-secret: testsecret" -X GET http://localhost:5800/testkey
    HTTP/1.1 200 OK

    test

Trying to retrieve the value with the wrong client secret:

    $ curl -i -H "Client-secret: wrongsecret" -X GET http://localhost:5800/testkey
    HTTP/1.1 401 Unauthorized

    Error in decryption. Incorrect Client-secret?

Deletion does not require Client-secret:

    $ curl -i -H "Client-secret: wrongsecret" -X DELETE http://localhost:5800/testkey
    HTTP/1.1 200 OK

TODO
====

Securebase is not production quality unless the following is implemented:

* Code cleanup, applying Go best practices

* Decent error handling

* Tests

* User authorisation

* HTTPS

* Rate limiting to prevent abuse

Furthermore it could use:

* An external backend datastore to easily make high availability scenarios possible

* Storing the server-key using asymmetric encryption, opening the possibility
  for using a HSM.

