# cloudlock

Regular filelock, but in cloud. 


## Supported environments

- Google Cloud


## Usage

```bash
cl lock ci-lock-deployment --wait --params timeout=30 
<locked command>
cl unlock 
```