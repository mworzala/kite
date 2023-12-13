# DEV NOTES

- [ ] Switch Server
    [ ] Return to configuration to swap worlds
    [ ] Switch world, no reconfiguration
    [ ] Stay in the same world.

## Return to configuration to swap worlds


1. Ensure the client is the right state (play or configuration) before commencing the switch.
1. Proxy Connects to the new server.
2. Proxy logs the player in on the new server, passing the previously used creds/params to the new server (but the client doesn't know anything yet)
3. As soon as the login is complete on the new server, send the client to the configuration phase from the proxy. Client is in the same state in both servers.
4. As soon as the client is in the configuration phase, disconnect the proxy from the old server.
5. Start sending the data to the client from the new server.

We should put the client to the configuration state as soon as we start login. In the future it's better to do it after encryption success, but we're just going to ignore it for now.

As soon as the login success packet is received, we should start passing new server data through to the client.

### TODO

- [ ] Know the client is the right state (play or configuration) before commencing the switch.
- [ ] Ensure the state is right
- [ ] Connect -> login to the new server
- [ ] Wait for the login to complete, put the client in the configuration state
- [ ] Start passing new server data to the client.
- [ ] Disconnect the client from the old server.