# filestore

mysql configuraion
https://coding.imooc.com/lesson/323.html#mid=23347
(
    docker run -d --name mysql-master -p 13306:3306 -v /Users/osir/mysql/conf/master.conf:/etc/mysql/mysql.conf.d/mysqld.cnf -v /Users/osir/mysql/datam:/var/lib/mysql  -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7
    
    docker run -d --name mysql-slave -p 13307:3306 -v /Users/osir/mysql/conf/slave.conf:/etc/mysql/mysql.conf.d/mysqld.cnf -v /Users/osir/mysql/datas:/var/lib/mysql  -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7

    config slave:
    CHANGE MASTER TO MASTER_HOST='',MASTER_PORT=3306,MASTER_USER='slave',MASTER_PASSWORD='slave',MASTER_LOG_FILE='',MASTER_LOG_POS=;
)


docker mysql doc
https://hub.docker.com/_/mysql/

setup master/slave repulication 
https://dev.mysql.com/doc/refman/5.7/en/replication-setup-slaves.html
