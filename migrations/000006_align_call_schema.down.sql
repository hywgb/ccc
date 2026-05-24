-- Revert 000006: best-effort rollback.

ALTER TABLE calls DROP INDEX idx_channel_uuid;
ALTER TABLE calls DROP INDEX idx_status_started;
ALTER TABLE calls DROP INDEX idx_campaign_started;
ALTER TABLE calls DROP INDEX idx_skill_started;
ALTER TABLE calls DROP INDEX idx_agent_started;
ALTER TABLE calls DROP INDEX idx_tenant_started;

ALTER TABLE calls ADD INDEX idx_campaign     (campaign_id, started_at);
ALTER TABLE calls ADD INDEX idx_skill_time   (skill_group_id, started_at);
ALTER TABLE calls ADD INDEX idx_agent_time   (agent_user_id, started_at);
ALTER TABLE calls ADD INDEX idx_tenant_time  (tenant_id, started_at);

ALTER TABLE calls
  DROP COLUMN custom_data,
  DROP COLUMN acw_duration_sec,
  DROP COLUMN hold_duration_sec,
  DROP COLUMN talk_duration_sec,
  DROP COLUMN duration_sec,
  DROP COLUMN satisfaction_rating,
  DROP COLUMN transfer_count,
  DROP COLUMN hold_count,
  DROP COLUMN sip_trunk_id,
  DROP COLUMN carrier_id,
  DROP COLUMN phone_number_id,
  DROP COLUMN status,
  DROP COLUMN direction,
  DROP COLUMN channel_uuid;

ALTER TABLE calls MODIFY COLUMN media_type ENUM('AUDIO','VIDEO') NOT NULL DEFAULT 'AUDIO';

ALTER TABLE calls ADD COLUMN duration_sec INT UNSIGNED GENERATED ALWAYS AS (TIMESTAMPDIFF(SECOND, answered_at, ended_at)) VIRTUAL;

ALTER TABLE calls CHANGE COLUMN hangup_reason hangup_cause VARCHAR(64);
ALTER TABLE calls CHANGE COLUMN caller        cli          VARCHAR(32);
ALTER TABLE calls CHANGE COLUMN ended_at      end_at       TIMESTAMP(3) NULL;
ALTER TABLE calls CHANGE COLUMN answered_at   answer_at    TIMESTAMP(3) NULL;
ALTER TABLE calls CHANGE COLUMN started_at    start_at     TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

-- Re-index report queries that reference renamed columns may need to be reverted manually.
