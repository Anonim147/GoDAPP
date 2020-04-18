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

