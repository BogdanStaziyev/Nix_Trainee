create table commentses
(
    id      serial primary key,
    post_id integer not null,
    name    varchar not null,
    email   varchar not null,
    body    varchar not null
);

alter table commentses
    owner to postgres;
