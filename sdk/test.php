<?php
/**
 * Created by PhpStorm.
 * User: axoford12
 * Date: 17-8-4
 * Time: 下午10:15
 */
include "FrozenGo.php";
$fg = new FrozenGo("127.0.0.1", 52023, "gqf5JeYiGHY4ew5qF2biOWxnZORdYs0wlBepYSelG6sNlUTxHttQvOLjiwdDDEth");
#print_r($fg->keyRegister("112233",0));
//print_r($fg->getServerList());
#print_r($fg->createServer("test",1));
//print_r($fg->setServerConfig(0, [[
//    "AttrName" => "MaxCpuRate", "AttrValue" => "20"],
//    ["AttrName" => "MaxHardDiskCapacity", "AttrValue" => "10240"],
//    ["AttrName" => "Executable", "AttrValue" => "Minecraft"],
//    ["AttrName" => "MaxMem", "AttrValue" => "120"],
//    ["AttrName" => "Expire", "AttrValue" => "100000"]]));
#print_r($fg->getServerList());
//echo json_encode([["AttrName"=>"MaxHardDiskCapacity","AttrValue"=>"1000000"],["AttrName"=>"Executable","AttrValue"=>"Minecraft"]]);
print_r($fg->startServer(0));
//print_r($fg->inputToServer(0,"eng"));
//print_r($fg->getServerList());
//print_r($fg->stopServer(1));
