#!/usr/bin/env roundup

# helper
source test/helper/http.sh
source test/helper/response.sh
req="./test/helper/request"

describe "Apps"

before() {
  app_name="nyan"
}

it_creates_app() {
  resp=$($req -u "http://$addr/apps/$app_name" -v "PUT")
  code=$(status_code $resp)

  test $code = "201";
}

it_updates_app() {
  resp=$($req -u "http://$addr/apps/$app_name" -v "PUT")
  code=$(status_code $resp)

  test $code = 204;
}
