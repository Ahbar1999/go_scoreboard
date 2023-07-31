# Running using docker image 

```docker pull ahbar99/scoreboard```  

replace 'port' with the port number you wish to run the server on   
```docker run --publish <port>:8080 ahbar99/scoreboard```

# API Details
base url: localhost:8080/api/  
	example with port = 8080: list all players: [GET: localhost:8080/api/players/]  
							: get player with rank = 2: [GET: localhost:8080/api/players/rank/:2]  
							: delete player with rank = 2: [DELETE: localhost:8080/api/players/rank/:2]  
							: change name of the player with id = 2: [PUT: localhost:8080/api/players/:2]  

# A word of caution
this is an example project, various basic features have been omitted like strict pattern matching, error handling etc.
For example were you to send a POST request to any url starting with /api/players/, the api will accept the request

