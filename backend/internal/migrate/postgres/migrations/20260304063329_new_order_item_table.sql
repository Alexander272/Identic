-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE order_items (
    id uuid NOT NULL,
    order_id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    quantity INT NOT NULL,
    search_text text COLLATE pg_catalog."default" DEFAULT ''::text,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT order_items_pkey PRIMARY KEY (id),
    CONSTRAINT orders_order_items_id_fkey FOREIGN KEY (order_items_id)
        REFERENCES public.orders (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.order_items;

DROP EXTENSION IF EXISTS pg_trgm;
-- +goose StatementEnd
