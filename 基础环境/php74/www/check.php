<?php
highlight_file(__FILE__);
$world = "World";
echo "Hello, ", $world . "<br>";

new mysqli("172.17.0.1", "root", "root", "mysql","4000");
if (mysqli_connect_errno()) { #检查连接
    printf("Connect failed: %s\n", mysqli_connect_error());
    exit();
}
echo "连接成功";
if (!is_dir('test_mkdir')) {
    mkdir('test_mkdir');
    rmdir('test_mkdir');
}
phpinfo(3);