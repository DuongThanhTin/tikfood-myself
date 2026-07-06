-- Authentication (ADR-0007), milestone M1. Additive only; no discovery table is altered.
-- citext gives case-insensitive email uniqueness; pgcrypto (gen_random_uuid) is enabled in 001.
create extension if not exists citext;

create table if not exists users (
  id             uuid primary key default gen_random_uuid(),
  email          citext not null unique,
  password_hash  text,                       -- null for OAuth-only accounts
  display_name   text not null default '',
  google_sub     text,                       -- Google subject id when linked; null otherwise
  email_verified boolean not null default false,
  created_at     timestamptz not null default now(),
  updated_at     timestamptz not null default now()
);

-- Partial unique index: many password users have no google_sub (null), which must not collide.
create unique index if not exists users_google_sub_key on users (google_sub) where google_sub is not null;
