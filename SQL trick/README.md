# 数据库判断

1. `len`和`length`

`select * from xxx where id=1 and length('a')=1`

正常则mysql或mssql,否则为oracle

2. `version()`和`@@vserion`

mysql 都可以

mssql `@@version`

3. `substring`和`substr`

mysql 都可以

mssql `substring`

oracle `substr`

