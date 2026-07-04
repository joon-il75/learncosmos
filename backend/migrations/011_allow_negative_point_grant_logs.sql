ALTER TABLE point_grant_logs
    DROP CONSTRAINT IF EXISTS point_grant_logs_amount_check;

ALTER TABLE point_grant_logs
    ADD CONSTRAINT point_grant_logs_amount_check
    CHECK (amount <> 0 AND amount >= -50 AND amount <= 50);
