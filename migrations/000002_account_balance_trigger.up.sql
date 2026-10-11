CREATE OR REPLACE FUNCTION calculate_balance()
RETURNS TRIGGER AS $$
DECLARE
    v_account_id BIGINT;
    v_delta      NUMERIC(20,2) := 0;
BEGIN
    IF TG_OP = 'INSERT' THEN
        v_account_id := NEW.account_id;
        v_delta := NEW.sum;

    ELSIF TG_OP = 'DELETE' THEN
        v_account_id := OLD.account_id;
        v_delta := -OLD.sum;

    ELSIF TG_OP = 'UPDATE' THEN
        IF NEW.account_id <> OLD.account_id THEN
            RAISE EXCEPTION 'account_id cannot be changed (operation id=%)', OLD.id;
        END IF;

        v_account_id := NEW.account_id;
        v_delta := NEW.sum
                 - OLD.sum;
    END IF;

    IF v_delta <> 0 THEN
        UPDATE accounts
           SET balance = balance + v_delta,
               version = version + 1
         WHERE id = v_account_id;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_balance
AFTER INSERT OR UPDATE OR DELETE ON operations
FOR EACH ROW EXECUTE FUNCTION calculate_balance();