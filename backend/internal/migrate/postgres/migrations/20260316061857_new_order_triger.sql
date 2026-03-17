-- +goose Up
-- +goose StatementBegin
-- 1. Создаем функцию, которая будет отправлять JSON с данными заказа
CREATE OR REPLACE FUNCTION notify_order_changes()
RETURNS TRIGGER AS $$
DECLARE
    payload JSON;
BEGIN
    -- Формируем JSON. Можно отправить весь ряд (new) или только нужные поля.
    -- TG_OP — это операция (INSERT, UPDATE, DELETE)
    payload = json_build_object(
        'id', NEW.id,
        'status', NEW.status,
        'updated_at', NEW.updated_at,
        'action', TG_OP
    );

    -- Отправляем уведомление в канал 'order_updates'
    PERFORM pg_notify('order_updates', payload::text);
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Привязываем функцию к таблице orders
CREATE TRIGGER trg_orders_changed
AFTER INSERT OR UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION notify_order_changes();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_orders_changed ON orders;
DROP FUNCTION IF EXISTS notify_order_changes();
-- +goose StatementEnd
