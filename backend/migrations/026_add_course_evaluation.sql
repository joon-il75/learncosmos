-- 026: 코스 종합 평가 AI 기능
-- point_settings에 course_eval_cost 키 추가
INSERT INTO point_settings (key, value, updated_at)
VALUES ('course_eval_cost', 3, NOW())
ON CONFLICT (key) DO NOTHING;

-- course_drafts에 평가 결과 컬럼 추가
ALTER TABLE course_drafts
    ADD COLUMN IF NOT EXISTS evaluation_result JSONB,
    ADD COLUMN IF NOT EXISTS evaluated_at TIMESTAMPTZ;
