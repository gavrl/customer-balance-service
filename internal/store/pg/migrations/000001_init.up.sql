create table balance
(
    id          serial not null unique,
    customer_id int    not null unique,
    balance     bigint not null
);