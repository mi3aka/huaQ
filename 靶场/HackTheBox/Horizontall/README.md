nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151336836.png)

要把`10.10.11.105 horizontall.htb`加到`/etc/hosts`中

---

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151337452.png)

没找到啥突破口,看了wp才知道,要爆破子域名...

[爆破字典 subdomains-top1million-110000.txt](https://raw.githubusercontent.com/danielmiessler/SecLists/master/Discovery/DNS/subdomains-top1million-110000.txt)

[爆破工具 gobuster](https://github.com/OJ/gobuster)

`./gobuster vhost -u http://horizontall.htb -w /mnt/hgfs/Exploits/subdomains-top1million-110000.txt -t 100`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151415128.png)

把`10.10.11.105 api-prod.horizontall.htb`加到`/etc/hosts`中,识别出`Strapi cms`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151418002.png)

[http://api-prod.horizontall.htb/admin/auth/login](http://api-prod.horizontall.htb/admin/auth/login)

登录页面

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151440114.png)

通过[http://api-prod.horizontall.htb/admin/strapiVersion](http://api-prod.horizontall.htb/admin/strapiVersion)得知版本号为`{"strapiVersion":"3.0.0-beta.17.4"}`

主要有以下两个exp

[Strapi CMS 3.0.0-beta.17.4 - Remote Code Execution (RCE) (Unauthenticated)](https://www.exploit-db.com/exploits/50239)

[Strapi 3.0.0-beta - Set Password (Unauthenticated)](https://www.exploit-db.com/exploits/50237)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151506683.png)

检测是否存在python

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151517814.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151518611.png)

利用python反弹shell(反弹shell不稳定...)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151541440.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151545704.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151550216.png)

在`/home/developer/user.txt`中得到`user flag`

在`~/myapi/config/environments/development/database.json`中得到数据库连接密码,尝试以此作为ssh密码(无法连接)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151606420.png)

后面不会了,看了眼wp,wp说生成一个ssh密钥???

在本地用`ssh-keygen`生成了一对密码,将公钥传到靶机上,并将`authorized_keys`设置为公钥的内容

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151649265.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151648884.png)

成功使用`ssh -i ~/.ssh/id_rsa strapi@horizontall.htb`连接到靶机

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151651256.png)

同样使用`python3 -c 'import pty;pty.spawn("/bin/bash")'`将其转换为可交互的shell

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151654790.png)

`netstat -ano`查看端口情况

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151702857.png)

有个`8000`端口,不知道是干什么的,不能直接访问,但可以通过SSH端口转发把8000映射到本地

`ssh -i ~/.ssh/id_rsa -L 8000:127.0.0.1:8000 strapi@horizontall.htb`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151728455.png)

xray扫描

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151730316.png)

[https://github.com/zhzyker/CVE-2021-3129](https://github.com/zhzyker/CVE-2021-3129)

在本地起个python服务器,把exp传到靶机上

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151808891.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151808231.png)

改了一下RCE5的内容,改成了`python3 /tmp/a.py &`,将shell反弹

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151817974.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151820250.png)