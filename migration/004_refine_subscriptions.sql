-- Migration: Refine subscriptions table
-- Version: 004
-- Description: Add 'updated' status to subscription_status and remove UNIQUE constraint on owner_id

-- 1. Add 'updated' to subscription_status enum
-- Note: In PostgreSQL, you cannot easily add a value to an enum inside a transaction in some versions,
-- but for simple migrations it usually works or can be done with ALTER TYPE.
ALTER TYPE subscription_status ADD VALUE IF NOT EXISTS 'updated';

-- 2. Remove UNIQUE constraint on subscriptions.owner_id
-- First, find the constraint name. It's usually 'subscriptions_owner_id_key' or similar.
-- Based on the schema, it was defined as: owner_id UUID NOT NULL UNIQUE REFERENCES users(id)
ALTER TABLE subscriptions DROP CONSTRAINT IF EXISTS subscriptions_owner_id_key;

-- If it was created as a unique index instead of a constraint:
DROP INDEX IF EXISTS idx_subscriptions_owner_id_unique;
