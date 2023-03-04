drop table if exists categories;
drop table if exists deliveries;
DROP TABLE IF EXISTS public.delivery_items;
DROP TABLE IF EXISTS public.images;
DROP TABLE IF EXISTS public.items;
DROP TABLE IF EXISTS public.products;
DROP TABLE IF EXISTS public.sessions;
DROP TABLE IF EXISTS public.suppliers;
DROP TABLE IF EXISTS public.userdata;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.product_infos;


DROP SEQUENCE IF EXISTS public.categories_id_seq;

CREATE SEQUENCE IF NOT EXISTS public.categories_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1
    OWNED BY categories.id;

ALTER SEQUENCE public.categories_id_seq
    OWNER TO postgres;

DROP SEQUENCE IF EXISTS public.deliveries_id_seq;

CREATE SEQUENCE IF NOT EXISTS public.deliveries_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1
    OWNED BY deliveries.id;

ALTER SEQUENCE public.deliveries_id_seq
    OWNER TO postgres;

DROP SEQUENCE IF EXISTS public.images_id_seq;

CREATE SEQUENCE IF NOT EXISTS public.images_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1
    OWNED BY images.id;

ALTER SEQUENCE public.images_id_seq
    OWNER TO postgres;

DROP SEQUENCE IF EXISTS public.items_id_seq;

CREATE SEQUENCE IF NOT EXISTS public.items_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1
    OWNED BY items.id;

ALTER SEQUENCE public.items_id_seq
    OWNER TO postgres;

DROP SEQUENCE IF EXISTS public.product_infos_id_seq;

CREATE SEQUENCE IF NOT EXISTS public.product_infos_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1
    OWNED BY product_infos.id;

ALTER SEQUENCE public.product_infos_id_seq
    OWNER TO postgres;

DROP SEQUENCE IF EXISTS public.products_id_seq;

CREATE SEQUENCE IF NOT EXISTS public.products_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1
    OWNED BY products.id;

ALTER SEQUENCE public.products_id_seq
    OWNER TO postgres;

DROP SEQUENCE IF EXISTS public.sessions_id_seq;

CREATE SEQUENCE IF NOT EXISTS public.sessions_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1
    OWNED BY sessions.id;

ALTER SEQUENCE public.sessions_id_seq
    OWNER TO postgres;

DROP SEQUENCE IF EXISTS public.users_id_seq;

CREATE SEQUENCE IF NOT EXISTS public.users_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1
    OWNED BY users.id;

ALTER SEQUENCE public.users_id_seq
    OWNER TO postgres;







CREATE TABLE IF NOT EXISTS public.users
(
    id bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    login text COLLATE pg_catalog."default" NOT NULL,
    username text COLLATE pg_catalog."default" NOT NULL,
    password text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT "unique" UNIQUE (login, username)
    )

    TABLESPACE pg_default;




CREATE TABLE IF NOT EXISTS public.categories
(
    id bigint NOT NULL DEFAULT nextval('categories_id_seq'::regclass),
    name text COLLATE pg_catalog."default",
    CONSTRAINT categories_pkey PRIMARY KEY (id)
    )

    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS public.items
