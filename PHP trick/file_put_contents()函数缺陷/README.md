![](2021-12-06%2009-00-53%20的屏幕截图.png)

`file_put_contents`可以使用数组进行传参

```php
<?php
highlight_file(__FILE__);
$a = array('a', 'b', 'c');
file_put_contents('a', $a);#file_put_contents可以使用数组作为参数,从而绕过某些过滤函数的限制
var_dump(file_get_contents('a'));
```

![](2021-12-01%2009-41-32%20的屏幕截图.png)