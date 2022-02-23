函数定义`intval(mixed $value, int $base = 10): int`

>如果 base 是 0，通过检测 value 的格式来决定使用的进制：
>
>如果字符串包括了 "0x" (或 "0X") 的前缀，使用 16 进制 (hex)；否则，
>
>如果字符串以 "0" 开始，使用 8 进制(octal)；否则，
>
>将使用 10 进制 (decimal)。


>除非value是一个字符串,否则base不会起作用

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
var_dump(intval('0xa',0));
var_dump(intval("0xa.0",0));
var_dump(intval(" 0xa",0));
```

```
int 10
int 10
int 10
```

同时`intval`不能正常处理数组,当数组为空时返回0,不为空则返回1

```php
<?php
$num=array();
var_dump(intval($num));
```

```
int 0
```

```php
<?php
$num=array('123','456');
var_dump($num);
var_dump(intval($num));
```

```
array (size=2)
  0 => string '123' (length=3)
  1 => string '456' (length=3)
int 1
```

```php
<?php
$num=array('foo','bar');
var_dump($num);
var_dump(intval($num));
```

```
array (size=2)
  0 => string 'foo' (length=3)
  1 => string 'bar' (length=3)
int 1
```