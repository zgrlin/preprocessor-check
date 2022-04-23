# Go based Web supported Preprocessor-check
This tool helps upload a new version software and new version is automatically UP. Supported Web UI, cookie-based authentication and log output so stop and start method supported CURL, status is support JSON format.

- Mux Framework
- Fully supports SSL/TLS
- Supported JSON
- Advanced upload capabilities
- Cookie based authentication
- Supported OS.exec

Endpoint List:
- / main page
- /login login handler
- /dashboard authenticated user dashboard
- /upload fast and secure upload process
- /logout authenticated user logout point
- /process check process active or not
- /stop process stop 
- /start process start
- /status JSON output process status
- /log/ Static filesystem (log output)

How do first configure?

- set httpServer IP and PORT
- add TLS certificate and key
- create SESSION_SECRET (OS environment)
- add users to users variable
- write to OS.exec script (stop start)
- Fix the URL in GTPL (Golang Template File) 
