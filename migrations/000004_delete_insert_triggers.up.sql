CREATE OR REPLACE FUNCTION reassign_operations_before_category_delete()
RETURNS TRIGGER AS $$
DECLARE
    v_default_category_id BIGINT;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE id = OLD.user_id) THEN
        DELETE FROM operations WHERE category_id = OLD.id;
        RETURN OLD;
    END IF;

    SELECT id
      INTO v_default_category_id
      FROM operation_categories
     WHERE operation_type_id = OLD.operation_type_id
       AND user_id IS NOT DISTINCT FROM OLD.user_id
       AND id <> OLD.id
       AND title = 'no_category'
     LIMIT 1;

    IF v_default_category_id IS NULL THEN
        RAISE EXCEPTION 'No default category for operation type %', OLD.operation_type_id;
    END IF;

    UPDATE operations
       SET category_id = v_default_category_id,
           version = version + 1
     WHERE category_id = OLD.id;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_reassign_operations_before_category_delete
BEFORE DELETE ON operation_categories
FOR EACH ROW
EXECUTE FUNCTION reassign_operations_before_category_delete();



CREATE OR REPLACE FUNCTION operations_set_signed_sum()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
  v_type TEXT;
BEGIN
  SELECT ot.title INTO v_type
  FROM operation_categories oc
  JOIN operation_types ot ON ot.id = oc.operation_type_id
  WHERE oc.id = NEW.category_id;

  NEW.sum := CASE WHEN v_type = 'expense' THEN -ABS(NEW.sum) ELSE ABS(NEW.sum) END;
  RETURN NEW;
END $$;

CREATE TRIGGER trg_operations_signed_sum
BEFORE INSERT OR UPDATE OF sum, category_id ON operations
FOR EACH ROW EXECUTE FUNCTION operations_set_signed_sum();