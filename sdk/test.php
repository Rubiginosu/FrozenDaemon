<?php
/**
 * Created by IntelliJ IDEA.
 * User: axoford12
 * Date: 8/19/17
 * Time: 10:57 PM
 */
include "FrozenGo.php";
$fg = new FrozenGo("127.0.0.1","52023","B5Inr7RDZ6QfVbBKQgc6MqMqUqUClGScK75HavHmLU9yiXhVQsyzWBuAzRb6r79z");
#print_r($fg->createServer("Test",0));
//print_r($fg->setServerConfig(0,[["AttrName" => "Executable","AttrValue" => "Minecraft"]]));
print_r($fg->startServer(0));