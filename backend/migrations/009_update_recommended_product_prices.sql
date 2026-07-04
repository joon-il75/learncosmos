UPDATE point_products
SET price_krw = 1900
WHERE name = '스타터 팩' AND COALESCE(product_type, 'one_time') = 'one_time';

UPDATE point_products
SET price_krw = 4900
WHERE name = '베이직 팩' AND COALESCE(product_type, 'one_time') = 'one_time';

UPDATE point_products
SET price_krw = 9900
WHERE name = '프리미엄 팩' AND COALESCE(product_type, 'one_time') = 'one_time';

UPDATE point_products
SET price_krw = 19900
WHERE name = 'LearnWeaver Pro' AND product_type = 'subscription' AND billing_period = 'monthly';

UPDATE point_products
SET price_krw = 199000
WHERE name = 'LearnWeaver Pro' AND product_type = 'subscription' AND billing_period = 'yearly';
