>**SQL对大小写不敏感!!!**

SQL语句可以分为以下三类

1. DDL(Data Definition Language,数据定义语言)

用来创建或者删除存储数据用的数据库以及数据库中的表等对象,DDL包含以下几种指令

`CREATE`创建数据库和表等对象

`DROP` 删除数据库和表等对象

`ALTER` 修改数据库和表等对象的结构

2. DML(Data Manipulation Language,数据操纵语言)

用来查询或者变更表中的记录,DML包含以下几种指令

`SELECT`查询表中的数据

`INSERT`向表中插入新数据

`UPDATE`更新表中的数据

`DELETE`删除表中的数据

3. DCL(Data Control Language,数据控制语言)

用来确认或者取消对数据库中的数据进行的变更,对用户权限进行设定,DCL 包含以下几种指令

`COMMIT`确认对数据库中的数据进行的变更

`ROLLBACK`取消对数据库中的数据进行的变更

`GRANT`赋予用户操作权限

`REVOKE`取消用户的操作权限

## DDL

1. 创建数据库

`create database person;`创建一个名为`person`的数据库

2. 创建表

`create table students(ID INTEGER NOT NULL,Name VARCHAR(100) NOT NULL,City VARCHAR(100),PRIMARY KEY (ID));`创建一个名为`students`的表,一共有3列

第一列为`ID`,数据类型为整数型,非空

第二列为`Name`,数据类型为可变长字符串,非空

第三列为`City`,数据类型为可变长字符串

数据库表中的每个列都要求有名称和数据类型

主键为`ID`(数据库表中对储存数据对象予以唯一和完整标识的数据列或属性的键,一个数据表只能有一个主键,且主键的取值不能缺失,即不能为空值`Null`)

|数据类型|描述|
|:---:|:---:|
|`CHARACTER(n)`|字符/字符串,固定长度n|
|`VARCHAR(n)`或`CHARACTER VARYING(n)`|字符/字符串,可变长度,最大长度n|
|`BINARY(n)`|二进制串,固定长度n|
|`BOOLEAN`|存储TRUE或FALSE值|
|`VARBINARY(n)`或BINARY VARYING(n)|二进制串,可变长度,最大长度n|
|`INTEGER(p)`|整数值(没有小数点),精度`p`|
|`SMALLINT`|整数值(没有小数点),精度`5`|
|`INTEGER`|整数值(没有小数点),精度`10`|
|`BIGINT`|整数值(没有小数点),精度`19`|
|`DECIMAL(p,s)`|精确数值,精度`p`,小数点后位数`s`,例如:`decimal(5,2)`是一个小数点前有3位数,小数点后有2位数的数字|
|`NUMERIC(p,s)`|精确数值,精度`p`,小数点后位数`s`,(与`DECIMAL`相同)|
|`FLOAT(p)`|近似数值,尾数精度`p`|
|`REAL`|近似数值,尾数精度`7`|
|`FLOAT`|近似数值,尾数精度`16`|
|`DOUBLE PRECISION`|近似数值,尾数精度`16`|
|`DATE`|存储年,月,日的值|
|`TIME`|存储小时,分,秒的值|
|`TIMESTAMP`|存储年,月,日,小时,分,秒的值|
|`INTERVAL`|由一些整数字段组成,代表一段时间,取决于区间的类型|
|`ARRAY`|元素的固定长度的有序集合|
|`MULTISET`|元素的可变长度的无序集合|
|`XML`|存储XML数据|

3. 删除表

`drop table students`

4. 更新表

在表中增加一列`alter table students add column Age INTEGER;`

在表中删除一列`alter table students drop column Age;`

重命名表`rename table students to student;`

## DML

1. 向表中插入数据

`insert into students (ID,Name,City) values (1,'Mike','GZ');`

`insert into students values(2,'Ben','NY');`

2. 查询数据

`select id from students;`

```
+----+
| id |
+----+
|  1 |
|  2 |
+----+
```

`select * from students;`

```
+----+------+------+
| ID | Name | City |
+----+------+------+
|  1 | Mike | GZ   |
|  2 | Ben  | NY   |
+----+------+------+
```

`select * from students where ID=1;`

```
+----+------+------+
| ID | Name | City |
+----+------+------+
|  1 | Mike | GZ   |
+----+------+------+
```

3. `where`子句

```
SELECT field1, field2,...fieldN FROM table_name1, table_name2...
[WHERE condition1 [AND [OR]] condition2.....
```

A为10,B为20

|操作符|描述|实例|
|:---:|:---:|:---:|
|`=`|等号,检测两个值是否相等,如果相等返回true|`A=B`返回false|
|`<>`或`!=`|不等于,检测两个值是否相等,如果不相等返回true|`A!=B`返回true|
|`>`|大于号,检测左边的值是否大于右边的值,如果左边的值大于右边的值返回true|`A>B`返回false|
|`<`|小于号,检测左边的值是否小于右边的值,如果左边的值小于右边的值返回true|`A<B`返回true|
|`>=`|大于等于号,检测左边的值是否大于或等于右边的值,如果左边的值大于或等于右边的值返回true|`A>=B`返回false|
|`<=`|小于等于号,检测左边的值是否小于于或等于右边的值,如果左边的值小于或等于右边的值返回true|`A<=B`返回true|

`select * from students where ID=1;`

```
+----+------+------+
| ID | Name | City |
+----+------+------+
|  1 | Mike | GZ   |
+----+------+------+
```

