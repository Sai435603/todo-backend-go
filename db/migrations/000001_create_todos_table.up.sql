create table todos (
    id bigserial primary key,
    title varchar(255) not null,
    description text,
    completed boolean default false,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);
    