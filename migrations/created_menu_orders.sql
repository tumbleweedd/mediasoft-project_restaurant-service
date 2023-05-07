create table if not exists restaurant.menu_product
(
    id       serial primary key,
    menu_id  uuid not null
        constraint fk_product_menu
            references menu (uuid),
    order_id uuid not null
        constraint fk_menu_product
            references product (uuid)
);

drop table if exists restaurant.menu_product;