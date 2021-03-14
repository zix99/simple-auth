#!/bin/bash
set +e

# This is a wrapper that builds-and-bounces simple-auth with a mock DB
# In order to run integration tests against.  Alternatively you
# can run it manually with the npm script

cd $(dirname "$0")/..

SATEST_HOST=${SATEST_HOST:-localhost:9002}
export SATEST_HOST

if [[ $* != *--nobuild* ]]; then
  echo Building...
  go build -o simple-auth-server simple-auth/cmd/server
  if [[ $? != 0 ]]; then
    echo Build failed
    exit 2
  fi
fi

ARGS_TEST=
if [[ $* == *--astest* ]]; then
  echo Running as test...
  ARGS_TEST=(-test.run '^TestMain$' -test.coverprofile=integration.cover --)
fi

./simple-auth-server ${ARGS_TEST[@]} --include=tests/testconfig.yml &
echo "PID: $!"

sleep 0.5
echo "Waiting for server to come up..."
for i in `seq 1 50`; do
  if $(curl -o /dev/null --fail --silent "http://${SATEST_HOST}/health"); then
    break
  fi
  printf '.'
  sleep 0.5
done

npm run integration-test
RET=$?

kill -2 %1
wait %1

exit $RET
