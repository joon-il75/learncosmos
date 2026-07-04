-- 005_point_grant_audit.sql
-- 관리자 포인트 수동 지급 감사 로그 테이블

BEGIN;

CREATE TABLE IF NOT EXISTS point_grant_logs (
  id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  admin_id          TEXT        NOT NULL,  -- 지급한 관리자 ID (super_admin은 "super_admin" 문자열)
  admin_email       TEXT,                  -- 관리자 이메일 (super_admin은 환경변수 ID)
  target_user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  target_identifier TEXT        NOT NULL,  -- display_id 또는 email (조회 편의용)
  amount            INT         NOT NULL CHECK (amount BETWEEN 1 AND 50),
  memo              TEXT,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_point_grant_logs_created_at
  ON point_grant_logs (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_point_grant_logs_admin_id
  ON point_grant_logs (admin_id);

CREATE INDEX IF NOT EXISTS idx_point_grant_logs_target_user
  ON point_grant_logs (target_user_id);

COMMIT;
