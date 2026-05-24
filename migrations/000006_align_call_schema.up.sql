-- 000006: Align `calls` table schema with Go ORM expectations.
-- Round 6 P0-1 fix: the original schema columns (start_at/answer_at/end_at/cli/hangup_cause
-- and missing direction/status/channel_uuid/hold_count/transfer_count/...) did not match the
-- columns the Go entity and repositories actually read/write, causing INSERT/UPDATE to fail
-- at runtime.

-- 1. Rename timestamps and related columns to match Go db tags.
ALTER TABLE calls CHANGE COLUMN start_at     started_at    TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE calls CHANGE COLUMN answer_at    answered_at   TIMESTAMP(3) NULL;
ALTER TABLE calls CHANGE COLUMN end_at       ended_at      TIMESTAMP(3) NULL;
ALTER TABLE calls CHANGE COLUMN cli          caller        VARCHAR(32);
ALTER TABLE calls CHANGE COLUMN hangup_cause hangup_reason VARCHAR(64);

-- 2. Drop VIRTUAL duration_sec (we will store it explicitly so post-call workers can write it).
ALTER TABLE calls DROP COLUMN duration_sec;

-- 3. Lowercase media_type ENUM so it matches Go MediaTypeAudio/MediaTypeVideo string values.
ALTER TABLE calls MODIFY COLUMN media_type ENUM('audio','video') NOT NULL DEFAULT 'audio';

-- 4. Add columns Go entity/repos expect but were missing in the original schema.
ALTER TABLE calls
  ADD COLUMN channel_uuid        VARCHAR(64)         NULL AFTER tenant_id,
  ADD COLUMN direction           ENUM('inbound','outbound') NOT NULL DEFAULT 'inbound' AFTER call_type,
  ADD COLUMN status              ENUM('ivr','queue','ringing','active','held','consulting','conference','completed','abandoned','failed') NOT NULL DEFAULT 'ivr' AFTER direction,
  ADD COLUMN phone_number_id     BIGINT UNSIGNED     NULL AFTER ivr_flow_id,
  ADD COLUMN carrier_id          BIGINT UNSIGNED     NULL AFTER phone_number_id,
  ADD COLUMN sip_trunk_id        BIGINT UNSIGNED     NULL AFTER carrier_id,
  ADD COLUMN hold_count          INT UNSIGNED        NOT NULL DEFAULT 0,
  ADD COLUMN transfer_count      INT UNSIGNED        NOT NULL DEFAULT 0,
  ADD COLUMN satisfaction_rating TINYINT UNSIGNED    NULL,
  ADD COLUMN duration_sec        INT UNSIGNED        NOT NULL DEFAULT 0,
  ADD COLUMN talk_duration_sec   INT UNSIGNED        NOT NULL DEFAULT 0,
  ADD COLUMN hold_duration_sec   INT UNSIGNED        NOT NULL DEFAULT 0,
  ADD COLUMN acw_duration_sec    INT UNSIGNED        NOT NULL DEFAULT 0,
  ADD COLUMN custom_data         JSON                NULL;

-- 5. Recreate indexes whose column names changed.
ALTER TABLE calls DROP INDEX idx_tenant_time;
ALTER TABLE calls DROP INDEX idx_agent_time;
ALTER TABLE calls DROP INDEX idx_skill_time;
ALTER TABLE calls DROP INDEX idx_campaign;

ALTER TABLE calls ADD INDEX idx_tenant_started   (tenant_id, started_at);
ALTER TABLE calls ADD INDEX idx_agent_started    (agent_user_id, started_at);
ALTER TABLE calls ADD INDEX idx_skill_started    (skill_group_id, started_at);
ALTER TABLE calls ADD INDEX idx_campaign_started (campaign_id, started_at);
ALTER TABLE calls ADD INDEX idx_status_started   (status, started_at);
ALTER TABLE calls ADD INDEX idx_channel_uuid     (channel_uuid);
