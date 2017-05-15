create table user_info (
	id bigserial primary key,
	username text,
	fullname text,
	passwordhash text,
	passwordsalt text,
	role text,
	isdisabled bool
);
