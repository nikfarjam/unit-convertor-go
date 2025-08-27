#!/bin/bash

if [ "$#" -ne 2 ]; then
  echo "Usage: sh $0 [F|C] [degree]"
  exit 1
fi

from=$1
value=$2

if [ "$from" == "F" ]; then
  from="FAHRENHEIT"
  to="CELSIUS"
elif [ "$from" == "C" ]; then
  from="CELSIUS"
  to="FAHRENHEIT"
else
  echo "Invalid unit. Use 'F' for Fahrenheit or 'C' for Celsius."
  exit 1
fi

if ! [[ "$value" =~ ^-?[0-9]+([.][0-9]+)?$ ]]; then
  echo "Invalid degree value. Please provide a numeric value."
  exit 1
fi

request_body="{ \"value\": $value, \"from\": \"$from\", \"to\": \"$to\" }"

echo "Converting $value degrees $from to $to"
curl --request POST -d "$request_body" localhost:9090/converter
