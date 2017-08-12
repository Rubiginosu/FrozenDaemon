case $1 in 
	"cg") case $2 in
			"init")
				mkdir /sys/fs/cgroup/cpu/${3} /sys/fs/cgroup/memory/${3} /sys/fs/cgroup/blkio/${3} /sys/fs/cgroup/net_cls/${3}
				# ${4} cpu $5 memmax $6 $7 $8 blkio blkio.throttle.read_bps_device $9 netcls
				tmp=$(cat /sys/fs/cgroup/cpu/${3}/cpu.cfs_period_us)
				cpux=$4
				rmb=$6
				wmb=$7
				mmb=$5
				let rmb=1024*1024*rmb
				let wmb=1024*1024*wmb
				let mmb=1024*1024*mmb
				let tmp=tmp*cpux/100
				echo $tmp > /sys/fs/cgroup/cpu/${3}/cpu.cfs_quota_us
				if test $? then
				 return 1;
				fi
				echo $mmb > /sys/fs/cgroup/memory/${3}/memory.max_usage_in_bytes
				if test $? then
				 return 1;
				fi
				echo "0x0001${9}" > /sys/fs/cgroup/net_cls/${3}/net_cls.classid
				echo "${8} ${rmb}" > /sys/fs/cgroup/blkio/${3}/blkio.throttle.read_bps_device
				echo "${8} ${wmb}" > /sys/fs/cgroup/blkio/${3}/blkio.throttle.write_bps_device
				;;
			"del")
				rmdir /sys/fs/cgroup/cpu/${3} /sys/fs/cgroup/memory/${3} /sys/fs/cgroup/blkio/${3} /sys/fs/cgroup/net_cls/${3}
				;;
			"run")
				/bin/echo ${4} |tee /sys/fs/cgroup/cpu/${3}/tasks /sys/fs/cgroup/memory/${3}/tasks /sys/fs/cgroup/blkio/${3}/tasks /sys/fs/cgroup/net_cls/${3}/tasks
				;;
			esac;;
	"net") DEV=$6;
			case $2 in
				"add")
				tc class add dev $DEV parent 1: classid 1:${3} htb rate ${4}mbit ceil ${5}mbit;
				tc filter add dev $DEV protocol ip parent 1:0 prio 1 handle 1:${3} cgroup;;
				"change")
				tc class change dev $DEV parent 1: classid 1:${3} htb rate ${4}mbit ceil ${5}mbit;;
				"del")
				tc class del dev $DEV parent 1: classid 1:${3};
				tc filter del dev $DEV protocol ip parent 1:0 prio 1 handle 1:${3} cgroup;;
				
			esac;;
	"init") 
	DEV=$2;
	#tc qdisc del dev $DEV root
	tc qdisc add dev $DEV root handle 1: htb;
	tc class add dev $DEV parent 1: classid 1: htb rate 10000mbit ceil 10000mbit;
	service cgconfig restart;
	;;
	esac
	return 0;
