#!/usr/bin/env roundup

describe "Webreduce/Apps"

req="./test/helper/request"

before() {
  app_name="nyan"
}

it_creates_app() {
  resp=$($req -u "http://localhost:5000/apps/$app_name" -v "PUT")
  code=$(extract_code $resp)

  test $code = "201";
}

it_updates_app() {
  resp=$($req -u "http://localhost:5000/apps/$app_name" -v "PUT")
  code=$(extract_code $resp)

  test $code = 204;
}

extract_code() {
  echo $2
}
