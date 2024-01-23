function main() {
  init_vault_dir=$1
  schema_vault_dir=$2


  echo "INFO: checking rabbitmq instance.."
  test_rabbitmq_env "rabbitmq_local" "rabbitmq_local" 2

  echo "INFO: checking mysql instance.."
  test_mysql_env "mysql_local" "mysql_local" 2

  echo "INFO: all depedency up"
}

function wait_for_http() {
  name=$1
  addr=$2
  sleep_time=$3

  while [[ true ]]; do
    local code=`curl -s -o /dev/null -w "%{http_code}" $addr`
    if [ "$code" -eq "200" ]; then
      echo "INFO: $name is ready"
      break
    fi
    echo "INFO: $name is not ready...sleeping for a while before checking again"
    sleep $sleep_time
  done
}

function test_mysql_env() {
  name=$1
  container=$2
  sleep_time=$3

  while [[ true ]]; do
    local code=$(docker exec $container mysqladmin ping -uroot -p"password"  2>&1)
    if [[ "$code" == *"mysqld is alive"* ]]; then
      echo "INFO: $name is ready"
      break
    fi
    echo "INFO: $name is not ready...sleeping for a while before checking again"
    sleep $sleep_time
  done
}

function test_rabbitmq_env() {
  name=$1
  container=$2
  sleep_time=$3

  while [[ true ]]; do
  #!/bin/bash

    local code=$(docker exec -it $container rabbitmqctl ping 2>&1)
    if [[ "$code" == *"Ping succeeded"* ]]; then
      echo "INFO: $name is ready"
      break
    fi
    echo "INFO: $name is not ready...sleeping for a while before checking again"
    sleep $sleep_time
  done
}




function test_redis_env() {
  name=$1
  container=$2
  sleep_time=$3

  while true; do
    local code=$(docker exec $container redis-cli ping)
    if [[ "$code" == "PONG" ]]; then
      echo "INFO: $name is ready"
      break
    fi
    echo "INFO: $name is not ready... sleeping for a while before checking again"
    sleep $sleep_time
  done
}

main "$@"