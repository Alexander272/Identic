-- +goose Up
-- +goose StatementBegin
UPDATE public.orders 
SET 
    is_bargaining = CASE 
        WHEN notes ILIKE '%лот%' THEN TRUE 
        ELSE is_bargaining 
    END,
    is_budget = CASE 
        WHEN notes ILIKE '%бюджет%' THEN TRUE 
        ELSE is_budget 
    END
WHERE (notes ILIKE '%лот%' OR notes ILIKE '%бюджет%') AND consumer != '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE your_table_name
SET is_bargaining = FALSE, is_budget = FALSE
WHERE (notes ILIKE '%лот%' OR notes ILIKE '%бюджет%') AND consumer != '';
-- +goose StatementEnd
