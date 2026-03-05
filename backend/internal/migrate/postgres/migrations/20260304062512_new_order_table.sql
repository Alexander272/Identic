-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.orders (
    id uuid NOT NULL,
    customer text COLLATE pg_catalog."default" DEFAULT ''::text,
    consumer text COLLATE pg_catalog."default" DEFAULT ''::text,
    manager text COLLATE pg_catalog."default" DEFAULT ''::text,
    bill text COLLATE pg_catalog."default" DEFAULT ''::text,
    date timestamp with time zone DEFAULT now(),
    notes text COLLATE pg_catalog."default" DEFAULT ''::text,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT orders_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.orders
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.orders;
-- +goose StatementEnd
