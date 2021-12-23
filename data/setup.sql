create table posts(
	id		serial primary key,
	body		text,
	user_id	 	integer,
	created_at	timestamp not null
);
	

