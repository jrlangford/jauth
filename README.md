# Sessions test

## Multiple application support
Redis is used as key storage backend. It locks any key being used, avoiding concurrency issues.

## Test service

### Get cookie
curl -v -c cookie localhost:8080/cookie/save

### Embed cookie in request
curl -v -b cookie localhost:8080/cookie/read

## References
* Password storage: https://astaxie.gitbooks.io/build-web-application-with-golang/en/09.5.html
* PostgresSQL driver: https://astaxie.gitbooks.io/build-web-application-with-golang/en/05.4.html
* Naming conventions: http://stackoverflow.com/questions/4702728/relational-table-naming-convention/4703155#4703155
* SQL injection: https://astaxie.gitbooks.io/build-web-application-with-golang/en/09.4.html
* Secure password: https://crackstation.net/hashing-security.htm
