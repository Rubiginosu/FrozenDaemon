# FrozenGo Daemon 安装手册
## 安装依赖
1.cgroup
debian/other:apt-get -y install cgroup-bin
centos:yum -y install libcgroup
## 安装步骤
1.请先到[Releases](https://github.com/Rubiginosu/frozen-go/releases)下载您所希望的FrozenGo版本<br />
2.请选择一个适合您操作系统的二进制文件包进行下载<br />
3.解压二进制文件，您可以看到FrozenGo的目录<br />
4.第一次运行后才会生成fg.json,请修改conf/fg.json,，将其修改为您所喜欢的配置，<strong>为什么会选择52023作为默认端口呢</strong>，请参见[config.go](https://github.com/Rubiginosu/frozen-go/blob/master/daemon/src/conf/config.go)第46行<br />
5.删除data/servers.json文件中的所有内容，您可以参考此文件写入一些您自己的服务器<br />
6.<strong>找到fg.json中的VerifyCode,Panel应该会用到这个VerifyCode</strong><br />
7.请编辑exec/目录下的配置文件，已经包含了一个示例文件了，您可以根据自己的要求再添加一些新的文件，FrozenGo可以读取这些文件并开服哦~
8.在bin目录下运行./frozen ，<strong>一定要在bin目录下运行哦！</strong>
## 后记
如果您是妹纸，还对画Logo感兴趣，欢迎联系QQ:847072154.<br />
如果您遇到了bug,请提交[Issues](https://github.com/Rubiginosu/frozen-go/issues)<br />
如果您有更好的建议或想法，欢迎Fork 我们的仓库并提出Pull Request!<br />
