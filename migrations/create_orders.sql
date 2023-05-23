create table if not exists restaurant.orders
(
    uuid  uuid primary key,
    user_uuid   uuid unique not null,
    office_uuid uuid        not null
        constraint fk_offices_orders
            references restaurant.offices_of_delivery (uuid),
    created_at  timestamp   not null
);



drop table restaurant.orders;