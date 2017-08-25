<?php
/**
 * Created by PhpStorm.
 * User: seth8277
 * Date: 17-8-24
 * Time: 下午9:26
 */

namespace FrozenSdk\Exceptions;


use InvalidArgumentException;

class SocketCreationFailed extends InvalidArgumentException
{
    public static function create(){
        return new static("Coundn't create socket!");
    }

}