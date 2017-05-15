#!/bin/bash

DATA_DIR='data'
function create_user {
	curl -v -XPOST --header "Content-Type: application/json" -d @$DATA_DIR/goodUser.json localhost:8080/users
}

function create_user_no_body {
	curl -v -XPOST --header "Content-Type: application/json" localhost:8080/users
}

