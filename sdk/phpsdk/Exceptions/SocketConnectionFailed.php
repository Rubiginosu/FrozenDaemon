<?php
/**
 * Created by PhpStorm.
 * User: seth8277
 * Date: 17-8-24
 * Time: 下午9:34
 */

namespace Seth8277\FrozenSdk\Exceptions;


use InvalidArgumentException;

class SocketConnectionFailed extends InvalidArgumentException
{
    public static function create($ip,$port,$reason){
        return new static("Unable to connect to Server($ip:$port),reason: $reason");
    }

}