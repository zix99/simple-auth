#!/bin/bash
set +e

if [[ $* != *--nobuild* ]]; then
  echo Building...
  go build -o simple-auth-server simple-auth/cmd/server
fi

./simple-auth-server --verbose --staticfromdisk \
  --api-external=true --api-sharedsecret=super-secret &
echo "PID: $!"

sleep 0.5
echo "Waiting for server to come up..."
for i in `seq 1 10`; do
  if $(curl -o /dev/null --fail --silent "http://localhost:9002/health"); then
    break
  fi
  printf '.'
  sleep 0.5
done

npm run integration-test
RET=$?

kill %1
wait %1

exit $RET
