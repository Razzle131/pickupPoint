-- users --
CREATE TABLE user_roles (
    user_role varchar(20) PRIMARY KEY
);

INSERT INTO user_roles
VALUES ('moderator');

INSERT INTO user_roles
VALUES ('employee');

CREATE TABLE users (
    user_id uuid PRIMARY KEY,
    user_email varchar(20) NOT NULL UNIQUE,
    user_passwd varchar(255) NOT NULL,
    user_role varchar(20) NOT NULL REFERENCES user_roles (user_role) 
);

-- pvz --
CREATE TABLE pvz_cities (
    pvz_city varchar(20) PRIMARY KEY
);

INSERT INTO pvz_cities
VALUES ('Москва');

INSERT INTO pvz_cities
VALUES ('Санкт-Петербург');

INSERT INTO pvz_cities
VALUES ('Казань');

CREATE TABLE pvz (
    pvz_id uuid PRIMARY KEY,
    pvz_city varchar(20) NOT NULL REFERENCES pvz_cities (pvz_city),
    pvz_date timestamptz NOT NULL
);

-- receptions --
CREATE TABLE reception_statuses (
    reception_status varchar(20) PRIMARY KEY
);

INSERT INTO reception_statuses
VALUES ('in_progress');

INSERT INTO reception_statuses
VALUES ('close');

CREATE TABLE receptions (
    reception_id uuid PRIMARY KEY,
    reception_date timestamptz NOT NULL,
    reception_status varchar(20) NOT NULL REFERENCES reception_statuses (reception_status),
    pvz_id uuid NOT NULL REFERENCES pvz (pvz_id)
);

-- products --
CREATE TABLE product_types (
    product_type varchar(20) PRIMARY KEY
);

INSERT INTO product_types
VALUES ('электроника');

INSERT INTO product_types
VALUES ('одежда');

INSERT INTO product_types
VALUES ('обувь');

CREATE TABLE products (
    product_id uuid PRIMARY KEY,
    product_date timestamptz NOT NULL,
    product_type varchar(20) NOT NULL REFERENCES product_types (product_type),
    reception_id uuid NOT NULL REFERENCES receptions (reception_id)
);

