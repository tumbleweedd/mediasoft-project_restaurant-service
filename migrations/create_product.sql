create type restaurant.ProductType as enum (
    'PRODUCT_TYPE_UNSPECIFIED',
    'PRODUCT_TYPE_SALAD',
    'PRODUCT_TYPE_GARNISH',
    'PRODUCT_TYPE_MEAT',
    'PRODUCT_TYPE_SOUP',
    'PRODUCT_TYPE_DRINK',
    'PRODUCT_TYPE_DESSERT'
    );


drop table if exists restaurant.product;

create table if not exists restaurant.product
(
    uuid uuid primary key,
    name         varchar(50) not null,
    description  varchar(255),
    type         restaurant.ProductType default 'PRODUCT_TYPE_UNSPECIFIED'::restaurant.ProductType,
    weight       int,
    price        double precision,
    created_at   timestamp              default now()
)