# Battleship Game

This repository contains a server for running a Battleship Game, where players can log in, make any number of attacks,
and log out. The server accepts TCP connections over a specified port (default 8080) and broadcasts actions to every
logged in players.

To run the server, go to the root of the repository and execute:

```
make all
make run
```

That's it. The first command will test and build the code. The second command starts the server.

By default, the server listens at `localhost:8080`. 
To specify the address where the server should listen:
```
./battleship-game serve --listen_address=localhost:8888
```


The board has a 10x10 size with a 25% of space occupied with 1x1 ships.
The ship positioning is random. To change that, you can use a different strategy. 
The `fixed` strategy will always fill the "starting" cells (i.e. {0, 1}, {0, 2}). 
To use this strategy:
```
./battleship-game serve --board_gen_strategy=fixed
```

After setting up the server, it can be tested using `netcat`:
```
nc localhost 8080
login my_name
attack 0 0
attack 1 1
logout
```