Related configuration files: 
  backend:
    config.yml: Server port, db, aeskey(16, 24, 32 chars), jwt-key and other information can be configured here.
    It is recommended to replace aes.key(resource password encryption such as machine, database, redis password) and jwT.key (jwt secret key) with a random string.


Service Management:
  Start service:   ./mayfly-go.sh start
  Stop service:    ./mayfly-go.sh stop
  Restart service: ./mayfly-go.sh restart
  Check status:    ./mayfly-go.sh status

Directory Structure:
  bin/mayfly-go  - Main program binary
  mayfly-go.sh   - Service management script
  config.yml     - Configuration file


The project can be accessed directly via host:port (port is server.port as configured in config.yml).
Initial account: admin/admin123.