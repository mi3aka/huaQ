<?php

ini_set("open_basedir", getcwd() . ":/etc:/tmp");

$DATABASE = array(

    "host" => "127.0.0.1",
    "username" => "root",
    "password" => "root",
    "dbname" =>"ctfusers"
);

$db = new mysqli($DATABASE['host'],$DATABASE['username'],$DATABASE['password'],$DATABASE['dbname']);