使用算数运算符`+-*/`

> 所有包含NULL 的计算,结果肯定是NULL

`select id*2 as 'id2' from students;`

```
+-----+
| id2 |
+-----+
|   2 |
|   4 |
+-----+
```

> 不能对NULL使用比较运算符,应该使用`where xxx is null;`

使用逻辑运算符`not and or`

`select * from students where not id>1;`

```
+----+------+------+
| ID | Name | City |
+----+------+------+
|  1 | Mike | GZ   |
+----+------+------+
```

`select * from students where id=1 and Name="Mike";`

```
+----+------+------+
| ID | Name | City |
+----+------+------+
|  1 | Mike | GZ   |
+----+------+------+
```

`select * from students where id=1 or 2;`

```
+----+------+------+
| ID | Name | City |
+----+------+------+
|  1 | Mike | GZ   |
|  2 | Ben  | NY   |
+----+------+------+
```

`select * from students where id=1 or Name="Ben";`

```
+----+------+------+
| ID | Name | City |
+----+------+------+
|  1 | Mike | GZ   |
|  2 | Ben  | NY   |
+----+------+------+
```

4. 删除数据

`DELETE FROM table_name [WHERE Clause]`

> 如果没有指定 WHERE 子句,MySQL表中的所有记录将被删除

5. 更新数据

`update students set Name="Sam" where id=1;`

```
+----+------+------+
| ID | Name | City |
+----+------+------+
|  1 | Sam  | GZ   |
|  2 | Ben  | NY   |
+----+------+------+
```

6. 聚合查询

|聚合函数|作用|
|:---:|:---:|
|COUNT|计算表中的记录数(行数)|
|SUM|计算表中数值列中数据的合计值|
|AVG|计算表中数值列中数据的平均值|
|MAX|求出表中任意列中数据的最大值|
|MIN|求出表中任意列中数据的最小值|

`select count(*) from students;`

```
+----------+
| count(*) |
+----------+
|        2 |
+----------+
```

7. 对结果排序

`select * from students order by id;`

```
+----+------+------+
| ID | Name | City |
+----+------+------+
|  1 | Sam  | GZ   |
|  2 | Ben  | NY   |
+----+------+------+
```

`select * from students order by id DESC;`(默认为升序`ASC`)

```
+----+------+------+
| ID | Name | City |
+----+------+------+
|  2 | Ben  | NY   |
|  1 | Sam  | GZ   |
+----+------+------+
```

8. 谓词

- LIKE谓词

前方一致查询:

`SELECT * FROM SampleLike WHERE strcol LIKE 'ddd%';`

也可用`_`代替`%`,但`_`只能代表一个字符

`SELECT * FROM SampleLike WHERE strcol LIKE 'abc_';`

中间一致查询:

`SELECT * FROM SampleLike WHERE strcol LIKE '%ddd%';`

后方一致查询:

`SELECT * FROM SampleLike WHERE strcol LIKE '%ddd';`

- BETWEEN谓词

`SELECT product_name, sale_price FROM Product WHERE sale_price BETWEEN 100 AND 1000;`

BETWEEN的特点就是结果中会包含100 和1000 这两个临界值

- IS NULL和IS NOT NULL谓词

为了选取出某些值为NULL 的列的数据,不能使用`=`,而只能使用特定的谓词`IS NULL`

`SELECT product_name, purchase_price FROM Product WHERE purchase_price IS NULL;`

- IN谓词

`SELECT product_name, purchase_price FROM Product WHERE purchase_price IN (320, 500, 5000);`

也可以用`NOT IN`

`SELECT product_name, purchase_price FROM Product WHERE purchase_price NOT IN (320, 500, 5000);`

>在使用`IN`和`NOT IN`时是无法选取出NULL 数据的

- EXIST谓词

```
SELECT product_name, sale_price
  FROM Product AS P 
 WHERE EXISTS (SELECT *
                  FROM ShopProduct AS SP 
                 WHERE SP.shop_id = '000C'
                   AND SP.product_id = P.product_id);
```

也可以用NOT EXIST

9. `UNION`

> UNION操作符用于连接两个以上的SELECT语句的结果组合到一个结果集合中,多个SELECT语句会删除重复的数据,如果想保留重复记录,可以在UNION后面加ALL

`SELECT product_id, product_name FROM Product UNION SELECT product_id, product_name FROM Product2;`

## DCL

1. 创建事务(START TRANSACTION) - 提交处理(COMMIT)

```text
START TRANSACTION;
    -- 将运动T恤的销售单价降低1000日元
    UPDATE Product
       SET sale_price = sale_price - 1000
     WHERE product_name = '运动T恤';
    -- 将T恤衫的销售单价上浮1000日元
    UPDATE Product
       SET sale_price = sale_price + 1000
     WHERE product_name = 'T恤衫';
COMMIT;
```

2. 取消处理(ROLLBACK)

```text
START TRANSACTION;
    -- 将运动T恤的销售单价降低1000日元
    UPDATE Product
       SET sale_price = sale_price - 1000
     WHERE product_name = '运动T恤';
    -- 将T恤衫的销售单价上浮1000日元
    UPDATE Product
       SET sale_price = sale_price + 1000
     WHERE product_name = 'T恤衫';
ROLLBACK;
```

---

## 注释

1. `#xxx`

2. `-- xxx`注意有一个空格

3. `/*xxx*/`