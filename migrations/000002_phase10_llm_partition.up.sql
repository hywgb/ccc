-- Phase 10: LLM Gateway model configs + MySQL partitioning

-- New table for multi-model LLM gateway
CREATE TABLE IF NOT EXISTS llm_model_configs (
  id              BIGINT UNSIGNED PRIMARY KEY,
  tenant_id       BIGINT UNSIGNED NOT NULL,
  name            VARCHAR(128) NOT NULL,
  provider_type   ENUM('tongyi','openai','bailian','self_hosted') NOT NULL,
  endpoint        VARCHAR(512) NOT NULL DEFAULT '',
  api_key         VARCHAR(512) NOT NULL DEFAULT '',
  model_name      VARCHAR(128) NOT NULL,
  is_default      BOOLEAN NOT NULL DEFAULT FALSE,
  is_active       BOOLEAN NOT NULL DEFAULT TRUE,
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_llm_tenant (tenant_id),
  INDEX idx_llm_default (tenant_id, is_default, is_active),
  CONSTRAINT fk_llm_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- MySQL Partitioning for high-volume tables (Phase 10 scaling)
-- Note: Partitioning requires dropping foreign keys first, then re-creating without FK.
-- These ALTER TABLE statements convert to RANGE partitioning by month.
-- Run during a maintenance window on production.

-- calls: partition by created_at month
-- (FK constraint on calls prevents direct ALTER — in production, drop FK first)
-- ALTER TABLE calls PARTITION BY RANGE (UNIX_TIMESTAMP(created_at)) (
--   PARTITION p202501 VALUES LESS THAN (UNIX_TIMESTAMP('2025-02-01')),
--   PARTITION p202502 VALUES LESS THAN (UNIX_TIMESTAMP('2025-03-01')),
--   ...
--   PARTITION pmax VALUES LESS THAN MAXVALUE
-- );

-- Provided as reference DDL. Production partition management should use
-- an automated partition maintenance script (e.g., cron adding monthly partitions).
