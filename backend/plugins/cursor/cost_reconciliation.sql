-- Cursor cost reconciliation queries
-- Run against the lake database after a full billing-cycle event backfill.
--
-- Note: SUM(charged_cents) from events only reconciles with /teams/spend when
-- ALL usage events for the billing cycle are collected (paginate filtered-usage-events).

-- Team-level rollup comparison
SELECT
  ROUND(SUM(spend_cents) / 100, 2) AS spend_usd,
  ROUND(SUM(included_spend_cents) / 100, 2) AS included_usd,
  ROUND(SUM(spend_cents + included_spend_cents) / 100, 2) AS total_cycle_usd
FROM _tool_cursor_user_spend;

SELECT
  ROUND(SUM(charged_cents) / 100, 2) AS event_charged_usd,
  ROUND(SUM(total_cents) / 100, 2) AS event_model_usd,
  COUNT(*) AS event_count,
  MIN(event_time) AS earliest_event,
  MAX(event_time) AS latest_event
FROM _tool_cursor_usage_events;

-- Per-user: billing cycle spend vs event charged (since billing cycle start)
SELECT
  s.email,
  ROUND(s.spend_cents / 100, 2) AS spend_usd,
  ROUND(s.included_spend_cents / 100, 2) AS included_usd,
  ROUND((s.spend_cents + s.included_spend_cents) / 100, 2) AS total_cycle_usd,
  ROUND(COALESCE(SUM(e.charged_cents), 0) / 100, 2) AS event_charged_usd,
  COUNT(e.event_id) AS event_count
FROM _tool_cursor_user_spend s
LEFT JOIN _tool_cursor_usage_events e
  ON s.email = e.user_email
 AND e.event_time >= s.billing_cycle_start
GROUP BY s.email, s.spend_cents, s.included_spend_cents
ORDER BY total_cycle_usd DESC;

-- Users with events but no spend row (should be empty)
SELECT DISTINCT e.user_email
FROM _tool_cursor_usage_events e
LEFT JOIN _tool_cursor_user_spend s ON e.user_email = s.email
WHERE s.email IS NULL;
