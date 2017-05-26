# jAuth

A microservice designed to manage users and sessions.

## Goals
* Allow user sign-up, sign-in
* Allow authorized services to query user data
* Allow authorized services to query session data

## Features
* User sign-up/sign-in
* Salted password storage
* Automated db schema generation
* Session data storage in Redis

## Notes
* Concurrent access to session data is hanled internally by Redis

## TODO
* Add unit tests
* Load settings through configuration files
* Return user data on log in
* Add endpoint that returns user data if user is logged in
* Add endpoint to set session data
