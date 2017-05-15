create table p_user (
	 id bigserial primary key,
	 username text,
	 fullname text,
	 passwordhash text,
	 passwordsalt text,
	 isdisabled bool
);
