create table products(
    product_id serial primary key not null,
    name varchar(255) not null unique,
    description text not null,
    added_at timestamptz default now(),
    removed_at timestamptz,
    tags text
);

create table product_variants(
    variant_id serial primary key not null,
    product_id int references products(product_id),
    weight int not null,
    unit varchar(255) not null,
    added_at timestamptz default now(),
    removed_at timestamptz,
    unique (product_id, weight, unit)
);

create table product_prices(
    price_id serial primary key,
    variant_id int references product_variants(variant_id),
    start_date timestamptz not null default now(),
    end_date timestamptz,
    price decimal(10, 2) not null,
    unique (variant_id, start_date)
);

create table storages(
    storage_id serial primary key,
    name varchar(255) not null unique,
    added_at timestamptz not null default now(),
    removed_at timestamptz
);

CREATE TABLE products_in_storage (
    pis_id serial primary key,
    variant_id int references product_variants(variant_id),
    storage_id int references storages(storage_id),
    added_at timestamptz not null default now(),
    removed_at timestamptz,
    quantity int not null
);

create table sales (
    sales_id serial primary key,
    variant_id int references product_variants(variant_id),
    storage_id int references storages(storage_id),
    sold_at timestamptz not null default now(),
    quantity int,
    total_price decimal(10, 2) not null
);

create table log_table (
    log_id serial primary key,
    logtime timestamptz not null,
    flag varchar(255),
    msg text not null,
    module varchar(255),
    fl varchar(255) not null,
    ln varchar(255) not null
);

create table log_details (
    LOG_ID integer not null,
    NAME text not null,
    VALUE text not null
);

INSERT INTO
    products (name, description, tags)
VALUES
    (
        'Порошок Ариэль',
        'Моющий порошок для стирки',
        'стирка,моющее средство'
    ),
    (
        'Вода Hydrolife',
        'Питьевая вода',
        'вода,напиток'
    ),
    (
        'Чай Ahmad',
        'Черный чай',
        'чай,напиток'
    );

INSERT INTO
    product_variants (product_id, weight, unit)
VALUES
    (1, 500, 'г'),
    (1, 3, 'кг'),
    (1, 10, 'кг'),
    (2, 1.5, 'л'),
    (2, 1, 'л'),
    (2, 0.5, 'мл'),
    (3, 100, 'г');

INSERT INTO
    product_prices (variant_id, start_date, end_date, price)
VALUES
    (1, '2023-07-01', '2024-07-25', 5.99),
    (2, '2023-07-01', '2024-07-25', 19.99),
    (3, '2023-07-01', '2024-07-25', 49.99),
    (4, '2023-07-01', '2024-07-25', 1.99),
    (5, '2023-07-01', '2024-07-25', 1.49),
    (6, '2023-07-01', '2024-07-25', 0.99),
    (7, '2023-07-01', '2024-07-25', 2.99);

INSERT INTO
    storages (name, added_at)
VALUES
    ('Склад 1', '2023-07-01'),
    ('Склад 2', '2023-07-01'),
    ('Склад 3', '2023-07-01');

INSERT INTO
    products_in_storage (variant_id, storage_id, added_at, quantity)
VALUES
    (1, 1, '2023-07-01', 2),
    (2, 1, '2023-07-01', 3),
    (3, 1, '2023-07-01', 5),
    (4, 1, '2023-07-01', 6),
    (5, 2, '2023-07-01', 1),
    (6, 2, '2023-07-01', 4),
    (7, 2, '2023-07-01', 2),
    (7, 3, '2023-07-01', 3);

INSERT INTO
    sales (
        variant_id,
        storage_id,
        sold_at,
        quantity,
        total_price
    )
VALUES
    (1, 1, '2023-07-03 10:30:00', 2, 10.89),
    (4, 1, '2023-07-04 14:15:00', 3, 5.97),
    (5, 2, '2023-07-02 09:45:00', 1, 1.49);