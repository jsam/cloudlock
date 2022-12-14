# cloudlock

Regular filelock, but in cloud. 
-------

Cloudlock is a distributed lock based on cloud storage. It uses a bucket to store a lock file. When a lock is acquired, the lock file is created. When the lock is released, the lock file is deleted. Great for distributed CI/CD pipelines or any kind of coordination between multiple processes.

## Supported environments

- [x] Google Cloud 
- [ ] AWS S3       
- [ ] Minio        
- [ ] Github


## Instalation

```bash
go install github.com/jsam/cloudlock
```

## GCS configuration

To start using cloudlock with GCS, there are 2 required settings - project identifier and service account.

They can be provided via environment:
```bash
export GOOGLE_APPLICATION_CREDENTIALS=
export GCP_PROJECT_ID=
```

or via the command line:
```bash
cloudlock [command] [args] --project <project_id> --service-account <svc_account_path>
```



## Usage

### Create a lock
```bash
cloudlock lock [lock-name]
```

### Release a lock
```bash
cloudlock unlock [lock-name]
```
