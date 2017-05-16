#!/bin/bash

DATA_DIR='data'
function create_user {
	curl -v -XPOST --header "Content-Type: application/json" -d @$DATA_DIR/goodUser.json localhost:8080/users
}

function create_user_no_body {
	curl -v -XPOST --header "Content-Type: application/json" localhost:8080/users
}

function log_in {
	curl -v -b /tmp/jcookie -c /tmp/jcookie --header "Content-Type: application/json" -d @$DATA_DIR/goodLoginData.json localhost:8080/login
}

function create_admin {
	curl -v -XPOST --header "Content-Type: application/json" -d @$DATA_DIR/admin.json localhost:8080/users
}

function log_in_admin {
	curl -v -b /tmp/jcookie -c /tmp/jcookie --header "Content-Type: application/json" -d @$DATA_DIR/adminLogin.json localhost:8080/login
}

function log_in_bad_data {
	curl -v -XPOST --header "Content-Type: application/json" -d @$DATA_DIR/badLoginData.json localhost:8080/login
}

function log_out {
	curl -v -XPOST -b /tmp/jcookie localhost:8080/logout
}

function get_users {
	curl -v -b /tmp/jcookie localhost:8080/admins/users
}

function get_user_by_email {
	curl -v -b /tmp/jcookie localhost:8080/admins/users/jrobin@gmail.com
}
