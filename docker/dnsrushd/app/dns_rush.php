<?php
function handle_result($id,$uuid) {
        $dt="/var/www/html/deploy/$id/node/$uuid/";
	if(file_exists($dt)) {
		if(isset($_FILES["result"]["tmp_name"])) {
			if (!move_uploaded_file($_FILES["result"]["tmp_name"], $dt . "result")) {
				$err="unable to access uploaded result file";
			}
		} else {
			$err="missing result data";
		}
		if($err!="") {
			file_put_contents($dt . "error.msg",$err);
			file_put_contents($dt . "error","oops");
		}
		echo "OK";
	}
}
function handle_error($id,$uuid,$info) {
        $dt="/var/www/html/deploy/$id/node/$uuid/";
	if(file_exists($dt)) {
		file_put_contents($dt . "error",$info);
		if($info=="download") {
			file_put_contents($dt . "error.msg","unable to dowanload playlist");
		}
		if($info=="verify") {
			file_put_contents($dt . "error.msg","unable to parse playlist");
		}
		if($info=="result") {
			$err="";
			if(isset($_FILES["message"]["tmp_name"])) {
				if (!move_uploaded_file($_FILES["message"]["tmp_name"], $dt . "error.msg")) {
					$err="unable to access uploaded message";
				}
			} else {
				$err="missing error description";
			}
			if($err!="") {
				file_put_contents($dt . "error.msg",$err);
			}
		}
	}
	echo "OK";
}
function handle_status($id,$uuid,$info) {
        $dt="/var/www/html/deploy/$id/node/$uuid/";
	if(file_exists($dt)) {
		file_put_contents($dt . "status",$info);
	}
}

function handle_query($id,$uuid) {
        $dt="/var/www/html/deploy/$id/";
	// global command
	if(file_exists($dt . "command")) {
		$cmd=trim(file_get_contents($dt . "command"));
		if($cmd=="quit") {
			die('quit');
		} else {
			list($cmd,$rest)=explode(" ",$cmd);
			if($cmd=="play") {
				$sched=file_get_contents($dt . "schedule");
				if(time()<=$sched) {
					$mode=file_get_contents($dt . "mode");
					$count=file_get_contents($dt . "count");
					$ns=file_get_contents($dt . "ns");
					die("play $sched $ns $count $mode");
				} else {
					// missing locking
					@unlink($dt . "command");
					@unlink($dt . "schedule");
					@unlink($dt . "count");
					@unlink($dt . "mode");
				}
			}
		}
		
	}
}
	if(!isset($_REQUEST['id'])) {
		die('-');
	}
	$id=$_REQUEST['id'];
	$dt="/var/www/html/deploy/$id/";
	if(!file_exists($dt)) {
		die('-');
	}
	if(isset($_REQUEST['register'])) {
		$ip=$_SERVER['REMOTE_ADDR'];
		if(isset($_REQUEST['hn'])) {
			$hn=$_REQUEST['hn'];
		} else {
			$hn="unknown";
		}
		$uuid = trim(file_get_contents('/proc/sys/kernel/random/uuid'));
		while(file_exists($dt . $uuid)) $uuid = trim(file_get_contents('/proc/sys/kernel/random/uuid'));
		$node = $dt . "/node/" .  $uuid;
		if(mkdir($node,0755,true)) {
			file_put_contents($node . "/ip",$ip);
			file_put_contents($node . "/hostname",$hn);
			die($uuid);
		}
		die('-');
	}
	if(isset($_REQUEST['action'])) {
		$action = $_REQUEST['action'];
		$uuid="";
		if(isset($_REQUEST['uuid'])) {
			$uuid=$_REQUEST['uuid'];
		}
		if($uuid=="") {
			die("-");
		}
		if(!file_exists($dt . "/node/" . $uuid . "/")) {
			die("-");
		}	
		$node=$dt . "/node/" . $uuid . "/";
		file_put_contents($node . "ping","OK");
		switch($action) {
			case "query":
				handle_query($id,$uuid);
				break;
			case "status":
				$info="";
				if(isset($_REQUEST['info'])) {
					$info = $_REQUEST['info'];
				}
				if(strlen($info)) {
					handle_status($id,$uuid,$info);
				}
				break;
			case "error":
				$info="";
				if(isset($_REQUEST['info'])) {
					$info = $_REQUEST['info'];
				}
				handle_error($id,$uuid,$info);
				break;
			case "result":
				handle_result($id,$uuid);
				break;
		}
	}
