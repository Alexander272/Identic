-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.role_hierarchy (
    parent_role_id uuid NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    role_id  uuid NOT NULL REFERENCES roles(id) ON DELETE CASCADE,

    PRIMARY KEY (parent_role_id, role_id),

    CHECK (parent_role_id <> role_id)
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.role_hierarchy
    OWNER to postgres;

-- === ЗАЩИТА ОТ ЦИКЛИЧЕСКОГО НАСЛЕДОВАНИЯ (триггер) ===
CREATE OR REPLACE FUNCTION check_role_hierarchy_cycle()
RETURNS TRIGGER AS $$
BEGIN
    -- Ищем: нет ли среди предков нашего нового ПАПЫ нашего же РЕБЕНКА?
    -- (Если ребенок уже является предком своего будущего папы — это цикл)
    IF EXISTS (
        WITH RECURSIVE parents AS (
            -- Начинаем от того, КТО станет родителем в новой записи
            SELECT parent_role_id 
            FROM role_hierarchy 
            WHERE role_id = NEW.parent_role_id 
            
            UNION ALL
            
            SELECT rh.parent_role_id
            FROM role_hierarchy rh
            JOIN parents p ON rh.role_id = p.parent_role_id
        )
        SELECT 1 FROM parents WHERE parent_role_id = NEW.role_id
    ) THEN
        RAISE EXCEPTION 'ERR_CIRCULAR: Role % is already a parent in the chain for %', 
            NEW.role_id, NEW.parent_role_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Удаляем старый триггер (если он был), чтобы не было конфликта
DROP TRIGGER IF EXISTS trg_role_hierarchy_cycle ON public.role_hierarchy;
-- Создаем триггер
CREATE TRIGGER trg_role_hierarchy_cycle
    BEFORE INSERT OR UPDATE ON role_hierarchy
    FOR EACH ROW EXECUTE FUNCTION check_role_hierarchy_cycle();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_role_hierarchy_cycle ON role_hierarchy;
DROP FUNCTION IF EXISTS check_role_hierarchy_cycle;

DROP TABLE IF EXISTS public.role_hierarchy;
-- +goose StatementEnd
