
REQUETE API Exemple

Add a room :
curl  'http://gopiko.fr:9000/room/add' \ 
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"name": "TestARthur", "creator_id": 1 }'

Get 1 room
curl http://localhost:9000/room/1 \
    --header "Content-Type: application/json" \
    --request "GET"


Get all room
curl http://localhost:9000/room \
    --header "Content-Type: application/json" \
    --request "GET"


Add a song
curl http://localhost:9000/room/song/add \
    --include \    
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"name": "song1", "roomid": 1, "author": "meheu", "trackid": "Icjid" }'
 
Add invite
curl http://localhost:9000/room/invite/add \ 
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"invite": 2, "room_id": 1 }'

curl http://localhost:9000/room/delete \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"id": 5}'


curl --location --request POST 'http://gopiko.fr:9000/room/add' \
--include \
--header "Content-Type: application/json" \
--data '{"name": "Amazing Playlist", "creator_id": 1 }'

curl --location 'http://gopiko.fr:9000/room/2' \
--header "Content-Type: application/json" \
--request "GET"


curl --location 'http://gopiko.fr:9000/room/add' \ 
    
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"name": "TestARthur", "creator_id": 1 }'
curl http://localhost:9000/room \
--header "Content-Type: application/json" \
--request "GET"