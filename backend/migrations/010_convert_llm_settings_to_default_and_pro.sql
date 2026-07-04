INSERT INTO system_llm_settings (feature, provider, model, endpoint_url, updated_at)
SELECT 'default', provider, model, endpoint_url, NOW()
FROM system_llm_settings
WHERE feature = 'default'
ON CONFLICT (feature) DO NOTHING;

INSERT INTO system_llm_settings (feature, provider, model, endpoint_url, updated_at)
SELECT 'pro_curriculum', provider, model, endpoint_url, NOW()
FROM system_llm_settings
WHERE feature = 'curriculum'
ON CONFLICT (feature) DO UPDATE
SET provider = EXCLUDED.provider, model = EXCLUDED.model, endpoint_url = EXCLUDED.endpoint_url, updated_at = NOW();

INSERT INTO system_llm_settings (feature, provider, model, endpoint_url, updated_at)
SELECT 'pro_tutor', provider, model, endpoint_url, NOW()
FROM system_llm_settings
WHERE feature = 'tutor'
ON CONFLICT (feature) DO UPDATE
SET provider = EXCLUDED.provider, model = EXCLUDED.model, endpoint_url = EXCLUDED.endpoint_url, updated_at = NOW();

INSERT INTO system_llm_settings (feature, provider, model, endpoint_url, updated_at)
SELECT 'pro_vision', provider, model, endpoint_url, NOW()
FROM system_llm_settings
WHERE feature = 'vision'
ON CONFLICT (feature) DO UPDATE
SET provider = EXCLUDED.provider, model = EXCLUDED.model, endpoint_url = EXCLUDED.endpoint_url, updated_at = NOW();

INSERT INTO system_llm_settings (feature, provider, model, endpoint_url, updated_at)
SELECT 'pro_quiz', provider, model, endpoint_url, NOW()
FROM system_llm_settings
WHERE feature = 'quiz'
ON CONFLICT (feature) DO UPDATE
SET provider = EXCLUDED.provider, model = EXCLUDED.model, endpoint_url = EXCLUDED.endpoint_url, updated_at = NOW();

INSERT INTO system_llm_settings (feature, provider, model, updated_at)
VALUES
    ('pro_curriculum', 'anthropic', 'claude-3-5-sonnet', NOW()),
    ('pro_tutor', 'grok', 'grok-2', NOW()),
    ('pro_vision', 'openai', 'gpt-4o', NOW()),
    ('pro_quiz', 'openai', 'gpt-4o-mini', NOW())
ON CONFLICT (feature) DO NOTHING;

DELETE FROM system_llm_settings
WHERE feature IN ('curriculum', 'tutor', 'vision', 'quiz');
