# 文件上传

>部分内容可能会和upload-labs里面的内容重复

## 无校验

传就完事了,但要注意上传的路径能不能解析该文件(把webshell传上去却发现没有解析...)

有时候一些上传点可以控制上传路径,可以尝试能不能进行路径穿越,把文件上传到其他目录

例如`C:\Users\用户名\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`,传个vbs上去进行启动项提权

## 前端校验

burpsuite改就完事了

## 后端校验

### MIME检查

要求`$_FILES['upload_file']['type']==='xxx'`,burpsuite改`content-type`就完事了

### 文件后缀检查

#### 黑名单

1. 黑名单是否完整,比如过滤`php`但没有过滤`phtml`

2. 大小写绕过,比如`PhP`

3. 如果是将黑名单后缀替换为空,可以尝试双写绕过,`pphphp`替换后得到`php`

4. windows平台特性

```
.php(空格)

.php.

.php::$DATA

>tofinish
```

>todo


#### 白名单

### 文件内容检查

#### 文件头检查

