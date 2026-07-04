-- LLM 설정 테이블 (기능별 LLM 분리)
CREATE TABLE IF NOT EXISTS system_llm_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feature TEXT NOT NULL UNIQUE, -- 'default' | 'curriculum' | 'tutor' | 'quiz'
    provider TEXT NOT NULL DEFAULT 'openai',
    model TEXT NOT NULL DEFAULT 'gpt-4o-mini',
    api_key_encrypted TEXT,       -- AES-256 암호화
    endpoint_url TEXT,            -- Ollama 자체 호스팅 URL
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 기본값 시드
INSERT INTO system_llm_settings (feature, provider, model)
VALUES
    ('default',    'openai', 'gpt-4o-mini'),
    ('curriculum', 'openai', 'gpt-4o-mini'),
    ('tutor',      'openai', 'gpt-4o-mini'),
    ('quiz',       'openai', 'gpt-4o-mini')
ON CONFLICT (feature) DO NOTHING;

-- 제공자별 API 키 테이블
CREATE TABLE IF NOT EXISTS system_api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider TEXT NOT NULL UNIQUE, -- 'openai'|'anthropic'|'google'|'solar'|'hyperclova'|'llama'|'exaone'
    api_key_encrypted TEXT,
    api_key_secondary_encrypted TEXT, -- HyperCLOVA APIGW Key 등 2번째 키
    endpoint_url TEXT,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- API 키 기본 행 삽입
INSERT INTO system_api_keys (provider) VALUES
    ('openai'), ('anthropic'), ('google'),
    ('solar'), ('hyperclova'), ('llama'), ('exaone')
ON CONFLICT (provider) DO NOTHING;

-- 포인트 정책 테이블
CREATE TABLE IF NOT EXISTS point_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT NOT NULL UNIQUE,
    value INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 포인트 정책 기본값
INSERT INTO point_settings (key, value) VALUES
    ('welcome_points',  30),
    ('course_gen_cost', 5),
    ('lesson_rec_cost', 1),
    ('admin_max_grant', 50)
ON CONFLICT (key) DO NOTHING;

-- 포인트 구매 상품 테이블
CREATE TABLE IF NOT EXISTS point_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    points INT NOT NULL,
    price_krw INT NOT NULL,
    is_active BOOLEAN DEFAULT false,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 기본 상품 시드
INSERT INTO point_products (name, points, price_krw, is_active, sort_order)
VALUES
    ('스타터 팩',   30,  1900, false, 1),
    ('베이직 팩',   80,  3900, false, 2),
    ('프리미엄 팩', 200, 7900, false, 3)
ON CONFLICT DO NOTHING;

-- 광고 슬롯 테이블
CREATE TABLE IF NOT EXISTS ad_slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_type TEXT NOT NULL, -- 'category_sponsor'|'step_complete'|'dashboard_banner'|'creator_sponsor'
    title TEXT NOT NULL,
    link_url TEXT NOT NULL,
    category_slug TEXT,      -- 카테고리 연동
    is_active BOOLEAN DEFAULT false,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 제휴 설정 테이블
CREATE TABLE IF NOT EXISTS affiliate_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider TEXT NOT NULL UNIQUE, -- 'coupang'|'naver'|'class101'|'kyobo'|'adpick'
    tracking_id TEXT,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO affiliate_settings (provider) VALUES
    ('coupang'), ('naver'), ('class101'), ('kyobo'), ('adpick')
ON CONFLICT (provider) DO NOTHING;
