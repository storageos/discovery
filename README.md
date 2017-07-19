# Discovery service

Discovery service is used by StorageOS to help forming clusters.

## API reference

### Create new cluster

Creates new cluster:

```
curl --request POST \
  --url http://discovery.storageos.cloud/clusters
```

Response:
```
{
	"id": "8976384d-08c3-4c3a-b3a9-5e3a6def7062", 
	"size": 3, 
	"createdAt": "2017-07-19T14:08:59.724988221Z",
	"updatedAt": "2017-07-19T14:08:59.724988559Z"
}
```

Here:
* __id__ - cluster ID, should be supplied to StorageOS through env variable CLUSTER_ID
* __size__ - expected cluster size, StorageOS will wait for 3 members to register before starting

### Get cluster status

Get cluster status (expected size, creation date and registered member info):

```
curl --request GET \
  --url http://discovery.storageos.cloud/clusters/8976384d-08c3-4c3a-b3a9-5e3a6def7062
```

Response:

```
{
	"id": "8976384d-08c3-4c3a-b3a9-5e3a6def7062",
	"size": 3,
	"nodes": [
		{
			"id": "node-id",
			"name": "storageos-1",
			"advertiseAddress": "http://1.1.1.1:2380",
			"createdAt": "2017-07-19T14:13:29.182503707Z",
			"updatedAt": "2017-07-19T14:13:29.182503807Z"
		}
	],
	"createdAt": "2017-07-19T14:08:59.724988221Z",
	"updatedAt": "2017-07-19T14:08:59.724988559Z"
}
```

### Register node (internal, used by StorageOS)

StorageOS is using this API for node registration but in some cases it can be useful for debugging:

```
curl --request PUT \
  --url http://discovery.storageos.cloud/clusters/8976384d-08c3-4c3a-b3a9-5e3a6def7062 \
  --header 'content-type: application/json' \
  --data '{\n	"id": "node-id",\n	"name":"storageos-1",\n	"advertiseAddress": "http://1.1.1.1:2380"\n}'
```

Response is same as status API call.

## Building an image

There is a cloudbuild.yaml for [Google Cloud Container Builder](https://cloud.google.com/container-builder/docs/) that can work for you with little modification (project ID). But recommended solution is to use multi-stage Dockerfile:

```
docker build -t <your org name>/discovery:latest -f Dockerfile.multi .
```

## Deployment on Kubernetes

Discovery wants persistent storage for storing cluster related data, create gcloud disk:

    gcloud compute disks create --size 10GB cluster-db

Then, use supplied `hack/deployment.yml` file to create a deployment:

    kubectl create -f hack/deployment.yml

To make it reachable, create a service:

    kubectl create -f hack/svc.yml    