(
    id bigint NOT NULL DEFAULT nextval('items_id_seq'::regclass),
    category_id bigint NOT NULL,
    color text COLLATE pg_catalog."default",
    size text COLLATE pg_catalog."default",
    quantity bigint,
    CONSTRAINT items_pkey PRIMARY KEY (id),
    CONSTRAINT category_id FOREIGN KEY (category_id)
    REFERENCES public.categories (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID
    )

    TABLESPACE pg_default;


CREATE TABLE IF NOT EXISTS public.deliveries
(
    id bigint NOT NULL DEFAULT nextval('deliveries_id_seq'::regclass),
    address text COLLATE pg_catalog."default" NOT NULL,
    phone_number text COLLATE pg_catalog."default" NOT NULL,
    status text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    CONSTRAINT deliveries_pkey PRIMARY KEY (id)
    )

    TABLESPACE pg_default;




CREATE TABLE IF NOT EXISTS public.delivery_items
(
    delivery_id bigint NOT NULL,
    item_id bigint NOT NULL,
    quantity bigint NOT NULL,
    CONSTRAINT delivery_id FOREIGN KEY (delivery_id)
    REFERENCES public.deliveries (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID,
    CONSTRAINT item_id FOREIGN KEY (item_id)
    REFERENCES public.items (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID
    )

    TABLESPACE pg_default;



CREATE TABLE IF NOT EXISTS public.images
(
    id bigint NOT NULL DEFAULT nextval('images_id_seq'::regclass),
    item_id bigint NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    data bytea NOT NULL,
    CONSTRAINT images_pkey PRIMARY KEY (id)
    )

    TABLESPACE pg_default;

CREATE TABLE IF NOT EXISTS public.products
(
    id bigint NOT NULL DEFAULT nextval('products_id_seq'::regclass),
    name text COLLATE pg_catalog."default" NOT NULL,
    category_id bigint NOT NULL,
    price numeric NOT NULL,
    created_at timestamp with time zone,
    CONSTRAINT products_pkey PRIMARY KEY (id),
    CONSTRAINT category_id FOREIGN KEY (category_id)
    REFERENCES public.categories (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID
    )

    TABLESPACE pg_default;



CREATE TABLE IF NOT EXISTS public.product_infos
(
    id bigint NOT NULL DEFAULT nextval('product_infos_id_seq'::regclass),
    product_id bigint NOT NULL,
    created_at timestamp with time zone,
    CONSTRAINT product_infos_pkey PRIMARY KEY (id),
    CONSTRAINT product_id FOREIGN KEY (product_id)
    REFERENCES public.products (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.product_infos
    OWNER to postgres;






CREATE TABLE IF NOT EXISTS public.sessions
(
    id bigint NOT NULL DEFAULT nextval('sessions_id_seq'::regclass),
    user_id bigint NOT NULL,
    uuid text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone,
    deleted_at timestamp with time zone,
    last_login timestamp with time zone,
    ip bytea,
    CONSTRAINT sessions_pkey PRIMARY KEY (id),
    CONSTRAINT user_id FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.sessions
    OWNER to postgres;




CREATE TABLE IF NOT EXISTS public.suppliers
(
    user_id bigint NOT NULL,
    product_id bigint NOT NULL,
    created_at timestamp with time zone,
    CONSTRAINT product_id FOREIGN KEY (product_id)
    REFERENCES public.products (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID,
    CONSTRAINT user_id FOREIGN KEY (user_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID
    )

    TABLESPACE pg_default;



CREATE TABLE IF NOT EXISTS public.userdata
(
    user_id bigint NOT NULL,
    firstname text COLLATE pg_catalog."default" NOT NULL,
    lastname text COLLATE pg_catalog."default",
    phone_number text COLLATE pg_catalog."default" NOT NULL,
    country_code text COLLATE pg_catalog."default" NOT NULL,
    zip_code text COLLATE pg_catalog."default",
    birthday timestamp with time zone,
    CONSTRAINT user_id FOREIGN KEY (user_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.suppliers
    OWNER to postgres;
ALTER TABLE IF EXISTS public.sessions
    OWNER to postgres;
ALTER TABLE IF EXISTS public.userdata
    OWNER to postgres;
ALTER TABLE IF EXISTS public.products
    OWNER to postgres;
ALTER TABLE IF EXISTS public.items
    OWNER to postgres;
ALTER TABLE IF EXISTS public.images
    OWNER to postgres;
ALTER TABLE IF EXISTS public.delivery_items
    OWNER to postgres;
ALTER TABLE IF EXISTS public.deliveries
    OWNER to postgres;
ALTER TABLE IF EXISTS public.categories
    OWNER to postgres;
ALTER TABLE IF EXISTS public.users
    OWNER to postgres;
