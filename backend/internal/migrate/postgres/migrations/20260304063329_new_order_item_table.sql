-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS public.positions (
    id uuid NOT NULL,
    order_id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    quantity real NOT NULL,
    notes text COLLATE pg_catalog."default" DEFAULT ''::text,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT positions_pkey PRIMARY KEY (id),
    CONSTRAINT orders_positions_id_fkey FOREIGN KEY (order_id)
        REFERENCES public.orders (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.positions
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.positions;

DROP EXTENSION IF EXISTS pg_trgm;
-- +goose StatementEnd
