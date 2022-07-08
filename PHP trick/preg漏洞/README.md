## preg_replace

`preg_replace`执行一个正则表达式的搜索和替换

`preg_replace ( mixed $pattern , mixed $replacement , mixed $subject , int $limit = -1 , int &$count = ? ) : mixed`

搜索 subject 中匹配 pattern 的部分,以 replacement 进行替换

### 双写绕过

```php
<?php
    $ext="pphphp";
    $ext=preg_replace('/php/i','',$ext);
    var_dump($ext);//php
?>
```

> 仅将关键词替换,用嵌套双写即可绕过

```php
<?php
    $ext="pphphp";
    $ext=preg_replace('/(.*)p(.*)h(.*)p/i','',$ext);
    var_dump($ext);
?>
```

### 大小写绕过

```php
<?php
    $ext="phP";
    $ext=preg_replace('/php/','',$ext);
    var_dump($ext);//phP
?>
```

修饰符`i`会使模式中的字母进行**大小写不敏感**匹配

### 命令执行

修饰符`e`会将替换后的字符串作为php代码评估执行(eval 函数方式),并使用执行结果作为实际参与替换的字符串,仅 `preg_replace()`使用此修饰符

```php
<?php
    $a="/a/e";
    $b="phpinfo()";
    $c="a";
    preg_replace($a,$b,$c);//执行phpinfo()
?>
```

>todo



## preg_match

>preg_match()返回pattern的匹配次数。 它的值将是0次（不匹配）或1次，因为preg_match()在第一次匹配后 将会停止搜索。preg_match_all()不同于此，它会一直搜索subject 直到到达结尾
>
>如果发生错误preg_match()返回false

### 数组绕过

preg_match无法处理数组会产生错误,从而返回false,进而可以绕过某些限制

```php
<?php
$a=array(1,2,3);
if(preg_match("/[0-9]/", $a)) {
    echo "no";
}else {
    echo "yes";
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202231951467.png)

>todo tofinish