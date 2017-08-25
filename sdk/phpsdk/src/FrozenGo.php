<?php
/**
 * Created by PhpStorm.
 * User: seth8277
 * Date: 17-8-24
 * Time: 下午8:55
 */

namespace FrozenSdk;

use Illuminate\Support\Collection;
use FrozenSdk\Exceptions\SocketConnectionFailed;
use FrozenSdk\Exceptions\SocketCreationFailed;

class FrozenGo
{
    private $ip, $port, $verifyCode;

    public function __construct($ip, $port, $verifyCode)
    {
        $this->setIp($ip);
        $this->setPort($port);
        $this->setVerifyCode($verifyCode);
    }


    /**
     * 获取所有服务器与配置信息
     *
     * @return Collection
     */
    public function getServerList()
    {
        return $this->sendMethod('List');

    }


    /**
     * 创建服务器
     *
     * @param int $id 服务器ID （请确保唯一性）
     * @param string $name
     *
     * @return Collection
     * 成功时 Message 返回 OK ,否则返回错误信息
     */
    public function createServer($id = 0, $name = '')
    {
        return $this->sendMethod('Create', $id, $name);
    }


    /**
     * 删除服务器
     *
     * @param int $id 服务器ID
     *
     * @return Collection
     * 成功时 Message 返回 OK ,否则返回错误信息
     */
    public function deleteServer($id)
    {
        return $this->sendMethod("Delete", $id);
    }


    /**
     * 在Daemon上安装一个配置文件
     *
     * @param string $name 配置文件名
     * @return Collection
     */

    public function execInstall($name)
    {
        return $this->sendMethod("ExecInstall", 0, $name);

    }


    /**
     * 启动一个服务器
     *
     * @param int $id 服务器ID
     * @return Collection
     */
    public function startServer($id)
    {
        return $this->sendMethod("Start", $id);
    }


    /**
     * 停止一个服务器
     *
     * @param int $id 服务器ID
     * @return Collection
     */
    public function stopServer($id)
    {
        return $this->sendMethod("Stop", $id);
    }

    /**
     * 设置一个服务器
     *
     * @param int $id 服务器iD
     * @param Collection|array|json $elements 配置
     * Example:
     * [
     *      [
     *          "AttrName" => FrozenGo::SERVER_NAME,
     *          "AttrValue" => "Axoford12"
     *      ]
     * ]
     *
     * @return Collection
     * 成功 Message 返回更新的配置项的数目,否则返回错误信息
     */
    public function setServerConfig($id, $elements)
    {
        return $this->sendMethod("SetServerConfig", $id, self::toJson($elements));
    }

    /**
     * 获取一个服务器的配置
     *
     * @param int $id 服务器ID
     * @return Collection
     */
    public function getServerConfig($id)
    {
        $result = $this->sendMethod("GetServerConfig", $id);
        if ($result->has('Message'))
            $result->put('Message',json_decode($result->get('Message')));
        return $result;
    }

    /**
     * 向服务器发送命令
     *
     * @param int $id 服务器ID
     * @param string $command 命令
     * @return Collection
     */
    public function inputToServer($id, $command)
    {
        return $this->sendMethod("InputLineToServer", $id, $command . "\n");
    }

    /**
     * 获取在线玩家
     * @param int $id 服务器ID
     * @return Collection
     */
    public function getOnlinePlayers($id)
    {
        $result = $this->sendMethod("GetServerPlayers", $id);
        if ($result->has('Message'))
            $result->put('Message',json_decode($result->get('Message')));
        return $result;
    }

    /**
     * 删除服务器数据目录下的文件
     * @param int $id 服务器ID
     * @param String $file 文件路径
     * Example:
     *  %FGO%/servers/server$id/serverData/$file
     *
     * @return Collection
     */
    public function deleteServerFile($id, $file)
    {
        return $this->sendMethod("DeleteServerFile", $id, $file);
    }

    /**
     * 获取安装的配置文件列表
     *
     * @return Collection
     */
    public function getExecList()
    {
        $result = $this->sendMethod("GetExecList");
        if ($result->has('Message'))
            $result->put('Message',json_decode($result->get('Message')));
        return $result;
    }


    /*
     * Setters & Getters
     * */


    /**
     * @return string
     */
    public function getIp()
    {
        return $this->ip;
    }

    /**
     * @param string $ip
     */
    public function setIp($ip)
    {
        $this->ip = $ip;
    }

    /**
     * @return int
     */
    public function getPort()
    {
        return $this->port;
    }

    /**
     * @param int $port
     */
    public function setPort($port)
    {
        $this->port = $port;
    }

    /**
     * @return string
     */
    public function getVerifyCode()
    {
        return $this->verifyCode;
    }

    /**
     * @param string $verifyCode
     */
    public function setVerifyCode($verifyCode)
    {
        $this->verifyCode = $verifyCode;
    }


    /*
     * Protected Functions
     * */
    protected static function toJson($value)
    {
        if ($value instanceof Collection) return $value->toJson();
        if (is_array($value)) return json_encode($value);
        $json = json_encode($value);
        if (json_last_error_msg() == JSON_ERROR_NONE) return $json;
        return collect([$value])->toJson();
    }

    /**
     * 发送一个.. 呃我也不好怎么形容
     *
     * @param String $method
     * @param int $operateID
     * @param string $message
     * @return Collection
     */

    protected function sendMethod($method, $operateID = 0, $message = '')
    {
        $socket = $this->getSocket();
        $data = $this->build($method, $operateID, $message);

        return $this->send($socket, $data);
    }


    /**
     * 获取可以直接啪的 Socket
     *
     * @param String $ip
     * @param Int $port
     * @return resource
     */
    protected function getSocket($ip = null, $port = null)
    {

        $socket = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
        if ($socket === false)
            throw SocketCreationFailed::create();

        $connection = socket_connect($socket, $ip ? $ip : $this->getIp(), $port ? $port : $this->getPort());
        if ($connection === false)
            throw SocketConnectionFailed::create($this->getIp(), $this->getPort(),  socket_strerror(socket_last_error($socket)));

        return $socket;
    }

    /**
     * 构造数据包
     *
     * @param String $method
     * @param int $operateID
     * @param string $message
     * @return Collection
     */
    protected function build($method, $operateID = 0, $message = '')
    {

        $data = collect([
            'Auth' => $this->getVerifyCode(),
        ]);

        $req = collect([
            'Method' => $method,
            'OperateID' => $operateID,
            'Message' => $message
        ]);

        $data->put("Req", $req);

        return $data;
    }


    /**
     * 发射！
     *
     * @param resource $socket
     * @param Collection|array $data
     * @return Collection
     */
    protected function send($socket, $data)
    {
        if ($data instanceof Collection) $data = $data->toJson();
        if (is_array($data)) $data = json_encode($data);

        socket_write($socket, $data, strlen($data));
        $result = '';
        while ($resultBuf = socket_read($socket, 1024)) {
            $result .= $resultBuf;
        }
        socket_close($socket);

        return collect(json_decode($result));
    }
}