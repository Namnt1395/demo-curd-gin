
create table curds
(
    id         bigint not null primary key,
    name       varchar(255),
    email      varchar(255),
    phone      varchar(11),
    city       varchar(255),
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
)
