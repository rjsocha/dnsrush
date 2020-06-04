<?php
if($_SERVER['REMOTE_ADDR'] == "127.0.0.1") {
	$data="/var/www/html/deploy/";
	$uuid = trim(file_get_contents('/proc/sys/kernel/random/uuid'));
	while(file_exists($data . $uuid)) $uuid = trim(file_get_contents('/proc/sys/kernel/random/uuid'));
	$prj = $data . $uuid;
	if(mkdir($prj,0755,true))
		die($uuid);
}
die();
