
--!!warning!! delete all spaces first
\copy temp from 'D:\temp\output.json' csv quote e'\x01' delimiter e'\x02';

select value from temp, json_array_elements(temp.data::json) as elem;