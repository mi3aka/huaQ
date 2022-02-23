[$_REQUEST — HTTP Request 变量](https://www.php.net/manual/zh/reserved.variables.request.php)

`$_REQUEST`中包含了`$_GET,$_POST,$_COOKIE`中的值,但要注意的是`$_REQUEST`的调用方式是**复制**这三个数组中的值**而非引用**

对这三个数组中的值进行修改不会影响到`$_REQUEST`数组

```php
<?php
var_dump($_GET);
var_dump($_POST);
var_dump($_COOKIE);
var_dump($_REQUEST);
$_GET = array_map('md5', $_GET);
$_POST = array_map('md5', $_POST);
var_dump($_GET);
var_dump($_POST);
var_dump($_COOKIE);
var_dump($_REQUEST);
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202191558608.png)

---

phpinfo中注意两个值,一个是`request_order`,另一个是`variables_order`

`request_order`设置了PHP将GET,POST和Cookie中的值复制到REQUEST的顺序,方向为从左到右完成,新的值会覆盖
旧的值,如果没有设置`request_order`则会根据`variables_order`

`variables_order`默认设置为`EGPCS`即启用了`$_ENV,$_GET,$_POST,$_COOKIE,$_SERVER`

如果`variables_order`设置为`SP`那么`PHP`将启用`$_SERVER,$_POST`,但不启用`$_ENV,$_GET,$_COOKIE`

```php
<?php
var_dump($_GET);
var_dump($_POST);
var_dump($_REQUEST);
```

1. 默认情况

`request_order`为`no value`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202191625078.png)

2. `php.ini`设置`request_order = "PG"`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202201636936.png)

3. `php.ini`设置`request_order = "GP"`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202201637266.png)