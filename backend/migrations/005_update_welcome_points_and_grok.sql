-- 004_update_welcome_points_and_grok.sql
-- v5: welcome_points 30pt 반영 + Grok provider 추가

BEGIN;

-- ① point_settings: welcome_points 값 30으로 업데이트
INSERT INTO point_settings (key, value)
VALUES ('welcome_points', 30)
ON CONFLICT (key) DO UPDATE
  SET value      = 30,
      updated_at = NOW();

-- ② system_api_keys: grok 키 슬롯 추가 (없을 때만)
INSERT INTO system_api_keys (provider)
VALUES ('grok')
ON CONFLICT (provider) DO NOTHING;

COMMIT;
