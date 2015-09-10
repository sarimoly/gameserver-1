#!/bin/bash

path=`pwd`
echo ${path}


rm -rf ${path}/Gate/bin/log/*
rm -rf ${path}/Login/bin/log/*
rm -rf ${path}/Center/bin/log/*

cd ${path}/Gate/bin/
./gxgate&

cd ${path}/Login/bin/
./gxlogin&

cd ${path}/Center/bin/
./gxcenter&

cd ${path}