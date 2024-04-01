## kubernetes 初始化zookeeper系统
### 用法
```shell
HOSTNAME=online-test-zookeeper-0
go run main.go --config=/usr/local/apache-zookeeper/conf/zoo.cfg \
             --dataDir=/var/lib/zookeeper/${HOSTNAME} \
             --max-node=3 \
             --server=six-zookeeper \
             --stateful_name=online-six-zookeeper \
             --initLimit=20 \
             --tickTime=2000 \
             --syncLimit=10 \
             --clientPort=2181 \
             --maxClientCnxns=1024 \
             --autopurge.snapRetainCount=10 \
             --autopurge.purgeInterval=1 \
             --metricsProvider.className=org.apache.zookeeper.metrics.prometheus.PrometheusMetricsProvider \
             --metricsProvider.httpHost=0.0.0.0 \
             --metricsProvider.httpPort=7000 \
             --metricsProvider.exportJvmInfo=true   
```