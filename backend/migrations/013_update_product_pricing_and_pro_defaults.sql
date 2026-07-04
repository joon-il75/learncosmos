UPDATE point_products
SET price_krw = 11900
WHERE name = '프리미엄 팩'
  AND COALESCE(product_type, 'one_time') = 'one_time';

UPDATE point_products
SET monthly_points = 300,
    description = '월 300pt · 커리큘럼 5pt · 튜터 1pt · 비전 코칭 5pt(긴 변 1280px 제한) · 퀴즈 2pt · 커리큘럼: Claude 3.5 Sonnet · 튜터: Grok · 비전 코칭: GPT-4o · 퀴즈: GPT-4o-mini · 광고 제거'
WHERE name = 'LearnWeaver Pro'
  AND COALESCE(product_type, 'subscription') = 'subscription'
  AND billing_period IN ('monthly', 'yearly');
