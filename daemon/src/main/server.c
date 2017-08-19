//
// Created by axoford12 on 8/19/17.
// C实现的Wrapper 分离的操作更方便～
//
#define _GNU_SOURCE
#include <sched.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <errno.h>
#include <grp.h>
int main(int argc, char** argv){
    if(argc < 5){
	// 若参数不够，可能在获取Argv[5]的时候出现指针越界....然后蜜汁错误
	// 原版本没这一句，但删掉后果自负。
        return -1;
    }
    int uid=atoi(argv[1]); // 转为整型
    if(errno){
        printf("%d\n",errno);
    }

    char* home=argv[2];
    char* sv=argv[3];
    unshare(CLONE_NEWUTS|CLONE_NEWIPC|CLONE_FS|CLONE_FILES);
    if(errno){
        printf("Unsharing error:%d\n",errno);
    }
    chroot(home); // chroot到指定目录
    chdir(sv);    // chdir 到指定目录 比如/serverData
    // setgroups不用会出现蜜汁root组
    if (!(setgroups(1,&uid))){
        printf("set groups error.");
    }
    // 在降权之前setgroups.
    while(setgid(uid)!=0) sleep(1);
    while(setuid(uid)!=0) sleep(1);
    execvp(argv[4],&argv[4]); // 直接execvp即可.
    if(errno){
        printf("Error occurred at execvp:%d\n",errno);
    }
}
