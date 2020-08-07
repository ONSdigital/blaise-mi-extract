#!/bin/bash

sandbox="ons-blaise-dev-pds-20:europe-west2"

cloud_sql_proxy -instances=$sandbox:blaise-dev-28475bb5=tcp:3306
