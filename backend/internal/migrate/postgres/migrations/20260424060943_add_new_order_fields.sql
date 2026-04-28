-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.orders 
ADD COLUMN IF NOT EXISTS is_bargaining BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS is_budget BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.orders  
DROP COLUMN IF EXISTS is_bargaining,
DROP COLUMN IF EXISTS is_budget;
-- +goose StatementEnd
