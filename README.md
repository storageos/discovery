# Discovery service

## Deployment on Kubernetes

Discovery wants persistent storage for storing cluster related data, create gcloud disk:

    gcloud compute disks create --size 10GB cluster-db