WITH RECURSIVE doc_key_and_value_recursive(key, value) AS (
  SELECT
    t.key,
    t.value
  FROM store, jsonb_each(store.data) AS t

  UNION ALL

  SELECT
    CONCAT(doc_key_and_value_recursive.key, '.', t.key),
    t.value
  FROM doc_key_and_value_recursive,
    jsonb_each(
      CASE 
        WHEN jsonb_typeof(doc_key_and_value_recursive.value) <> 'object' THEN '{}' :: JSONB
        ELSE doc_key_and_value_recursive.value
      END
      ) AS t
)
SELECT DISTINCT key, jsonb_typeof(value)
FROM doc_key_and_value_recursive
WHERE jsonb_typeof(doc_key_and_value_recursive.value) NOT IN ( 'object')   --'array',
ORDER BY key



--CASE
--WHEN (json_typeof(t.value)='array')
--THEN json_array_elements(t.value)
--ELSE t.value
--END

--BTW : as said in the post, change ligne 12 to the following code in order to avoid table issues