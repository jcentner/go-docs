#!/bin/bash

# NOTE: run with "make test" in top level directory
# Expects clean test database loaded

# source lib
if [[ -f "${0%/*}/lib.sh" ]]; then
    . "${0%/*}/lib.sh"
else
    echo "Error! lib.sh not found!"
    exit 1
fi


# remove old outputs
rm -f *.out
rm -f *.header


# --------------------------------------------------------------------------------
# /api/drafts

name="draftspost"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--post-data '{"name": "Cookie Recipe", "content": "Go to Wholefoods"}' \
	 'http://localhost:8080/api/drafts' 


check_http_status "$name" "/api/$name" "200 OK"


# --------------------------------------------------------------------------------
# /api/findindrafts

name="findindrafts"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--post-data '{"content": "some text for document 1"}' \
	 'http://localhost:8080/api/findindrafts' 

check_http_status "$name" "/api/$name" "200 OK"

if [[ $(jq -r '.[0].name' test/$name.out) = "doc 1" ]] && \
	[[ $(jq -r '.[0].content' test/$name.out) = "some text for document 1" ]] ; then
	:
else
	echo ""
	echo "FAILED - /api/$name returned unexpected json result" 
	echo ""
fi


# --------------------------------------------------------------------------------
# /api/drafts

name="draftsget"

wget -o test/$name.header -O test/$name.out \
	'http://127.0.0.1:8080/api/drafts/'

check_http_status "$name" "/api/$name" "200 OK"

if [[ $(jq -r '.[0].name' test/$name.out) = "doc 0" ]] && \
	[[ $(jq -r '.[1].content' test/$name.out) = "some text for document 1" ]] ; then
	:
else
	echo ""
	echo "FAILED - /api/$name returned unexpected json result" 
	echo ""
fi


# --------------------------------------------------------------------------------
# /api/createcomment

name="createcomment"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--header "Authorization: Bearer thisisatesttoken" \
	--post-data '{"id": "1", "content": "comment on draft", "user": "1"}' \
	 'http://localhost:8080/api/createcomment' 

check_http_status "$name" "/api/$name" "200 OK"


# /api/createcomment
# test as non-admin user should fail

name="createcomment2"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--header "Authorization: Bearer regularusertoken" \
	--post-data '{"id": "3", "content": "comment on draft", "user": "2"}' \
	'http://localhost:8080/api/createcomment'

check_http_status "$name" "/api/$name" "401 Unauthorized"


# /api/createcomment
# test with no token should fail

name="createcomment3"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--post-data '{"id": "3", "content": "comment on draft", "user": "1"}' \
	'http://localhost:8080/api/createcomment'

check_http_status "$name" "/api/$name" "401 Unauthorized"


# /api/createcomment
# test with wrong token should fail

name="createcomment4"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--header "Authorization: Bearer regularusertoken" \
	--post-data '{"id": "3", "content": "comment on draft", "user": "1"}' \
	 'http://localhost:8080/api/createcomment' 

check_http_status "$name" "/api/$name" "401 Unauthorized"

# --------------------------------------------------------------------------------
# /api/commentcomment
# test as admin

name="commentcomment"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--header "Authorization: Bearer thisisatesttoken" \
	--post-data '{"id": "3", "content": "comment on comment", "user": "1"}' \
	'http://localhost:8080/api/commentcomment' 

check_http_status "$name" "/api/$name" "200 OK"


# /api/commentcomment
# test as non-admin user should fail

name="commentcomment2"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--header "Authorization: Bearer regularusertoken" \
	--post-data '{"id": "3", "content": "comment on comment", "user": "2"}' \
	'http://localhost:8080/api/commentcomment'

check_http_status "$name" "/api/$name" "401 Unauthorized"


# /api/commentcomment
# test with no token should fail

name="commentcomment3"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--post-data '{"id": "3", "content": "comment on comment", "user": "1"}' \
	'http://localhost:8080/api/commentcomment'

check_http_status "$name" "/api/$name" "401 Unauthorized"


# /api/commentcomment
# test with wrong token should fail

name="commentcomment4"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--header "Authorization: Bearer regularusertoken" \
	--post-data '{"id": "3", "content": "comment on comment", "user": "1"}' \
	 'http://localhost:8080/api/commentcomment' 

check_http_status "$name" "/api/$name" "401 Unauthorized"


# --------------------------------------------------------------------------------
# /api/commentreaction

name="commentreaction"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--post-data '{"id": "1"}' \
	 'http://localhost:8080/api/commentreaction' 

check_http_status "$name" "/api/$name" "200 OK"


# --------------------------------------------------------------------------------
# /api/comments

name="comments"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--post-data '{"id": "1"}' \
	 'http://localhost:8080/api/comments' 

check_http_status "$name" "/api/$name" "200 OK"


if [[ $(jq -r '.[0].id' test/$name.out) = "1" ]] && \
	[[ $(jq -r '.[1].content' test/$name.out) = "a comment on the comment" ]] && \
	[[ $(jq -r '.[1].user' test/$name.out) = "1" ]] && \
	[[ $(jq -r '.[1].parent.Valid' test/$name.out) = "true" ]] && \
	[[ $(jq -r '.[3].parent.String' test/$name.out) = "3" ]] ; then
	:
else
	echo ""
	echo "FAILED - /api/$name returned unexpected json result" 
	echo ""
fi


# --------------------------------------------------------------------------------
# /api/reaction

name="reaction"

wget -o test/$name.header -O test/$name.out \
	--header "Content-Type: application/json" \
	--post-data '{"id": "1", "content": "😀 ", "user": "1"}' \
	 'http://localhost:8080/api/reaction' 

check_http_status "$name" "/api/$name" "200 OK"

