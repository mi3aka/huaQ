![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202261946613.png)

当url中有多个`@`符号时,`parse_url`匹配最后一个`@`后面符合格式的`host`,而`libcurl`则是匹配第一个@后面符合格式的`host`

```php
<?php
$host="http://foo@evil.com@google.com/";
var_dump(parse_url($host));
?>
```

```
array (size=4)
  'scheme' => string 'http' (length=4)
  'host' => string 'google.com' (length=10)
  'user' => string 'foo@evil.com' (length=12)
  'path' => string '/' (length=1)
```

```cpp
#include <stdio.h>
#include <curl/curl.h>
int main()
{
    CURLU *h;
    CURLUcode uc;
    char *host;
    char *path;

    h = curl_url();
    uc = curl_url_set(h, CURLUPART_URL, "http://foo@evil.com@google.com/", 0);
    uc = curl_url_get(h, CURLUPART_HOST, &host, 0);
    printf("Host name: %s\n", host);
    curl_free(host);

    uc = curl_url_get(h, CURLUPART_PATH, &path, 0);
    printf("Path: %s\n", path);
    curl_free(path);
    return 0;
}
```

```
Host name: evil.com@google.com
Path: /
```

---

```php
<?php
$host='http://foo@evil.com@google.com/';
var_dump(filter_var($host, FILTER_VALIDATE_URL));
var_dump(parse_url($host));
?>
```

```
boolean false
array (size=4)
  'scheme' => string 'http' (length=4)
  'host' => string 'google.com' (length=10)
  'user' => string 'foo@evil.com' (length=12)
  'path' => string '/' (length=1)
```

可见`http://foo@evil.com@google.com/`无法通过`filter_var`的检测

但是URL具有许多特殊字符例如`@#;:\,`等,可以利用这些特殊字符去构造特殊的URL从而绕过`filter_var`

```php
<?php
$host='http://evil.com;google.com';
var_dump(filter_var($host, FILTER_VALIDATE_URL));
var_dump(parse_url($host));
$host='0://evil.com;google.com';
var_dump(filter_var($host, FILTER_VALIDATE_URL));
var_dump(parse_url($host));
var_dump(null);
$host='http://evil.com,google.com';
var_dump(filter_var($host, FILTER_VALIDATE_URL));
var_dump(parse_url($host));
$host='0://evil.com,google.com';
var_dump(filter_var($host, FILTER_VALIDATE_URL));
var_dump(parse_url($host));
var_dump(null);
$host='http://evil.com\google.com';
var_dump(filter_var($host, FILTER_VALIDATE_URL));
var_dump(parse_url($host));
$host='0://evil.com\google.com';
var_dump(filter_var($host, FILTER_VALIDATE_URL));
var_dump(parse_url($host));
var_dump(null);
$host='http://evil$google.com';
var_dump(filter_var($host, FILTER_VALIDATE_URL));
var_dump(parse_url($host));
$host='0://evil$google.com';
var_dump(filter_var($host, FILTER_VALIDATE_URL));
var_dump(parse_url($host));
?>
```

```
boolean false
array (size=2)
  'scheme' => string 'http' (length=4)
  'host' => string 'evil.com;google.com' (length=19)
string '0://evil.com;google.com' (length=23)
array (size=2)
  'scheme' => string '0' (length=1)
  'host' => string 'evil.com;google.com' (length=19)
null
boolean false
array (size=2)
  'scheme' => string 'http' (length=4)
  'host' => string 'evil.com,google.com' (length=19)
string '0://evil.com,google.com' (length=23)
array (size=2)
  'scheme' => string '0' (length=1)
  'host' => string 'evil.com,google.com' (length=19)
null
boolean false
array (size=2)
  'scheme' => string 'http' (length=4)
  'host' => string 'evil.com\google.com' (length=19)
string '0://evil.com\google.com' (length=23)
array (size=2)
  'scheme' => string '0' (length=1)
  'host' => string 'evil.com\google.com' (length=19)
null
boolean false
array (size=2)
  'scheme' => string 'http' (length=4)
  'host' => string 'evil$google.com' (length=15)
string '0://evil$google.com' (length=19)
array (size=2)
  'scheme' => string '0' (length=1)
  'host' => string 'evil$google.com' (length=15)
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202262114676.png)