[https://www.vulnhub.com/entry/beelzebub-1,742/](https://www.vulnhub.com/entry/beelzebub-1,742/)

`dirsearch`显示`index.php`返回`200`,但是在浏览器打开`index.php`返回`404 NOT FOUND`

查看网页源码显示

```html
<html><head>
<title>404 Not Found</title>
</head><body>
<h1>Not Found</h1>
<!--My heart was encrypted, "beelzebub" somehow hacked and decoded it.-md5-->
<p>The requested URL was not found on this server.</p>
<hr>
<address>Apache/2.4.30 (Ubuntu)</address>
</body></html>
```

对`beelzebub`计算md5 `echo -n beelzebub | md5sum`

返回`d18e1e22becbd915b45e0e655429d487`

尝试以`beelzebub`作为用户名,`d18e1e22becbd915b45e0e655429d487`作为密码进行ssh连接但是ssh连接失败

对路径`d18e1e22becbd915b45e0e655429d487`的扫描结果如下

```
[23:53:07] 403 -  278B  - /d18e1e22becbd915b45e0e655429d487//.htaccess     
[23:53:07] 403 -  278B  - /d18e1e22becbd915b45e0e655429d487//.htpasswd     
[23:53:07] 301 -    0B  - /d18e1e22becbd915b45e0e655429d487//?m=a  ->  http://192.168.148.7/d18e1e22becbd915b45e0e655429d487/?m=a
[23:53:08] 200 -   58KB - /d18e1e22becbd915b45e0e655429d487//?s=d          
[23:53:08] 200 -   56KB - /d18e1e22becbd915b45e0e655429d487//?wp-html-rend 
[23:53:08] 200 -   56KB - /d18e1e22becbd915b45e0e655429d487//?pageservices 
[23:53:12] 200 -   56KB - /d18e1e22becbd915b45e0e655429d487//index.php/    
[23:53:12] 200 -   56KB - /d18e1e22becbd915b45e0e655429d487//index.php?chemin=..%2f..%2f..%2f..%2f..%2f..%2f%2fetc
[23:53:12] 200 -   56KB - /d18e1e22becbd915b45e0e655429d487//index.php?file=../../../../../../etc/passwd
[23:53:12] 200 -   56KB - /d18e1e22becbd915b45e0e655429d487//index.php?file=/etc/passwd
[23:53:12] 200 -   56KB - /d18e1e22becbd915b45e0e655429d487//index.php?page=../../../../etc/passwd
[23:53:18] 200 -   56KB - /d18e1e22becbd915b45e0e655429d487/index.php      
[23:53:27] 200 -   45KB - /d18e1e22becbd915b45e0e655429d487/wp-includes/    
[23:53:27] 200 -    7KB - /d18e1e22becbd915b45e0e655429d487/readme.html     
[23:53:27] 200 -   19KB - /d18e1e22becbd915b45e0e655429d487/license.txt     
[23:53:41] 302 -    0B  - /d18e1e22becbd915b45e0e655429d487/wp-admin/  ->  http://192.168.1.6/d18e1e22becbd915b45e0e655429d487/wp-login.php?redirect_to=http%3A%2F%2F192.168.148.7%2Fd18e1e22becbd915b45e0e655429d487%2Fwp-admin%2F&reauth=1
[23:53:44] 301 -  350B  - /d18e1e22becbd915b45e0e655429d487/wp-admin  ->  http://192.168.148.7/d18e1e22becbd915b45e0e655429d487/wp-admin/
[23:53:57] 200 -    0B  - /d18e1e22becbd915b45e0e655429d487/wp-config.php   
[23:53:57] 301 -  352B  - /d18e1e22becbd915b45e0e655429d487/wp-content  ->  http://192.168.148.7/d18e1e22becbd915b45e0e655429d487/wp-content/
[23:53:57] 200 -    0B  - /d18e1e22becbd915b45e0e655429d487/wp-content/
[23:53:57] 301 -  353B  - /d18e1e22becbd915b45e0e655429d487/wp-includes  ->  http://192.168.148.7/d18e1e22becbd915b45e0e655429d487/wp-includes/
[23:53:57] 200 -    6KB - /d18e1e22becbd915b45e0e655429d487/wp-login.php  
```

显示其可能为`wordpress`的后台,使用[wpscan](https://github.com/wpscanteam/wpscan)对其进行扫描,检测其是否存在漏洞

`docker run -it --rm wpscanteam/wpscan --url http://192.168.148.7/d18e1e22becbd915b45e0e655429d487/ -e --plugins-detection aggressive --ignore-main-redirect --force vp`

```
[+] XML-RPC seems to be enabled: http://192.168.148.7/d18e1e22becbd915b45e0e655429d487/xmlrpc.php
 | Found By: Direct Access (Aggressive Detection)
 | Confidence: 100%
 | References:
 |  - http://codex.wordpress.org/XML-RPC_Pingback_API
 |  - https://www.rapid7.com/db/modules/auxiliary/scanner/http/wordpress_ghost_scanner/
 |  - https://www.rapid7.com/db/modules/auxiliary/dos/http/wordpress_xmlrpc_dos/
 |  - https://www.rapid7.com/db/modules/auxiliary/scanner/http/wordpress_xmlrpc_login/
 |  - https://www.rapid7.com/db/modules/auxiliary/scanner/http/wordpress_pingback_access/

...

[+] valak
 | Found By: Author Id Brute Forcing - Author Pattern (Aggressive Detection)
 | Confirmed By: Login Error Messages (Aggressive Detection)

[+] krampus
 | Found By: Author Id Brute Forcing - Author Pattern (Aggressive Detection)
 | Confirmed By: Login Error Messages (Aggressive Detection)
```

`dirsearch`显示存在`/d18e1e22becbd915b45e0e655429d487/wp-content/uploads/`路径

`http://192.168.148.7/d18e1e22becbd915b45e0e655429d487/wp-content/uploads/Talk%20To%20VALAK/`存在一个特殊的网页

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111010049029.png)

传入数据后,发现返回的数据包中含有password,`M4k3Ad3a1`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111011245699.png)

测试得知该密码为ssh密码,对应用户为`krampus`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111011248435.png)

用[linux-exploit-suggester](https://github.com/mzet-/linux-exploit-suggester)检测,发现存在CVE-2021-3156

使用该exp [https://github.com/worawit/CVE-2021-3156](https://github.com/worawit/CVE-2021-3156) 对CVE-2021-3156进行利用

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111011428826.png)

而在writeup中还提到可以使用`[CVE-2019-12181] Serv-U FTP Server`进行提权

[https://github.com/guywhataguy/CVE-2019-12181](https://github.com/guywhataguy/CVE-2019-12181)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202111011438511.png)