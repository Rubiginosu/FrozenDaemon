#include <stdio.h>
#include<libcgroup.h>
#include <stdlib.h>
#include <sys/stat.h>
#include<unistd.h>
#include <memory.h>
#include <errno.h>

bool err_handler_code(int code, char *by);
enum ERROR {
    ERR_CG_INIT_FAILED = 2,
    ERR_ARGUMENTS_NOT_FILLED,
    ERR_NOT_ROOT_RUNNER,
    ERR_CPU_SET_FAILED,
    ERR_MEMORY_SET_FAILED,
    ERR_INVALID_ARGUMENTS
};

enum METHOD{
    METHOD_CREATE = 0,
    METHOD_DELETE,
    METHOD_ATTACH,
    METHOD_MODIFY
};
/**
 *
 * @param argc
 * @param argv
 * argv[1] 操作Method ID;Create,Delete,Attach
 * argv[2] Controller控制族群名称.
 * argv[3] CPU 使用率
 * argv[4] Memory大小
 * @return
 */
int main(int argc, char **argv) {

    char *control_group;
    struct cgroup *cg;
    if(argc <= 2) return ERR_ARGUMENTS_NOT_FILLED;

    // 检查身份

    if (getuid() != 0) {
        return ERR_NOT_ROOT_RUNNER;
    }

    int method = atoi(argv[1]);
    switch(method){
        case METHOD_CREATE:
            control_group = argv[2];
            long cpu = 20, memory = 10240;
            cpu = atoi(argv[3]);
            // 解析参数
            memory = atoi(argv[4]);
            if(!err_handler_code(cgroup_init(), "initial")) return ERR_CG_INIT_FAILED;
            // 检查挂载
            char* cpu_path ;

            char* cpu_mount_point;
            cgroup_get_subsys_mount_point("cpu",&cpu_mount_point);
            //char cpu_path[200];
            cpu_path = malloc(sizeof(char) * (strlen(cpu_mount_point) + strlen("/") + strlen(control_group)));
            sprintf(cpu_path,"%s/%s",cpu_mount_point,control_group);
            struct stat s;
            stat(cpu_path,&s);
            if(errno == ENOENT){
                if(mkdir(cpu_path, 0555)){
                    printf("Error occurred during making sub_controller,%d\n",errno);
                    printf("Reason:%s",strerror(errno));
                }
            }
            free(cpu_path);
            char* mem_path;
            char* mem_mount_point;
            cgroup_get_subsys_mount_point("memory",&mem_mount_point);
            mem_path = malloc(sizeof(char) * (strlen(mem_mount_point) + strlen("/") + strlen(control_group)));
            sprintf(mem_path,"%s/%s",mem_mount_point,control_group);
            stat(mem_path,&s);
            if(errno == ENOENT){
                if(!mkdir(mem_path, 0555)){
                    printf("Error occurred during making sub_controller,%d\n",errno);
                    printf("Reason:%s",strerror(errno));
                }
            }
            free(mem_path);

            cg = cgroup_new_cgroup(control_group);
            cgroup_add_controller(cg, "cpu");
            cgroup_add_controller(cg, "memory");

            if(!err_handler_code(cgroup_set_value_int64(cgroup_get_controller(cg,"cpu"),"cpu.cfs_quota_us",(int64_t)cpu*100),"Set CPU")){
                cgroup_free(&cg);
                exit(ERR_CPU_SET_FAILED);
            }
            if(!err_handler_code(cgroup_set_value_int64(cgroup_get_controller(cg,"memory"),"memory.max_usage_in_bytes",(int64_t)memory),"Set Memory")){
                cgroup_free(&cg);
                exit(ERR_MEMORY_SET_FAILED);
            }
            err_handler_code(cgroup_modify_cgroup(cg),"Save ..");
            cgroup_free(&cg);
            return 0;
        default:
            return ERR_INVALID_ARGUMENTS;
    }
}

/**
 *
 * @param code
 * 错误码
 * @param by
 * 打印调用者，方便查看.
 */


bool err_handler_code(int code, char *by) {
    if (code != 0) {
        printf("Caused by:%s ,Code：%d   Message:%s \n", by, code, cgroup_strerror(code));
        return false;
    }
    return true;
}

