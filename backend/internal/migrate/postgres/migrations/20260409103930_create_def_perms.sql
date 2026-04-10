-- +goose Up
-- +goose StatementBegin
-- 1. Создаем все возможные разрешения для ресурсов
INSERT INTO public.permissions (object, action, description) VALUES
-- Заявки (ResourceOrder)
('order', 'read', 'Просмотр заявок'),
('order', 'write', 'Создание и редактирование заявок'),
('order', 'delete', 'Удаление заявок'),
-- Пользователи (ResourceUser)
('user', 'read', 'Просмотр пользователей'),
('user', 'write', 'Редактирование пользователей'),
('user', 'delete', 'Удаление пользователей'),
-- Роли (ResourceRole)
('role', 'read', 'Просмотр ролей'),
('role', 'write', 'Редактирование ролей'),
('role', 'delete', 'Удаление ролей'),
-- Права (ResourcePerm)
('permission', 'read', 'Просмотр прав'),
('permission', 'write', 'Редактирование прав'),
('permission', 'delete', 'Удаление прав');

-- 2. Распределяем права по ролям
-- Reader: только чтение заявок
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.slug = 'reader' AND p.object = 'order' AND p.action = 'read';

-- User: полный доступ к заявкам (read, write, delete)
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.slug = 'user' AND p.object = 'order';

-- Admin: доступ к администрированию (users, roles, permissions) + заявки
INSERT INTO public.role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM public.roles r, public.permissions p
WHERE r.slug = 'admin' 
AND p.object IN ('user', 'role', 'permission');

-- Root: обычно получает всё (можно продублировать админа или дать расширенные системные права)
-- INSERT INTO public.role_permissions (role_id, permission_id)
-- SELECT r.id, p.id FROM public.roles r, public.permissions p
-- WHERE r.slug = 'root';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM public.role_permissions;
DELETE FROM public.permissions WHERE object IN ('order', 'user', 'role', 'permission');
-- +goose StatementEnd
