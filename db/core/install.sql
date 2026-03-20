-----------------------
-- CONFIGURE SCHEMAS --
-----------------------
create schema if not exists sys;
grant usage on schema sys to public;
create schema if not exists core;
grant usage on schema core to public;
----------------------
-- HELPER FUNCTIONS --
----------------------
create or replace procedure pg_temp.create_user(name text, password text) as
$$
declare
    current_database_name text;
begin
    select current_database from current_database() into current_database_name;
    begin
        execute format('create user %s with password ''%s''', name, password);
    exception
        when duplicate_object then null;
        when others then raise;
    end;
    execute format('revoke all privileges on database %s from %s', current_database_name, name);
    execute format('grant connect on database %s to %s', current_database_name, name);
end;
$$ language plpgsql;

create or replace procedure pg_temp.create_role(name text) as
$$
begin
    execute format('create role %s with nologin', name);
exception
    when duplicate_object then null;
    when others then raise;
end;
$$ language plpgsql;

create or replace procedure pg_temp.grant_role_to_user(role text, name text) as
$$
begin
    execute format('grant %s to %s', role, name);
end;
$$ language plpgsql;

create or replace procedure pg_temp.grant_schema_permissions_to_role(schema_name text, permissions text, role text) as
$$
begin
    execute format('alter default privileges for role current_user in schema %s grant %s to group %s', schema_name, permissions, role);
end;
$$ language plpgsql;

create or replace procedure sys.grant_permissions_to_role(permissions text, schema_table text, role text) as
$$
begin
    execute format('grant %s on %s to %s', permissions, schema_table, role);
end;
$$ language plpgsql;
------------------
-- CREATE USERS --
------------------
call pg_temp.create_user('{{.apiDatasourceUser}}', '{{.apiDatasourcePassword}}');
------------------
-- CREATE ROLES --
------------------
call pg_temp.create_role('{{.sysRWRole}}');
call pg_temp.create_role('{{.sysRRole}}');
call pg_temp.create_role('{{.coreRWRole}}');
call pg_temp.create_role('{{.coreRRole}}');
-----------------
-- GRANT ROLES --
-----------------
call pg_temp.grant_role_to_user('{{.coreRWRole}}', '{{.apiDatasourceUser}}');

-- create test user if necessary
do
$$
    begin
        if {{.createDBTestUser}} then
            call pg_temp.create_user('{{.testDatasourceUser}}', '{{.testDatasourcePassword}}');
            call pg_temp.grant_role_to_user('{{.coreRWRole}}', '{{.testDatasourceUser}}');
        end if;
    end
$$ language plpgsql;

----------------
-- PRIVILEGES --
----------------
call pg_temp.grant_schema_permissions_to_role('sys', 'select, insert, update on tables', '{{.sysRWRole}}');
call pg_temp.grant_schema_permissions_to_role('sys', 'usage on sequences', '{{.sysRWRole}}');
call pg_temp.grant_schema_permissions_to_role('sys', 'select on tables', '{{.sysRRole}}');
call pg_temp.grant_schema_permissions_to_role('core', 'select, insert, update on tables', '{{.coreRWRole}}');
call pg_temp.grant_schema_permissions_to_role('core', 'usage on sequences', '{{.coreRWRole}}');
call pg_temp.grant_schema_permissions_to_role('core', 'select on tables', '{{.coreRRole}}');
----------------
-- EXTENSIONS --
----------------
create extension if not exists "citext" schema sys;
---------------------------
-- CONFIGURE SEARCH PATH --
---------------------------
do
$$
    declare
        current_database_name text;
    begin
        select current_database from current_database() into current_database_name;
        execute format('alter database %s set search_path to sys', current_database_name);
    end;
$$ language plpgsql;