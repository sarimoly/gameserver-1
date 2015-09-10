#!/bin/bash

# ps -ewf | grep "gxgate" | grep $who | awk '{print $2}'| xargs kill -9
# ps -ewf | grep "gxlogin" | grep $who | awk '{print $2}'| xargs kill -9
# ps -ewf | grep "gxcenter" | grep $who | awk '{print $2}'| xargs kill -9

 ps -ef |grep gxgate |awk '{print $2}'|xargs kill -9
 ps -ef |grep gxlogin |awk '{print $2}'|xargs kill -9
 ps -ef |grep gxcenter |awk '{print $2}'|xargs kill -9