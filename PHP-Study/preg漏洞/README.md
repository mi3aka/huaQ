> `preg_replace`执行一个正则表达式的搜索和替换

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

