DROP TRIGGER IF EXISTS trg_update_balance ON operations;

DROP FUNCTION IF EXISTS calculate_balance();

DROP FUNCTION IF EXISTS operation_signed_sum(BIGINT, NUMERIC);