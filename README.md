# IMGCONV - A microservice for converting HEIF images

## Architecture 
---
## Router 
This package consists of a router implementation which takes in a ConversionManager to handle all the requests from clients and forward them to ConversionManager.
## Manager 
This package consists of a ConversionManager implementation which takes in a StorageRepo for storing the converted images. ConversionManager is responsible for handling the conversion of image from start to end.manager can also take in ConversionListeners for sending conversion notifications. manager puts data into image StorageRepo and also takes out and gives to router.\n 
## Storage
This package consists of StorageRepo implementation. StorageRepo is a simple abstraction for storing images by their IDs, their are two StorageRepo implementations provided, memory and file. 

## Deploying
```
./imgconv 
```