-- +goose Up
-- +goose StatementBegin
-- 1. Создаем роли
INSERT INTO public.roles (slug, name, description, level, is_system) VALUES
('reader', 'Читатель', 'Только чтение', 2, true),
('user',   'Пользователь',   'Обычный пользователь', 4, true),
('admin',  'Администратор',  'Администратор', 9, true),
('root',   'Root',   'Суперпользователь', 10, true);

-- 2. Выстраиваем иерархию (каждая следующая наследует предыдущую)
-- user -> reader
-- admin -> user
-- root -> admin
INSERT INTO public.role_hierarchy (parent_role_id, role_id)
VALUES 
    ((SELECT id FROM public.roles WHERE slug = 'user'),   (SELECT id FROM public.roles WHERE slug = 'reader')),
    ((SELECT id FROM public.roles WHERE slug = 'admin'),  (SELECT id FROM public.roles WHERE slug = 'user')),
    ((SELECT id FROM public.roles WHERE slug = 'root'),   (SELECT id FROM public.roles WHERE slug = 'admin'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM public.role_hierarchy WHERE parent_role_id IN (SELECT id FROM public.roles WHERE slug IN ('user', 'admin', 'root'));
DELETE FROM public.roles WHERE slug IN ('reader', 'user', 'admin', 'root');
-- +goose StatementEnd
