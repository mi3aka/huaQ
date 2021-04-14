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

在表中增加一列