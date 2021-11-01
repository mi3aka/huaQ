



`index.php`提示进行sql注入,但使用万能密码进行初步尝试,返回`Try Again`,推测闭合方式不正确或者存在过滤

`dirsearch`扫描得到的主要路径

```
http://192.168.148.9:80/add.php
http://192.168.148.9:80/c.php
http://192.168.148.9:80/index.php
http://192.168.148.9:80/show.php
http://192.168.148.9:80/test.php
http://192.168.148.9:80/head.php
http://192.168.148.9:80/in
http://192.168.148.9/phpmy
```

`http://192.168.148.9/phpmy`跳转到`phpmyadmin`的登录页面

打开`test.php`提示要向其传递`file`参数,推测其可能为文件包含或文件下载,尝试GET传参失败,使用POST传参成功

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111011728795.png)

`index.php`关键内容如下

```php
$uname=str_replace('\'','',urldecode($_POST['un']));
$pass=str_replace('\'','',urldecode($_POST['ps']));
$run='select * from auth where  pass=\''.$pass.'\' and uname=\''.$uname.'\'';
```

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

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111011740107.png)

`c.php`中含有数据库连接密码`$conn = mysqli_connect("127.0.0.1","billu","b0x_billu","ica_lab");`

尝试利用该密码进行登录和ssh连接均失败,使用该密码成功登录`phpmyadmin`

