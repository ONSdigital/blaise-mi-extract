#!/bin/bash

sandbox="ons-blaise-dev-pds-27:europe-west2"

cloud_sql_proxy -instances=$sandbox:blaise-dev-068d804a=tcp:3306
