#!/bin/bash

DATA_DIR='data'
function create_user {
	curl -v -XPOST --header "Content-Type: application/json" -d @$DATA_DIR/goodUser.json localhost:8080/users
}

function create_user_no_body {
	curl -v -XPOST --header "Content-Type: application/json" localhost:8080/users
}

function log_in {
	curl -v -c mycookie --header "Content-Type: application/json" -d @$DATA_DIR/goodLoginData.json localhost:8080/login
}

function log_in_bad_data {
	curl -v -XPOST --header "Content-Type: application/json" -d @$DATA_DIR/badLoginData.json localhost:8080/login
}

function log_out {
	curl -v -XPOST -b mycookie localhost:8080/logout
}

function get_user_data {
	curl -v -b mycookie localhost:8080/users/email/jrobin@gmail.com
}
