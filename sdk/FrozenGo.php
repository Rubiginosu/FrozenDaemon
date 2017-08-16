<?php
/**
 * Created by PhpStorm.
 * User: axoford12
 * Date: 17-7-29
 * Time: 下午9:38
 */

// 这是用于连接Daemon的SDK..
class Request
{
    public $Method;
    public $OperateID;
    public $Message;
}

class InterfaceRequest
{
    public $Auth;
    public $Req;
}

class ServerAttrElements
{
    public $AttrName;
    public $AttrValue;
}

/*
 * 公共返回对象字段：
 * Status : 包含运行的状态
 * Message : 当状态为-1时包含错误，其余包含返回信息(可能是string也可能是另一个对象)
 * 在下面的函数注释中，我会对Message进行说明.
 */

class FrozenGo
{
    /**
     * 这些常量用于设置服务器配置
     */
    const SERVER_MAX_MEMORY = "MaxMemory"; // 最大内存限制
    const SERVER_EXECUTABLE = "Executable"; // 可执行conf
    const SERVER_MAX_HARD_DISK = "MaxHardDisk"; // 磁盘空间
    const SERVER_NAME = "Name"; // 名称
    const SERVER_EXPIRE = "Expire"; // 到期时间，请传入int类型时间，得到的过期时间=Unix时间戳(传入值 + 开始时间戳)


    /**
     * @var $ip
     * 要连接的ip地址
     * @var $port
     * 端口
     * @var $verifyCode
     * 在Daemon中配置的验证Code.
     */
    private $ip;
    private $port;
    private $verifyCode;

    /**
     * FrozenGo constructor.
     * @param $ip
     * @param $port
     * @param $verifyCode
     */
    public function __construct($ip, $port, $verifyCode)
    {
        $this->ip = $ip;
        $this->port = $port;
        $this->verifyCode = $verifyCode;
    }

    /**
     * @return mixed|string
     * 返回服务器列表，类型是对象数组
     * 对象保存了如下字段：
     * ID: 服务器的ID
     * Name: 服务器名
     * Executable : 可执行文件配置
     * Status: 运行状态
     * UserUid: 运行时期制定用户的Uid
     */
    public function getServerList()
    {
        $servers = $this->SockResult("List");
        $servers->Message = json_decode($servers->Message);
        return $servers;
    }

    /**
     * @param $name
     * 服务器名。
     * @return mixed|string
     * 成功时返回"OK"
     * 失败返回错误信息
     */
    public function createServer($name,$id)
    {
        return $this->SockResult("Create", $id, $name);
    }

    /**
     * @param $id
     * 服务器id (只能是ID)否则daemon无法解析
     * @return mixed|string
     * 返回同上
     */
    public function deleteServer($id)
    {
        return $this->SockResult("Delete", $id);
    }

    /**
     * @param $id
     * 这个密钥对应的服务器id
     * @return mixed|string
     * 对象，包含两个字段，
     * ValidationKeyPair : 对象，包含两个字段：
     *         ID:这个Key对应的ID
     *         Key: 密钥，20位字符
     * GeneratedTIme: 生成的时间。格式大致如下：
     * 2017-07-29T22:35:15.184376223+08:00
     */
    public function getValidationKeyPairs($id)
    {
        $result = $this->SockResult("GetPairs", $id);
        $result->Message = json_decode($result->Message);
        return $result;
    }

    /**
     * @param $url
     * 配置文件url
     * @param $id
     * 要读取的配置文件id
     * @return mixed|string
     */
    public function execInstall($url, $id)
    {
        return $this->SockResult("ExecInstall", $id, $url);

    }

    /**
     * 开启一个未运行的服务器
     * @param $id
     * 要开启的服务器
     * @return mixed|string
     */
    public function startServer($id)
    {
        return $this->SockResult("Start", $id);
    }

    /**
     * 停止一个特定id的服务器。
     * @param $id
     * 停止的服务器id
     * @return mixed|string
     * 返回对象
     */
    public function stopServer($id)
    {
        return $this->SockResult("Stop", $id);
    }

    /**
     * @param $id
     * 服务器id
     * @return object
     */
    public function getServerConfig($id)
    {
        $result = $this->SockResult("GetServerConfig", $id);
        $result->Message = json_decode($result->Message);
        return $result;
    }

    /**
     * @param $id
     * 服务器id
     * @param $elements
     * 元素数组
     * Map类型数组
     * 一个简易的元素集如下：
     * [
     *      [
     *          "AttrName" => FrozenGo::SERVER_NAME,
     *          "AttrValue" => "Axoford12"
     *      ]
     * ]
     * AttrName和AttrValue用于帮助Daemon解析
     * ----------------------------------------------------
     *      支持的标记 详细参见 源码/dmserver/struct
     *
     *
     * @return object
     * 返回成功更新了的Attr的数目或者错误信息,储存在Object对象Message字段中
     */
    public function setServerConfig($id, $elements)
    {
        return $this->SockResult("SetServerConfig", $id, json_encode($elements));
    }

    /**
     * @param $id
     * 要输入的服务器
     * @param $message
     * 输入的命令，如stop
     * say Poi!
     * 都可以,随便大家玩
     * @return object
     * 返回daemon的提示结果
     */
    public function inputToServer($id,$message){
        return $this->SockResult("InputLineToServer",$id,$message."\n");
    }

    /**
     * 获取在线的人
     * @param $id
     * 服务器id
     * @return object
     */
    public function getOnlinePlayers($id){
        $result = $this->SockResult("GetServerPlayers",$id);
        $result->Message = json_decode($result->Message);
        return $result;
    }

    /**
     * 本函数用于向服务器注册key
     * @param $key
     * key，字符串
     * @param $id
     * id 对应服务器id整数
     * @return object
     * 返回结果对象
     */
    public function keyRegister($key,$id){
        return $this->SockResult("KeyRegister",$id,$key);
    }

    /**
     * 本函数用于删除服务器数据目录下的文件，删除 %FGO%/servers/server$id/serverData/$file
     * @param $id
     * 要删除的id
     * @param $file
     * 文件路径及其名称
     * @return object
     */
    public function deleteServerFile($id,$file){
        return $this->SockResult("DeleteServerFile",$id,$file);
    }

    private function SockResult($method, $operateId = 0, $message = "")
    {
        $socket = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
        $conn = socket_connect($socket, $this->ip, $this->port);
        if ($conn < 0) {
            return "5" . socket_strerror($conn);
        }
        $Req = new Request();
        $Req->Method = $method;
        $Req->OperateID = $operateId;
        $Req->Message = $message;
        $InReq = new InterfaceRequest();
        $InReq->Auth = $this->verifyCode;
        $InReq->Req = $Req;
        $sending = json_encode($InReq);
        socket_write($socket, $sending, strlen($sending));
        $result = "";
        while ($resultBuf = socket_read($socket, 1024)) {
            $result .= $resultBuf;
        }
        socket_close($socket);
        return json_decode($result);
    }
}
