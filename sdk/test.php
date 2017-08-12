<?php
/**
 * Created by PhpStorm.
 * User: axoford12
 * Date: 17-8-4
 * Time: 下午10:15
 */
include "FrozenGo.php";
$fg = new FrozenGo("127.0.0.1",52023,"hvGpglPuYs9xTfmgrTKQYlqzkBHlvMhJgwFr7PMwMZ6yz4LuhBnyXRk1HyMgaN07");
#print_r($fg->keyRegister("112233",0));
#print_r($fg->getServerList());
print_r($fg->createServer("test",1));
print_r($fg->setServerConfig(1,[["AttrName"=>"MaxHardDisk","AttrValue"=>"10240"],["AttrName"=>"Executable","AttrValue"=>"Minecraft"]]));
print_r($fg->getServerList());
print_r($fg->startServer(1));
//print_r($fg->inputToServer(0,"stop"));
//print_r($fg->getServerList());
