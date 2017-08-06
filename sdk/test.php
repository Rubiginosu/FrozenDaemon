<?php
/**
 * Created by PhpStorm.
 * User: axoford12
 * Date: 17-8-4
 * Time: 下午10:15
 */
include "FrozenGo.php";
$fg = new FrozenGo("47.92.90.152",52023,"987IsA1MWIa68yQcvMfy");
print_r($fg->startServer(0));
#print_r($fg->getServerList());
#print_r($fg->inputToServer(0,"stop"));
//print_r($fg->getServerList());
