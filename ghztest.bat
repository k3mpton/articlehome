@echo off 

ghz --insecure --proto ./proto/newspaper.proto --call PublishHome.CreatePublish -D create_data.json localhost:8080

pause 