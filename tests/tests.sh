#!/bin/bash

DATA_DIR='data'
function create_user {
	curl -v --header "Content-Type: application/json" -d @$DATA_DIR/goodUser.json localhost:8080/user/create
}

function create_user_no_body {
	curl -v --header "Content-Type: application/json" localhost:8080/user/create
}

