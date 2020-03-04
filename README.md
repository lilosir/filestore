# filestore

## mysql configuraion
https://coding.imooc.com/lesson/323.html#mid=23347
(   
    master replication: 
    docker run -d --name mysql-master -p 13306:3306 -v /Users/osir/mysql/conf/master.conf:/etc/mysql/mysql.conf.d/mysqld.cnf -v /Users/osir/mysql/datam:/var/lib/mysql  -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7
    
    slave replication
    docker run -d --name mysql-slave -p 13307:3306 -v /Users/osir/mysql/conf/slave.conf:/etc/mysql/mysql.conf.d/mysqld.cnf -v /Users/osir/mysql/datas:/var/lib/mysql  -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7

    config slave:
    docker inspect --format='{{.NetworkSettings.IPAddress}}' mysql-master
    CHANGE MASTER TO MASTER_HOST='',MASTER_PORT=3306,MASTER_USER='slave',MASTER_PASSWORD='slave',MASTER_LOG_FILE='',MASTER_LOG_POS=;
)


## docker mysql doc
https://hub.docker.com/_/mysql/

## setup master/slave repulication 
https://dev.mysql.com/doc/refman/5.7/en/replication-setup-slaves.html


## install redis 
docker run -itd --name redis-test -p 6379:6379 redis

## ceph
https://coding.imooc.com/lesson/323.html#mid=23291
change the load volum path to /Users/myname/ceph
mkdir -p /Users/myname/ceph/www/ceph Users/myname/ceph/var/lib/ceph/osd Users/myname/ceph/www/osd/

## create monitor
docker run -itd --name monnode --network ceph-network --ip 172.20.0.10 -e MON_NAME=monnode -e MON_IP=172.20.0.10 -e CEPH_PUBLIC_NETWORK=172.20.0.10/24 -v /Users/osir/ceph/www/ceph:/etc/ceph -v /Users/osir/ceph/var/lib/ceph/:/var/lib/ceph/ ceph/daemon mon

## create osd 
docker run -itd --name osdnode0 --network ceph-network -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v /Users/osir/ceph/www/ceph:/etc/ceph -v /Users/osir/ceph/www/osd/0:/var/lib/ceph/osd/ceph-0 ceph/daemon osd

docker run -itd --name osdnode1 --network ceph-network -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v /Users/osir/ceph/www/ceph:/etc/ceph -v /Users/osir/ceph/www/osd/1:/var/lib/ceph/osd/ceph-1 ceph/daemon osd

docker run -itd --name osdnode2 --network ceph-network -e MON_NAME=monnode -e MON_IP=172.20.0.10 -v /Users/osir/ceph/www/ceph:/etc/ceph -v /Users/osir/ceph/www/osd/2:/var/lib/ceph/osd/ceph-2 ceph/daemon osd

## create more monitors to make cluster
docker run -itd --name monnode_1 --network ceph-network --ip 172.20.0.11 -e MON_NAME=monnode_1 -e MON_IP=172.20.0.11 -e CEPH_PUBLIC_NETWORK=172.20.0.11/24 -v /Users/osir/ceph/www/ceph:/etc/ceph -v /Users/osir/ceph/var/lib/ceph/:/var/lib/ceph/ ceph/daemon mon
docker run -itd --name monnode_2 --network ceph-network --ip 172.20.0.12 -e MON_NAME=monnode_2 -e MON_IP=172.20.0.12 -e CEPH_PUBLIC_NETWORK=172.20.0.12/24 -v /Users/osir/ceph/www/ceph:/etc/ceph -v /Users/osir/ceph/var/lib/ceph/:/var/lib/ceph/ ceph/daemon mon

## create gateway
docker run -itd --name gwnode --network ceph-network --ip 172.20.0.9 -p 9080:80 -e RGW_NAME=gwnode -v /Users/osir/ceph/www/ceph:/etc/ceph ceph/daemon rgw