create extension if not exists pgcrypto;
create schema if not exists meta;
create schema if not exists data;

create table if not exists meta.datasets (
    id uuid primary key default gen_random_uuid(),
    code text not null unique,
    name text not null,
    description text not null default '',
    table_name text not null unique,
    mqtt_topic text,
    is_hidden boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists meta.dataset_columns (
    id uuid primary key default gen_random_uuid(),
    dataset_id uuid not null references meta.datasets(id) on delete cascade,
    code text not null,
    name text not null,
    data_type text not null,
    ordinal int not null,
    nullable boolean not null default true,
    unique(dataset_id, code)
);

create table if not exists meta.algorithms (
    id uuid primary key default gen_random_uuid(),
    code text not null unique,
    name text not null,
    description text not null default '',
    run_url text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists meta.transformations (
    id uuid primary key default gen_random_uuid(),
    code text not null unique,
    name text not null,
    description text not null default '',
    algorithm_id uuid not null references meta.algorithms(id),
    target_dataset_id uuid not null references meta.datasets(id),
    period_seconds int,
    enabled boolean not null default false,
    last_processed_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table if not exists meta.transformation_sources (
    transformation_id uuid not null references meta.transformations(id) on delete cascade,
    dataset_id uuid not null references meta.datasets(id),
    primary key (transformation_id, dataset_id)
);

create table if not exists meta.transformation_runs (
    id uuid primary key default gen_random_uuid(),
    transformation_id uuid not null references meta.transformations(id),
    status text not null,
    started_at timestamptz not null default now(),
    finished_at timestamptz,
    error_message text,
    request_body text,
    response_body text
);
