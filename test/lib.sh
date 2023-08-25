#!/bin/bash


# check for expected HTTP response
# $1 is filename
# $2 is endpoint name
# $3 is expected content
function check_http_status {

	if grep "$3" "test/$1.header" >/dev/null ; then
	    :
	else
    	echo ""
		echo "TEST FAILED - server returned wrong status code for $2" 
		echo "CHECK test/$1.header"
		echo "EXPECTED: $3"
    	echo ""
	fi
}

