create table core.categories
(
    id          integer      not null,
    name        varchar(200) not null,
    description varchar(500),
    ems_ids     text[]       not null default '{}',

    constraint categories_pk primary key (id)
);
call sys.grant_permissions_to_role('delete', 'core.categories', '{{.coreRWRole}}');

create table core.ems_categories
(
    id   uuid         not null,
    name varchar(200) not null,

    constraint ems_categories_pk primary key (id)
);
call sys.grant_permissions_to_role('delete', 'core.ems_categories', '{{.coreRWRole}}');

create table core.ems_themes
(
    id             uuid        not null default gen_random_uuid(),
    code           varchar(20) not null,
    datasets_count integer     not null default 0,

    constraint ems_themes_pk primary key (id),
    constraint ems_themes_code_unique unique (code)
);
call sys.grant_permissions_to_role('delete', 'core.ems_themes', '{{.coreRWRole}}');

create table core.ems_theme_translations
(
    id           integer      not null,
    ems_theme_id uuid         not null,
    language     varchar(10)  not null,
    value        varchar(200) not null,
    description  varchar(500),
    created_at   timestamptz,
    updated_at   timestamptz,

    constraint ems_theme_translations_pk primary key (id),
    constraint ems_theme_translations_theme_fk foreign key (ems_theme_id) references core.ems_themes (id)
);
call sys.grant_permissions_to_role('delete', 'core.ems_theme_translations', '{{.coreRWRole}}');

create index ems_theme_translations_theme_id_idx on core.ems_theme_translations (ems_theme_id);

create table core.ems_theme_ems_categories
(
    ems_theme_id    uuid not null,
    ems_category_id uuid not null,

    constraint ems_theme_ems_categories_pk primary key (ems_theme_id, ems_category_id),
    constraint ems_theme_ems_categories_theme_fk foreign key (ems_theme_id) references core.ems_themes (id),
    constraint ems_theme_ems_categories_cat_fk foreign key (ems_category_id) references core.ems_categories (id)
);
call sys.grant_permissions_to_role('delete', 'core.ems_theme_ems_categories', '{{.coreRWRole}}');

create index ems_theme_ems_categories_cat_id_idx on core.ems_theme_ems_categories (ems_category_id);

create table core.collection_metadata
(
    id                integer not null default 1,
    last_collected_at timestamptz,

    constraint collection_metadata_pk primary key (id),
    constraint collection_metadata_single_row check (id = 1)
);

insert into core.collection_metadata (id)
values (1);

---- create above / drop below ----

drop table if exists core.ems_theme_ems_categories;
drop table if exists core.ems_theme_translations;
drop table if exists core.ems_themes;
drop table if exists core.ems_categories;
drop table if exists core.categories;
drop table if exists core.collection_metadata;
