
# Overview

This package contains three seperate functions that operate as follows:

ExtractFunction => **ENCRYPT_LOCATION** => EncryptFunction => **ZIP_LOCATION** => ZIPFunction => **ENCRYPTED_LOCATION**

GCP storage triggers have been used to send notifications that a file has arrived in a bucket.

The application architecture has been modelled on the [```Hexagonal Architecture```](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)) pattern. 

# Configuration

### Google Functions Region Setting

Set the default functions region:

`gcloud config set functions/region europe-west2`

Otherwise functions will be created somewhere far away in the ether...

### Environment Variables

The following environment variables are available (see the testing section for details on how to create buckets):

* `PUBLIC_KEY=<path to gpg public key file>` - required to encrypt the zip file

* `ENCRYPT_LOCATION=<bucket>` - the GCloud bucket where the file that needs to be encrypted is located. 
Placed there by the  `extract_function`

* `ENCRYPTED_LOCATION=<bucket>` - the GCloud bucket where the file that has been encrypted is located. 
Placed there by the `encrypt_function`

* `ZIP_LOCATION=<bucket>` - the GCloud bucket where the file that needs to be zipped is located. Placed
there by the `encrypted_function`

* `GOOGLE_APPLICATION_CREDENTIALS=<file>` - google credentials file

* `LOG_FORMAT=Terminal|Json` - (json is the default) for logging messages. 
If you want pretty coloured output for local testing use `Terminal`

* `Debug=True|False|NotSet` - set debugging

* `DB_SERVER=<server>` - server address

* `DB_DATABASE=<database>` - name of the database, defaults to 'blaise'

* `DB_USER=<user>` - database user

* `DB_PASSWORD=password` - database password

## Testing

Under the `cmd` directory there is a main.go file which uses the google `FunctionsFramework` to provide some emulation of 
events. Note however that these events are triggered by an HTTP request rather than an item arriving on a queue. 
Nevertheless, it provides a means of testing. The corrosponding postman file can be found under the `scripts` directory. 

Note that you will need to run the cloud sql proxy application locally to connect to the mysql instance. A script to
do so, `run_proxy.sh` is in the `scripts` directory and you will need to change the sandbox name to your own.


Add 3 new storage buckets like so before running:

1. gsutil mb gs://<sandbox>-encrypt
2. gsutil mb gs://<sandbox>-encrypted
3. gsutil mb gs://<sandbox>-zip
