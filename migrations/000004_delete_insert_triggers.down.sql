DROP TRIGGER IF EXISTS trg_operations_signed_sum ON operations;

DROP FUNCTION IF EXISTS operations_set_signed_sum();

DROP TRIGGER IF EXISTS trg_reassign_operations_before_category_delete ON operation_categories;

DROP FUNCTION IF EXISTS reassign_operations_before_category_delete();