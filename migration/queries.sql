-- name: CreateUser :one
INSERT INTO users (
    email, password_hash, full_name, role, owner_id, restaurant_id, phone, avatar_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET 
    full_name = COALESCE(sqlc.narg('full_name'), full_name),
    phone = COALESCE(sqlc.narg('phone'), phone),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    email_verified = COALESCE(sqlc.narg('email_verified'), email_verified),
    last_login_at = COALESCE(sqlc.narg('last_login_at'), last_login_at),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteUser :exec
UPDATE users SET deleted_at = NOW() WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListUsersWithFilters :many
SELECT * FROM users
WHERE 
    (sqlc.narg('email')::text IS NULL OR email = sqlc.narg('email')) AND
    (sqlc.narg('role')::user_role IS NULL OR role = sqlc.narg('role')) AND
    (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active')) AND
    (sqlc.narg('search')::text IS NULL OR full_name ILIKE '%' || sqlc.narg('search') || '%' OR email ILIKE '%' || sqlc.narg('search') || '%') AND
    deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListStaffByOwner :many
SELECT * FROM users
WHERE owner_id = $1 AND role = 'staff' AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListStaffByRestaurant :many
SELECT * FROM users
WHERE restaurant_id = $1 AND role = 'staff' AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateStaffStatus :exec
UPDATE users
SET is_active = $3, updated_at = NOW()
WHERE id = $1 AND restaurant_id = $2 AND role = 'staff';

-- name: DeleteStaff :exec
UPDATE users
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND restaurant_id = $2 AND role = 'staff';

-- name: CreateRestaurant :one
INSERT INTO restaurants (
    owner_id, name, slug, description, cuisine_type, phone, email, website, address, city, country, logo_url, cover_image_url, theme_settings
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: GetRestaurantBySlug :one
SELECT * FROM restaurants
WHERE slug = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetRestaurantByID :one
SELECT * FROM restaurants
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListRestaurantsByOwner :many
SELECT * FROM restaurants
WHERE owner_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateRestaurant :one
UPDATE restaurants
SET 
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    cuisine_type = COALESCE(sqlc.narg('cuisine_type'), cuisine_type),
    phone = COALESCE(sqlc.narg('phone'), phone),
    email = COALESCE(sqlc.narg('email'), email),
    website = COALESCE(sqlc.narg('website'), website),
    address = COALESCE(sqlc.narg('address'), address),
    city = COALESCE(sqlc.narg('city'), city),
    country = COALESCE(sqlc.narg('country'), country),
    logo_url = COALESCE(sqlc.narg('logo_url'), logo_url),
    cover_image_url = COALESCE(sqlc.narg('cover_image_url'), cover_image_url),
    theme_settings = COALESCE(sqlc.narg('theme_settings'), theme_settings),
    is_published = COALESCE(sqlc.narg('is_published'), is_published),
    updated_at = NOW()
WHERE id = sqlc.arg('id') AND (owner_id = sqlc.arg('owner_id') OR sqlc.arg('is_admin')::boolean)
RETURNING *;



-- name: ListRestaurantsWithFilters :many
SELECT * FROM restaurants
WHERE 
    (sqlc.narg('owner_id')::uuid IS NULL OR owner_id = sqlc.narg('owner_id')) AND
    (sqlc.narg('cuisine_type')::text IS NULL OR cuisine_type = sqlc.narg('cuisine_type')) AND
    (sqlc.narg('city')::text IS NULL OR city = sqlc.narg('city')) AND
    (sqlc.narg('country')::text IS NULL OR country = sqlc.narg('country')) AND
    (sqlc.narg('is_published')::boolean IS NULL OR is_published = sqlc.narg('is_published')) AND
    (sqlc.narg('search')::text IS NULL OR name ILIKE '%' || sqlc.narg('search') || '%') AND
    deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeleteRestaurant :exec
UPDATE restaurants SET deleted_at = NOW() WHERE id = $1 AND owner_id = $2;

-- name: CreateCategory :one
INSERT INTO categories (
    restaurant_id, name, description, icon, display_order, is_active, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: ListCategoriesByRestaurant :many
SELECT * FROM categories
WHERE restaurant_id = $1
ORDER BY display_order ASC;

-- name: UpdateCategory :one
UPDATE categories
SET 
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    icon = COALESCE(sqlc.narg('icon'), icon),
    display_order = COALESCE(sqlc.narg('display_order'), display_order),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: CreateMenuItem :one
INSERT INTO menu_items (
    restaurant_id, category_id, name, description, price, currency, images, allergens, dietary_tags, spice_level, calories, is_available, display_order, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: GetMenuItemByID :one
SELECT * FROM menu_items
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListMenuItemsByCategory :many
SELECT * FROM menu_items
WHERE category_id = $1 AND deleted_at IS NULL
ORDER BY display_order ASC;

-- name: ListMenuItemsByRestaurant :many
SELECT * FROM menu_items
WHERE restaurant_id = $1 AND deleted_at IS NULL
ORDER BY category_id, display_order ASC;

-- name: UpdateMenuItem :one
UPDATE menu_items
SET 
    category_id = COALESCE(sqlc.narg('category_id'), category_id),
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    price = COALESCE(sqlc.narg('price'), price),
    images = COALESCE(sqlc.narg('images'), images),
    allergens = COALESCE(sqlc.narg('allergens'), allergens),
    dietary_tags = COALESCE(sqlc.narg('dietary_tags'), dietary_tags),
    spice_level = COALESCE(sqlc.narg('spice_level'), spice_level),
    calories = COALESCE(sqlc.narg('calories'), calories),
    is_available = COALESCE(sqlc.narg('is_available'), is_available),
    display_order = COALESCE(sqlc.narg('display_order'), display_order),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteMenuItem :exec
UPDATE menu_items SET deleted_at = NOW() WHERE id = $1;

-- name: CreateActivityLog :one
INSERT INTO activity_logs (
    restaurant_id, user_id, action_type, action_category, description, target_type, target_id, target_name, before_value, after_value, ip_address, user_agent, device_type, browser, os, success
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
) RETURNING *;

-- name: ListActivityLogsByRestaurant :many
SELECT al.*, u.full_name as user_name, u.email as user_email
FROM activity_logs al
JOIN users u ON al.user_id = u.id
WHERE al.restaurant_id = $1
ORDER BY al.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListActivityLogsWithFilters :many
SELECT al.*, u.full_name as user_name, u.email as user_email
FROM activity_logs al
JOIN users u ON al.user_id = u.id
WHERE 
    (sqlc.narg('restaurant_id')::uuid IS NULL OR al.restaurant_id = sqlc.narg('restaurant_id')) AND
    (sqlc.narg('user_id')::uuid IS NULL OR al.user_id = sqlc.narg('user_id')) AND
    (sqlc.narg('action_type')::text IS NULL OR al.action_type = sqlc.narg('action_type')) AND
    (sqlc.narg('action_category')::text IS NULL OR al.action_category = sqlc.narg('action_category')) AND
    (sqlc.narg('target_type')::text IS NULL OR al.target_type = sqlc.narg('target_type')) AND
    (sqlc.narg('target_id')::uuid IS NULL OR al.target_id = sqlc.narg('target_id')) AND
    (sqlc.narg('success')::boolean IS NULL OR al.success = sqlc.narg('success')) AND
    (sqlc.narg('search')::text IS NULL OR al.description ILIKE '%' || sqlc.narg('search') || '%' OR al.target_name ILIKE '%' || sqlc.narg('search') || '%')
ORDER BY al.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateAnalyticsEvent :one
INSERT INTO analytics_events (
    restaurant_id, event_type, visitor_id, session_id, target_id, ip_address, device_type, browser, os, country, city
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: ListAnalyticsEventsWithFilters :many
SELECT * FROM analytics_events
WHERE 
    (sqlc.narg('restaurant_id')::uuid IS NULL OR restaurant_id = sqlc.narg('restaurant_id')) AND
    (sqlc.narg('event_type')::text IS NULL OR event_type = sqlc.narg('event_type')) AND
    (sqlc.narg('visitor_id')::text IS NULL OR visitor_id = sqlc.narg('visitor_id')) AND
    (sqlc.narg('session_id')::uuid IS NULL OR session_id = sqlc.narg('session_id')) AND
    (sqlc.narg('target_id')::uuid IS NULL OR target_id = sqlc.narg('target_id')) AND
    (sqlc.narg('country')::text IS NULL OR country = sqlc.narg('country')) AND
    (sqlc.narg('city')::text IS NULL OR city = sqlc.narg('city'))
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetAnalyticsAggregates :many
SELECT * FROM analytics_aggregates
WHERE restaurant_id = $1 AND date >= $2 AND date <= $3
ORDER BY date ASC, hour ASC;

-- name: UpsertAnalyticsAggregate :one
INSERT INTO analytics_aggregates (
    restaurant_id, date, hour, metric_type, target_id, value
) VALUES (
    $1, $2, $3, $4, $5, $6
) ON CONFLICT (restaurant_id, date, hour, metric_type, target_id)
DO UPDATE SET 
    value = analytics_aggregates.value + EXCLUDED.value,
    updated_at = NOW()
RETURNING *;

-- name: CreateSubscriptionPlan :one
INSERT INTO subscription_plans (
    name, slug, description, price_monthly, price_annual, currency, features, display_order
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: ListSubscriptionPlans :many
SELECT * FROM subscription_plans
WHERE is_active = TRUE
ORDER BY display_order ASC;

-- name: GetSubscriptionPlanBySlug :one
SELECT * FROM subscription_plans
WHERE slug = $1 LIMIT 1;

-- name: CreateSubscription :one
INSERT INTO subscriptions (
    owner_id, plan_id, status, current_period_start, current_period_end, trial_end
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetSubscriptionByOwner :one
SELECT s.*, sp.name as plan_name, sp.slug as plan_slug, sp.features
FROM subscriptions s
JOIN subscription_plans sp ON s.plan_id = sp.id
WHERE s.owner_id = $1 LIMIT 1;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET 
    plan_id = COALESCE(sqlc.narg('plan_id'), plan_id),
    status = COALESCE(sqlc.narg('status'), status),
    current_period_start = COALESCE(sqlc.narg('current_period_start'), current_period_start),
    current_period_end = COALESCE(sqlc.narg('current_period_end'), current_period_end),
    trial_end = COALESCE(sqlc.narg('trial_end'), trial_end),
    cancelled_at = COALESCE(sqlc.narg('cancelled_at'), cancelled_at),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: CreateInvoice :one
INSERT INTO invoices (
    subscription_id, owner_id, invoice_number, amount, currency, status, billing_period_start, billing_period_end
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;
-- name: CreatePaymentTransaction :one
INSERT INTO payment_transactions (
    owner_id, amount, currency, status, tx_ref, reference
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetPaymentTransactionByTxRef :one
SELECT * FROM payment_transactions
WHERE tx_ref = $1 LIMIT 1;

-- name: UpdatePaymentTransactionStatus :one
UPDATE payment_transactions
SET 
    status = $2,
    provider_transaction_ref = COALESCE(sqlc.narg('provider_transaction_ref'), provider_transaction_ref),
    updated_at = NOW()
WHERE tx_ref = $1
RETURNING *;

-- name: CreatePaymentWebhook :one
INSERT INTO payment_webhooks (
    provider_event_id, event_type, payload
) VALUES (
    $1, $2, $3
) ON CONFLICT (provider_event_id) DO NOTHING
RETURNING *;

-- name: MarkWebhookAsProcessed :exec
UPDATE payment_webhooks SET processed = TRUE WHERE provider_event_id = $1;

-- name: CreatePaymentRetryJob :one
INSERT INTO payment_retry_jobs (
    subscription_id, scheduled_for
) VALUES (
    $1, $2
) RETURNING *;

-- name: UpdatePaymentRetryJob :one
UPDATE payment_retry_jobs
SET 
    status = $2,
    retry_count = retry_count + 1,
    scheduled_for = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListInvoicesByOwner :many
SELECT * FROM invoices
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: ListInvoicesWithFilters :many
SELECT * FROM invoices
WHERE 
    (sqlc.narg('owner_id')::uuid IS NULL OR owner_id = sqlc.narg('owner_id')) AND
    (sqlc.narg('subscription_id')::uuid IS NULL OR subscription_id = sqlc.narg('subscription_id')) AND
    (sqlc.narg('status')::invoice_status IS NULL OR status = sqlc.narg('status'))
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateInvoiceStatus :one
UPDATE invoices
SET 
    status = $2,
    paid_at = CASE WHEN $2 = 'paid'::invoice_status THEN NOW() ELSE paid_at END,
    updated_at = NOW()
WHERE invoice_number = $1
RETURNING *;

-- name: IncrementRestaurantViewCount :exec
UPDATE restaurants SET view_count = view_count + 1 WHERE id = $1;

-- name: IncrementMenuItemViewCount :exec
UPDATE menu_items SET view_count = view_count + 1 WHERE id = $1;

-- name: GetAdminDashboardStats :one
SELECT 
    (SELECT COUNT(*) FROM users WHERE deleted_at IS NULL) as total_users,
    (SELECT COUNT(*) FROM restaurants WHERE deleted_at IS NULL) as total_restaurants,
    (SELECT COUNT(*) FROM subscriptions WHERE status = 'active') as active_subscriptions,
    (SELECT SUM(amount) FROM invoices WHERE status = 'paid') as total_revenue;

-- name: GetRecentAdminLogs :many
SELECT al.*, u.full_name as user_name, u.email as user_email
FROM activity_logs al
JOIN users u ON al.user_id = u.id
ORDER BY al.created_at DESC
LIMIT $1;

-- name: GetRestaurantDetailsForAdmin :one
SELECT r.*, u.full_name as owner_name, u.email as owner_email
FROM restaurants r
JOIN users u ON r.owner_id = u.id
WHERE r.id = $1;

-- name: GetAllAdminEmails :many
SELECT email FROM users
WHERE role = 'admin' AND deleted_at IS NULL;
