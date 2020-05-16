WITH RECURSIVE doc_key_and_value_recursive(key, value) 
AS 
( 
    SELECT 
    CASE WHEN JSON_TYPEOF(t.value) = 'array' 
        THEN CONCAT(t.key,'[]') 
    ELSE t.key END AS key, 
    CASE WHEN JSON_TYPEOF(t.value) = 'array' 
        THEN JSON_ARRAY_ELEMENTS(t.value) 
    ELSE JSON_BUILD_ARRAY(t.value) 
    END AS value, 
    CASE WHEN JSON_TYPEOF(t.value) = 'array' 
        THEN JSON_ARRAY_LENGTH(t.value) 
    ELSE NULL 
    END AS i 
    FROM store1, JSON_EACH(store1.data) AS t 
    -- WHERE something = 123 # apply filtering here 
    UNION ALL 
    SELECT 
        CASE WHEN JSON_TYPEOF(t.value) = 'array' 
            THEN CONCAT(doc_key_and_value_recursive.key, '.', t.key, '[]') 
            ELSE CONCAT(doc_key_and_value_recursive.key, '.', t.key) 
        END AS key,
        CASE WHEN JSON_TYPEOF(t.value) = 'array' 
            THEN JSON_ARRAY_ELEMENTS(t.value) 
            ELSE JSON_BUILD_ARRAY(t.value) 
        END AS value, 
        CASE WHEN JSON_TYPEOF(t.value) = 'array' 
        THEN JSON_ARRAY_LENGTH(t.value) 
        ELSE NULL 
        END AS i 
    FROM doc_key_and_value_recursive, 
        JSON_EACH( 
            CASE WHEN JSON_TYPEOF(doc_key_and_value_recursive.value) <> 'object' THEN '{}' :: JSON 
            ELSE doc_key_and_value_recursive.value 
            END 
        ) AS t
)
SELECT key, MAX(i) FROM doc_key_and_value_recursive GROUP BY key ORDER BY key