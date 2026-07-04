-- 002_admin_totp.sql
-- admin TOTP 설정을 위한 users 테이블 컬럼 추가

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS totp_secret TEXT,
    ADD COLUMN IF NOT EXISTS totp_enabled BOOLEAN NOT NULL DEFAULT false;

-- 확인용 쿼리:
-- SELECT id, email, role, totp_secret, totp_enabled FROM users;
