ALTER TABLE point_products
    ADD COLUMN IF NOT EXISTS product_type TEXT NOT NULL DEFAULT 'one_time',
    ADD COLUMN IF NOT EXISTS billing_period TEXT,
    ADD COLUMN IF NOT EXISTS monthly_points INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ad_free BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS premium_access BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS description TEXT;

UPDATE point_products
SET product_type = 'one_time'
WHERE product_type IS NULL OR product_type = '';

INSERT INTO system_llm_settings (feature, provider, model)
VALUES ('vision', 'openai', 'gpt-4o')
ON CONFLICT (feature) DO NOTHING;

UPDATE system_llm_settings
SET provider = 'anthropic', model = 'claude-3-5-sonnet', updated_at = NOW()
WHERE feature = 'curriculum';

UPDATE system_llm_settings
SET provider = 'grok', model = 'grok-2', updated_at = NOW()
WHERE feature = 'tutor';

UPDATE system_llm_settings
SET provider = 'openai', model = 'gpt-4o', updated_at = NOW()
WHERE feature = 'vision';

INSERT INTO point_products (
    name, points, price_krw, is_active, sort_order,
    product_type, billing_period, monthly_points, ad_free, premium_access, description
)
SELECT
    'LearnWeaver Pro',
    0,
    0,
    false,
    100,
    'subscription',
    'monthly',
    0,
    true,
    true,
    '커리큘럼: Claude 3.5 Sonnet · 튜터: Grok · 비전 코칭: GPT-4o · 임베딩: text-embedding-3-small · 광고 제거'
WHERE NOT EXISTS (
    SELECT 1 FROM point_products WHERE name = 'LearnWeaver Pro' AND billing_period = 'monthly'
);

INSERT INTO point_products (
    name, points, price_krw, is_active, sort_order,
    product_type, billing_period, monthly_points, ad_free, premium_access, description
)
SELECT
    'LearnWeaver Pro',
    0,
    0,
    false,
    101,
    'subscription',
    'yearly',
    0,
    true,
    true,
    '커리큘럼: Claude 3.5 Sonnet · 튜터: Grok · 비전 코칭: GPT-4o · 임베딩: text-embedding-3-small · 광고 제거'
WHERE NOT EXISTS (
    SELECT 1 FROM point_products WHERE name = 'LearnWeaver Pro' AND billing_period = 'yearly'
);
