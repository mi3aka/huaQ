>服务端接收了用户的恶意输入以后,未经任何处理就将其作为Web应用模板内容的一部分,模板引擎在进行目标编译渲染的过程中,执行了用户插入的可以破坏模板的语句

![img](1344396-20200911174631687-758048107.png)

## PHP-twig

> 1.x

```php
<?php
require_once(dirname(__FILE__) . '/Twig-1.44.1/lib/Twig/Autoloader.php');
Twig_Autoloader::register(true);
$twig = new Twig_Environment(new Twig_Loader_String());
$output = $twig->render("Hello {{name}}", array("name" => $_GET["name"]));//将name作为模版变量的值
echo $output;
?>
```

```php
<?php
require_once(dirname(__FILE__) . '/Twig-1.44.1/lib/Twig/Autoloader.php');
Twig_Autoloader::register(true);
$twig = new Twig_Environment(new Twig_Loader_String());
$output = $twig->render("Hello {$_GET['name']}");//将name作为模版变量的值
echo $output;
?>
```

```
?name={{2*3}}
Hello 6
```

