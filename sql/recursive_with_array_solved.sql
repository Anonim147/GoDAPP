with recursive extract_all as
(
    select 
        key as path, 
        value
    from test_csv
    cross join lateral jsonb_each(data)
union all
    select
        path || '.' || coalesce(obj_key, '[]'),
        coalesce(obj_value, arr_value)
    from extract_all
    left join lateral 
        jsonb_each(case jsonb_typeof(value) when 'object' then value end) 
        as o(obj_key, obj_value) 
        on jsonb_typeof(value) = 'object'
    left join lateral 
        jsonb_array_elements(case jsonb_typeof(value) when 'array' then value end) 
        with ordinality as a(arr_value, arr_key)
        on jsonb_typeof(value) = 'array'
    where obj_key is not null or arr_key is not null
)
select distinct path, jsonb_typeof(value)
from extract_all where jsonb_typeof(value)<>'null' order by path;