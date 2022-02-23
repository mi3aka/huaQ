![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202201823567.png)

strcmp函数无法处理数组,如果传入的其中一个变量是数组的话,函数会报错返回NULL,结合弱类型比较,可以绕过某些限制

>例子

```php
<?php
error_reporting(0);
$str1 = 'a';
var_dump(strcmp($str1, $str1));
if (strcmp($str1, $str1) == 0) {
    echo "wow";
}
$array1 = array('a', 'b');
var_dump(strcmp($str1, $array1));
if (strcmp($str1, $array1) == 0) {
    echo "wow";
}
$array2 = array('a', 'b');
var_dump(strcmp($array1, $array2));
if (strcmp($array1, $array2) == 0) {
    echo "wow";
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202201829727.png)