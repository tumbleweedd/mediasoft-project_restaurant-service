create table if not exists restaurant.offices_of_delivery
(
    uuid    uuid primary key,
    name    varchar(50)  not null,
    address varchar(100) not null
);

CREATE OR REPLACE FUNCTION check_office_exists() RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM restaurant.offices_of_delivery WHERE uuid = NEW.uuid) THEN
        -- Вставляем запись об офисе в таблицу офисов
        INSERT INTO restaurant.offices_of_delivery (uuid, name, address)
        VALUES (NEW.uuid, NEW.name, NEW.address);
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_office_trigger
    AFTER INSERT
    ON restaurant.offices_of_delivery
    FOR EACH ROW
EXECUTE FUNCTION check_office_exists();

drop table restaurant.offices_of_delivery;