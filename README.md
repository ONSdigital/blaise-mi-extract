
### Google functions Region

Set the default functions region:

`gcloud config set functions/region europe-west2`

Otherwise functions will be created somewhere else in the ether.

### Environment Variables

The following environment variables may be set:

* `MI_BUCKET_NAME=<bucket>` The name of the GCloud bucket

* `GOOGLE_APPLICATION_CREDENTIALS=<file>` Google credentials file for testing

* `LOG_FORMAT=Terminal|Json` (json is the default) for logging messages. 
If you want pretty coloured output for local testing (you do) use `Terminal`.

* `Debug=True|False|NotSet` Set debugging

