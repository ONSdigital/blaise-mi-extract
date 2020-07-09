
### Google functions Region

Set the default functions region:

`gcloud config set functions/region europe-west2`

otherwise functions will be created somewhere far away in the ether.

### Environment Variables

The following environment variables are available:

* `MI_BUCKET_NAME=<bucket>` - the name of the GCloud bucket

* `GOOGLE_APPLICATION_CREDENTIALS=<file>` - google credentials file for testing

* `LOG_FORMAT=Terminal|Json` - (json is the default) for logging messages. 
If you want pretty coloured output for local testing use `Terminal`

* `Debug=True|False|NotSet` - set debugging

