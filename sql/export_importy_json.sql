\copy (select row_to_json(d) 
    from 
    (SELECT data #> '{ ingredients }' as "ingredients", 
            data #> '{ dimensions}' as "dimensions", 
            data #> '{ name}' as "name" 
                FROM store 
                WHERE  data #>> '{ dimensions,weight }' like '%50%'  
                    OR  data #>> '{ name }' = 'Pizza') d) 
    to 'D:\temp\file.json';


insert into store(data) 
    select row_to_json(d) from 
     (SELECT data #> '{ ingredients }' as "ingredients", 
            data #> '{ dimensions}' as "dimensions", 
            data #> '{ name}' as "name" 
                FROM store 
                WHERE  data #>> '{ dimensions,weight }' like '%50%'  
                    OR  data #>> '{ name }' = 'Pizza') d 



\copy temp from 'D:\temp\output.json' csv quote e'\x01' delimiter e'\x02';

select value from temp, json_array_elements(temp.data::json) as elem;

insert into store(data) 

    select values from (select jsonb_array_elements(temp.data::jsonb) as values from temp) temp;