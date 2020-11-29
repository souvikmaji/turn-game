# Turn Game

Demo application for a turn based websocket game.

## Run Instruction

```
go build

./turn-game -addr :8080
```

Visit: <http://localhost:8080/>

The websocket connection can also be tested using any 3rd part socket connection tester.

Ex: [**WebSocket King Client**](https://chrome.google.com/webstore/detail/websocket-king-client/cbcbkhdmedgianpaifchdaddpnmgnknn?hl=en).

Websocket Address: <ws://127.0.0.1:8080/ws>

Maximum room size is 4. Minimum is 2.
Winning score is 61.
