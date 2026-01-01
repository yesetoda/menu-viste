-- Seed Subscription Plans
INSERT INTO subscription_plans (name, slug, description, price_monthly, price_annual, currency, features, display_order)
VALUES 
(
    'Free Trial', 
    'free-trial', 
    '7-day full access trial to explore all features.', 
    0, 
    0, 
    'ETB', 
    '{
        "max_restaurants": 1,
        "max_categories": 5,
        "max_items": 20,
        "analytics": "basic",
        "staff_management": false,
        "custom_branding": false,
        "trial_days": 7
    }'::jsonb, 
    0
),
(
    'Bronze', 
    'bronze-monthly', 
    'Perfect for small cafes and restaurants.', 
    500, 
    5000, 
    'ETB', 
    '{
        "max_restaurants": 1,
        "max_categories": 10,
        "max_items": 50,
        "analytics": "basic",
        "staff_management": false,
        "custom_branding": false
    }'::jsonb, 
    1
),
(
    'Silver', 
    'silver-monthly', 
    'Ideal for growing businesses with multiple locations.', 
    1500, 
    15000, 
    'ETB', 
    '{
        "max_restaurants": 3,
        "max_categories": -1,
        "max_items": 200,
        "analytics": "advanced",
        "staff_management": true,
        "custom_branding": false
    }'::jsonb, 
    2
),
(
    'Gold', 
    'gold-monthly', 
    'Unlimited power for large enterprises and chains.', 
    3000, 
    30000, 
    'ETB', 
    '{
        "max_restaurants": -1,
        "max_categories": -1,
        "max_items": -1,
        "analytics": "full",
        "staff_management": true,
        "custom_branding": true
    }'::jsonb, 
    3
)
ON CONFLICT (slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    price_monthly = EXCLUDED.price_monthly,
    price_annual = EXCLUDED.price_annual,
    currency = EXCLUDED.currency,
    features = EXCLUDED.features,
    display_order = EXCLUDED.display_order,
    updated_at = NOW();
