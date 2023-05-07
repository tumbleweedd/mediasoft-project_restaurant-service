drop table if exists restaurant.menu;

create table if not exists restaurant.menu
(
    uuid         uuid primary key,
    on_date           timestamp not null,
    opening_record_at timestamp not null,
    closing_record_at timestamp not null,
    created_at        timestamp default now()
)