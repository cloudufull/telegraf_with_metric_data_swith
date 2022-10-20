# telegraf_with_metric_data_swith
当使用telegraf 采集不同云厂商的数据库实例时，可以使用新增的Swith_data_type参数 强制转换入库冲突的metric的数据类型

 本项目是基于telegraf-1.18.0 修改的 
 目前只有  mysql 和 redis  的input 增加了 Swith_data_type 参数

[root@data-backup telegraf-1.18.0]# ./telegraf  --usage redis

# Read metrics from one or many redis servers
[[inputs.redis]]
  ## specify servers via a url matching:
  ##  [protocol://][:password]@address[:port]
  ##  e.g.
  ##    tcp://localhost:6379
  ##    tcp://:password@192.168.99.100
  ##    unix:///var/run/redis.sock
  ##
  ## If no servers are specified, then localhost is used as the host.
  ## If no port is specified, 6379 is used
  servers = ["tcp://localhost:6379"]

  ## Optional. Specify redis commands to retrieve values
  # [[inputs.redis.commands]]
  # command = ["get", "sample-key"]
  # field = "sample-key-value"
  # type = "string"

  ## specify server password
  # password = "s#cr@t%"

  ## Optional TLS Config
  # tls_ca = "/etc/telegraf/ca.pem"
  # tls_cert = "/etc/telegraf/cert.pem"
  # tls_key = "/etc/telegraf/key.pem"
  ## Use TLS but skip chain & host verification
  # insecure_skip_verify = true

  # when diffent instance metric have diffent data type ,use Swith_data_type force switch
  # eg: Swith_data_type = ["used_memory_peak_perc@string"]
  #     this means swith used_memory_peak_perc to string datatype
  #     当收集不同云厂商metric存入相同influxdb时，由于数据类型不通，
  #     可能有些云厂商的metric会发生写入冲突,无法存入采集的metric，通过强制转换metric数据类型能临时解决此类问题

