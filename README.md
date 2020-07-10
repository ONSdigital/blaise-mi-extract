
# Overview

This package contains three seperate functions that operate as follows:

```
MI_EXTRACT_BUCKET/incoming => ExtractFunction => MI_EXTRACT_BUCKET/zip => ZIPFunction => MI_EXTRACT_BUCKET/encrypt => EncryptFunction => MI_EXTRACT_BUCKET/encrypted
```

GCP events are used to receive a notification that a file has arrived in a bucket.

# Configuration

### Google Functions Region

Set the default functions region:

`gcloud config set functions/region europe-west2`

otherwise functions will be created somewhere far away in the ether.

### Environment Variables

The following environment variables are available:

* `FILE_PROVIDER=Google|other` - name name of the file storage provider, defaults to Google

* `MI_BUCKET_NAME=<bucket>` - the name of the GCloud bucket if using the Google file provider

* `GOOGLE_APPLICATION_CREDENTIALS=<file>` - google credentials file for testing

* `LOG_FORMAT=Terminal|Json` - (json is the default) for logging messages. 
If you want pretty coloured output for local testing use `Terminal`

* `Debug=True|False|NotSet` - set debugging

