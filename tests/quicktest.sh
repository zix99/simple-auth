#!/bin/bash
set +e

if [[ $* != *--nobuild* ]]; then
  echo Building...
  go build -o simple-auth-server simple-auth/cmd/server
  if [[ $? != 0 ]]; then
    echo Build failed
    exit 2
  fi
fi

rm quicktest.db
./simple-auth-server --verbose --staticfromdisk \
  --api-external=true --api-sharedsecret=super-secret \
  --web-login-twofactor-enabled \
  --email-enabled=false \
  --db-url=quicktest.db &
echo "PID: $!"

sleep 0.5
echo "Waiting for server to come up..."
for i in `seq 1 50`; do
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
