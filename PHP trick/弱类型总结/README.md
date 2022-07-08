php中有两种比较的符号`==`与`===`

`==`在进行比较的时候,会把两端变量类型转换成相同的,再进行比较(一个字符串与一个数字相比较时,字符串会转换成数值)

`===`在进行比较的时候,会先判断两种字符串的类型是否相等,再进行比较

```php
<?php
$a=1;
$b="a1";
$c="1a";
var_dump($a==$b);//$b被转换成0
var_dump($a==$c);
var_dump($a===$b);
var_dump($a===$c);
?>
```

```
boolean false
boolean true
boolean false
boolean false
```

`0e`开头且都是数字的字符串,**弱类型比较**都等于0

```php
<?php
var_dump(0=='1a');
var_dump(0=='a1');
var_dump(0==0e1);
var_dump(0=="0e1");
var_dump(0===0e1);
var_dump(0==="0e1");
?>
```

```
boolean false
boolean true
boolean true
boolean true
boolean false
boolean false
```

字符串转数值时,调用了`intval()`取整函数,而这个函数在转换字符串的时候即使碰到不能转换的字符串的时候它也不会报错,而是返回0

```php
<?php
var_dump(intval(4));
var_dump(intval("10asd"));
var_dump(intval("asd10"));
var_dump(intval(1.23));
?>
```

```
int 4
int 10
int 0
int 1
```

```php
<?php
var_dump(1.23=="1.23");
var_dump(15=="0xf");
?>
```

```
boolean true
boolean false
```