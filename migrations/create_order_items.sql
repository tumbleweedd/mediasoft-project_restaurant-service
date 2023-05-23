create table if not exists restaurant.order_items
(
    id           serial primary key,
    count        int  not null,
    product_uuid uuid not null
        constraint fk_product_orderItem
            references restaurant.product (uuid),
    order_uuid   uuid not null
        constraint fk_order_order_item
            references restaurant.orders (uuid)
);

drop table restaurant.order_items