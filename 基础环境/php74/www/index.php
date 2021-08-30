<?php
$world = "World";
echo "Hello, ", $world . "<br>";

new mysqli("172.16.172.202", "root", "root", "mysql");
if (mysqli_connect_errno()) { #检查连接
    printf("Connect failed: %s\n", mysqli_connect_error());
    exit();
}
echo "连接成功";
if (!is_dir('test_mkdir')) {
    mkdir('test_mkdir');
}
phpinfo();
?>