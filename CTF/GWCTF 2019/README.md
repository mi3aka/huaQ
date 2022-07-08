## 枯燥的抽奖

> mt_srand爆破

[https://www.openwall.com/php_mt_seed/](https://www.openwall.com/php_mt_seed/)下载4.0版

```php
<?php
#这不是抽奖程序的源代码！不许看！
header("Content-Type: text/html;charset=utf-8");
session_start();
if(!isset($_SESSION['seed'])){
$_SESSION['seed']=rand(0,999999999);
}

mt_srand($_SESSION['seed']);
$str_long1 = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ";
$str='';
$len1=20;
for ( $i = 0; $i < $len1; $i++ ){
    $str.=substr($str_long1, mt_rand(0, strlen($str_long1) - 1), 1);       
}
$str_show = substr($str, 0, 10);
echo "<p id='p1'>".$str_show."</p>";


if(isset($_POST['num'])){
    if($_POST['num']===$str){x
        echo "<p id=flag>抽奖，就是那么枯燥且无味，给你flag{xxxxxxxxx}</p>";
    }
    else{
        echo "<p id=flag>没抽中哦，再试试吧</p>";
    }
}
show_source("check.php"); 
```

跟新生赛那题一毛一样

```python
import os

s="07UW6BLYyB"
randstr="abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"#pay attention

seed=""

for i in range(len(s)):
    pos=randstr.index(s[i])
    seed+="%s %s %s %s "%(pos,pos,0,len(randstr)-1)

os.system("./php_mt_seed "+seed)
```

```php
<?php
$str = '';
$all_alpha = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ";#pay attention

mt_srand(630868830);

$flag_len = 20;#pay attention
for ( $i = 0; $i < $flag_len; $i++ ){
    $str.=substr($all_alpha, mt_rand(0, strlen($all_alpha) - 1), 1);       
}
echo $str
?>
```

