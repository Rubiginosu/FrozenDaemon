<?php
/**
 * Created by PhpStorm.
 * User: axoford12
 * Date: 17-8-4
 * Time: 下午10:15
 */
include "FrozenGo.php";
$fg = new FrozenGo("127.0.0.1",52023,"987IsA1MWIa68yQcvMfy");
print_r($fg->keyRegister("12345",0));
//print_r($fg->getServerList());
//print_r($fg->createServer("test",1));
//print_r($fg->getServerList());
//print_r($fg->inputToServer(0,"stop"));
//print_r($fg->getServerList());
