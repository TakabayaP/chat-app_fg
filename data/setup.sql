
create table chatrooms (
    id      serial primary key,
    uuid    varchar(64) not null unique,
    name    text 
);

create table posts (
	id		serial primary key,
	body		text,
	user_id	 	integer,
	created_at	timestamp not null
);
	

