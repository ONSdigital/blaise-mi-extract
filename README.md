
# Overview

This package contains three seperate functions that operate as follows:

```
MI_EXTRACT_BUCKET/incoming => ExtractFunction => MI_EXTRACT_BUCKET/zip => ZIPFunction => MI_EXTRACT_BUCKET/encrypt => EncryptFunction => MI_EXTRACT_BUCKET/encrypted
```

GCP storage triggers have been used to send notifications that a file has arrived in a bucket.

The application architecture has been modelled on the ```Hexagonal Architecture``` pattern. 
See [here](https://about.sourcegraph.com/go/gophercon-2018-how-do-you-structure-your-go-apps) and [here](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)) for details.

# Configuration

### Google Functions Region Setting

Set the default functions region:

`gcloud config set functions/region europe-west2`

otherwise functions will be created somewhere far away in the ether...

### Environment Variables

The following environment variables are available:

ENCRYPT_LOCATION=ons-blaise-dev-pds-18-mi-encrypt;ENCRYPTED_LOCATION=ons-blaise-dev-pds-18-mi-encrypted;LOG_FORMAT=Terminal;PUBLIC_KEY=pkg/encryption/keys/paul.gpg

* `ZIP_LOCATION=<bucket>` - the GCloud bucket where the file that needs to be zipped is located. Placed
there by the `extract_function`

* `ENCRYPT_LOCATION=<bucket>` - the GCloud bucket where the file that needs to be encrypted is located. 
Placed there by the  `zip_function`

* `ENCRYPTed_LOCATION=<bucket>` - the GCloud bucket where the file that has been encrypted is located. 
Placed there by the `encrypt_function`

* `GOOGLE_APPLICATION_CREDENTIALS=<file>` - google credentials file

* `LOG_FORMAT=Terminal|Json` - (json is the default) for logging messages. 
If you want pretty coloured output for local testing use `Terminal`

* `Debug=True|False|NotSet` - set debugging

## Testing

