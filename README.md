# Sessions test

## Get cookie
curl -v -c cookie localhost:8080/cookie/save

## Embed cookie in request
curl -v -b cookie localhost:8080/cookie/read
