CREATE ROLE todo_manager WITH LOGIN PASSWORD '1488';

GRANT USAGE ON SCHEMA public TO todo_manager;
GRANT CREATE ON SCHEMA public TO todo_manager;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO todo_manager;