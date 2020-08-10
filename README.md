
# Overview


> Note: Plant UML is required to view diagrams in this document. Available in Jetbrain IDEs under the `markdown` setting


This package contains three seperate functions that operate as follows:


<br/>

```plantuml
scale 2.0
left to right direction
skinparam backgroundColor transparent
skinparam shadowing false
scale 2048*1024
left to right direction
skinparam defaultFontSize 10

skinparam rectangleBackgroundColor white
skinparam queueBorderThickness 0.5
skinparam queueBorderColor black
skinparam queueBackgroundColor white
skinparam storageBackgroundColor white
skinparam databaseBackgroundColor white

queue "Pub/Sub" as pubSub
rectangle extractFunction as "Extract Function"
rectangle encryptFunction as "Encrypt Function"
rectangle zipFunction as "Zip Function"
storage encryptBucket as "Encrypt Bucket"
storage encryptedBucket as "Encrypted Bucket"
storage zipBucket as "Zip Bucket"

pubSub --> extractFunction : Event
extractFunction --> encryptBucket
encryptBucket --> encryptFunction
encryptFunction --> encryptedBucket
encryptedBucket --> zipFunction
zipFunction --> zipBucket

```
<br/>

The application architecture has been modelled on the [```Hexagonal Architecture```](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)) pattern. 

# Functions

## Extract Function

```plantuml

skinparam backgroundColor transparent
skinparam shadowing false
scale 2.0
left to right direction
skinparam defaultFontSize 10

skinparam rectangleBackgroundColor white
skinparam queueBorderThickness 0.5
skinparam queueBorderColor black
skinparam queueBackgroundColor white
skinparam storageBackgroundColor white
skinparam databaseBackgroundColor white


queue "Pub/Sub" as pubSub
rectangle extractFunction as "Extract Function"
storage encryptBucket as "Encrypt Bucket"
database db as "Cloud SQL"

pubSub --> extractFunction : Event
extractFunction --> encryptBucket : CSV File
db --> extractFunction : Response Data
extractFunction -[hidden]- db

```
<br/>
This function receives an event from Gcloud pub/pub with a payload describing the 
instrument the user wishes to extract csv response data for. The function  
matches response data to the MI fields in the database and writes a CSV file to the encrypt storage bucket.
 
## Encrypt Function

```plantuml
skinparam shadowing false
scale 2.0
left to right direction
skinparam defaultFontSize 10
skinparam backgroundColor transparent
skinparam rectangleBackgroundColor white
skinparam storageBackgroundColor white

rectangle encryptFunction as "Encrypt Function"
storage encryptBucket as "Encrypt Bucket"
storage encryptedBucket as "Encrypted Bucket"

encryptBucket --> encryptFunction : CSV file
encryptFunction --> encryptedBucket : x.csv.gpg

```

<br/>
The encrypt function is triggered when a file arrives in the `encrypt bucket`. The file is encrypted using the build-in 
Golang PGP encryption functions with the stipulated public key and the result placed in the encrypted bucket.

The Golang libraries allow for the streaming of data into and out of the encryption routines with the result being 
that any sized file can be encrypted without being constrained by memory 
or storage considerations.

## Zip Function

```plantuml
skinparam shadowing false
scale 2.0
left to right direction
skinparam defaultFontSize 10
skinparam backgroundColor transparent
skinparam rectangleBackgroundColor white
skinparam storageBackgroundColor white

rectangle zipFunction as "Zip Function"
storage encryptedBucket as "Encrypted Bucket"
storage zipBucket as "Zip Bucket"

encryptedBucket --> zipFunction : xxx.gpg file
zipFunction --> zipBucket : mi_[instrument]_[data]_[time].zip

```

<br/>

The zip function is triggered when a file arrives in the `encrypted` storage bucket.

The zip is zipped and added to the `zip bucket` using the following file format:

`mi[1]_[2]_[3].zip` where:
1. Is the name of the instrument
2. Is a date in the format DDMMYYYY
3. Is the time is the format HHMMSS

# Configuration

### Google Functions Region Setting

Set the default functions region:

`gcloud config set functions/region europe-west2`

Otherwise, functions will be created somewhere far away in the ether...

### Environment Variables

The following environment variables are available (see the testing section for details on how to create buckets):

* `PUBLIC_KEY=<path to gpg public key file>` - required to encrypt the zip file

* `ENCRYPT_LOCATION=<bucket>` - the GCloud bucket where the file that needs to be encrypted is located. 
Placed there by the `extract_function`.

* `ENCRYPTED_LOCATION=<bucket>` - the GCloud bucket where the file that has been encrypted is located. 
Placed there by the `encrypt_function`.

* `ZIP_LOCATION=<bucket>` - the GCloud bucket where the file that needs to be zipped is located. Placed
there by the `zip_function`.

* `GOOGLE_APPLICATION_CREDENTIALS=<file>` - google credentials file

* `LOG_FORMAT=Terminal|Json` - (json is the default) for logging messages. 
If you want pretty coloured output for local testing use `Terminal`

* `Debug=True|False|NotSet` - set debugging

* `DB_SERVER=<server>` - server address

* `DB_SOCKET_DIR` - the name of the Unix domain socket used by the GCloud SQL instance. Should be set to `/cloudsql` for 
production deployment, unset for testing. 

* `DB_DATABASE=<database>` - the name of the database, defaults to 'blaise'

* `DB_USER=<user>` - the database user

* `DB_PASSWORD=password` - the database password

## Testing

Under the `cmd` directory there is a main.go file which uses the google `FunctionsFramework` to provide some emulation of 
events. Note that these events are triggered by an HTTP request rather than an item arriving on a queue. 
Nevertheless, it provides a means of testing. The corrosponding postman file can be found under the `scripts` directory. 

Note that you will need to run the cloud sql proxy application locally to connect to the mysql instance. A script to
do so, `run_proxy.sh` or `run_proxy.ps1` is in the `scripts` directory. You will need to change the sandbox name to your own.

Add 3 new storage buckets like so before running:

1. gsutil mb gs://\<sandbox\>-encrypt
2. gsutil mb gs://\<sandbox\>-encrypted
3. gsutil mb gs://\<sandbox\>-zip
