CREATE OR REPLACE FUNCTION operation_signed_sum(
    p_category_id BIGINT,
    p_sum         NUMERIC
) RETURNS NUMERIC AS $$
DECLARE
    v_type_title VARCHAR(50);
BEGIN
    SELECT ot.title
      INTO STRICT v_type_title
      FROM operation_categories oc
      JOIN operation_types ot ON ot.id = oc.operation_type_id
     WHERE oc.id = p_category_id;

    CASE v_type_title
        WHEN 'income'  THEN RETURN  p_sum;
        WHEN 'expense' THEN RETURN -p_sum;
        ELSE
            RAISE EXCEPTION 'Unknown operation type title: %', v_type_title;
    END CASE;
END;
$$ LANGUAGE plpgsql STABLE;

CREATE OR REPLACE FUNCTION calculate_balance()
RETURNS TRIGGER AS $$
DECLARE
    v_account_id BIGINT;
    v_delta      NUMERIC(20,2) := 0;
BEGIN
    IF TG_OP = 'INSERT' THEN
        v_account_id := NEW.account_id;
        v_delta := operation_signed_sum(NEW.category_id, NEW.sum);

    ELSIF TG_OP = 'DELETE' THEN
        v_account_id := OLD.account_id;
        v_delta := -operation_signed_sum(OLD.category_id, OLD.sum);

    ELSIF TG_OP = 'UPDATE' THEN
        IF NEW.account_id <> OLD.account_id THEN
            RAISE EXCEPTION 'account_id cannot be changed (operation id=%)', OLD.id;
        END IF;

        v_account_id := NEW.account_id;
        v_delta := operation_signed_sum(NEW.category_id, NEW.sum)
                 - operation_signed_sum(OLD.category_id, OLD.sum);
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