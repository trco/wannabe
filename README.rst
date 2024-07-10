=======
Wannabe
=======

A versatile Go tool for effortlessly generating mock APIs.

.. contents:: **Table of Contents**
   :depth: 2
   :local:
   :backlinks: none

How does it work?
=================

Proxy mode
**********

In "proxy" mode, Wannabe operates as a proxy server. It derives a cURL command \
from the received request based on your request matching configuration and \
hashes it to create a unique identifier. Wannabe then proxies the received \
request to the host defined in the request and, upon receiving the response, \
stores a record in the configured storage provider using the previously \
generated hash as the key. Each record includes the original request and its \
corresponding response from the upstream server.

Server mode
***********

In "server" mode, Wannabe functions as a standalone server. Upon receiving a \
request, it generates a cURL command from it based on your request matching \
configuration and generates a hash from the prepared cURL command. Wannabe \
then looks up the matching record in the storage provider using the hash as a \
record key and responds with the stored response if it finds a match, or with \
an error if a matching record is not found.

Mixed mode
**********

In "mixed" mode, Wannabe functions as both a standalone server and a proxy server. \
Upon receiving a request, it generates a cURL command from it based on your request \
matching configuration and generates a hash from the prepared cURL command. If it finds \
a matching record for the received request using the hash as a record key, Wannabe \
responds with the recorded response. If no matching records are found in the storage, \
Wannabe proxies the received request to the host defined in the request and, upon receiving \
the response, stores a record in the configured storage provider using the previously generated \
hash as the key.

Usage examples
==============

Wannabe seamlessly mimics any desired API, whether external or internal, existing \
or still in development, and without business logic. It can effectively become the \
API you need for faster and better development and testing processes.

Mocking external APIs
*********************

Wannabe allows developers to record and simulate the behavior of external services, \
eliminating the need for reliance on those services during development and testing. \
This spans from initial development to regression testing.

Mocking internal APIs
*********************

Developers can use Wannabe to prepare mocks of non-existing APIs and share them with \
other teams before implementing any business logic. These mocks facilitate development \
and testing processes, spanning from initial development to regression testing.

Reusability
***********

Wannabe records, along with their underlying configuration files, can be shared among \
developers, teams, and businesses. This accelerates development processes by providing \
robust and well-tested mocks.

Wannabe supports numerous other use cases. If you discover an innovative use case for Wannabe, \
please share it with us

Getting started
===============

Running in Docker
*****************

Wannabe provides an official Docker image for running the application within a container.

To ensure the application starts correctly, it is mandatory to mount the following files \
from the filesystem where the Wannabe container is running to the root folder of the container:

1. “config.json” configuration file
2. “wannabe.crt”
3. “wannabe.key”

The wannabe.crt and wannabe.key files represent the self-signed certificate that you need \
to create and mount to the Wannabe container. This is necessary for Wannabe to proxy requests \
to other servers using the HTTPS protocol. This certificate should also be added to the service \
that will be using Wannabe. For details on the configuration file usage, refer to the Configuration \
section.

Creation of self-signed certificate
***********************************

.. code-block:: shell
    $ openssl genrsa -out wannabe.key 2048
    $ openssl req -new -x509 -key wannabe.key -out wannabe.crt -days 3650

Example of adding the certificate to the containerized service
***************************************************************

.. code-block:: shell
    $ docker cp ./wannabe.crt containerized-service:/usr/local/share/ca-certificates/
    $ update-ca-certificates

Running Wannabe container
*************************

The Wannabe server runs on port 6789 within the container, while the API runs on port 6790.

.. code-block:: shell
    $ docker run -d \
    -p 6789:6789 \
    -p 6790:6790 \
    -v $(pwd)/config.json:/config.json \
    -v $(pwd)/wannabe.crt:/wannabe.crt \
    -v $(pwd)/wannabe.key:/wannabe.key \
    --name wannabe \
    wannabe // add official image

Running as a standalone server
******************************

Like any Go program, Wannabe can be launched by simply cloning the repository, adding a custom \
config.json, wannabe.crt and wannabe.key to the root of the cloned repository, compiling the source \
code into an executable binary file using the “go build” command, and then running the program with \
the “go run” command.

Example
=======

TODO

Configuration
=============

Wannabe requires a config.json configuration file. Any changes made to the \
configuration file will only take effect after restarting the standalone Wannabe \
server or the one running in the container.

You can configure the following root fields: mode, storageProvider and wannabes. \
See following subsections for the details of all the options that can be configured using those root fields.

Mode
****
.. code-block:: javascript
    {
        "mode": string // “proxy”, “server”, “mixed”; defaults to “mixed”
    }

The mode field defines how a Wannabe container operates. Refer to the "How does it work?" \
section for details.

Storage provider
****************

.. code-block:: javascript
    {
        "storageProvider": {
            "type": string, // “filesystem”; defaults to “filesystem”
            "filesystemConfig": filesystemConfig // see below
        }
    }

The storageProvider field configures the storage for saving the records. Based \
on the specified "type," relevant configuration should be defined. \
For "type": "filesystem," the "filesystemConfig" is required.

Type
----

The type field defines the type of storage provider Wannabe should use.

FilesystemConfig
----------------

.. code-block:: javascript
    {
        "filesystemConfig": {
            "folder": string, // path to the folder, relative to the configuration file
            "regenerateFolder": string, //  path to the folder, relative to the configuration file
            "format": string // “json”
        }
    }

The filesystemConfig field defines the configuration of the file system storage provider.

Folder
------

The folder field defines the folder for storing the records.

RegenerateFolder
----------------
The regenerateFolder field defines a folder for storing the regenerated records.

Format
------

The format field defines the format in which the records are stored.

Wannabes
********

.. code-block:: javascript
    {
        "wannabes": {
            "example.com": {
                "requestMatching": {...}, // see Request matching section
                “records”: {...} // see Records section
            },
            "api.github.com": {
                "requestMatching": {...},
                “records”: {...}
            },
            ...
    }

Wannabes are a map of configurations for request matching and \
records for the hosts that Wannabe mocks, where the host name \
should be used as a key in the map.

Request matching
****************

The "Request matching" field configures the generation of cURL \
commands and the underlying unique hash identifier for each request \
received by Wannabe. It allows you to include or exclude specific parts \
of the requests, whether static or dynamic, from the generation of cURL \
commands corresponding to those requests, or replace specific request \
parts with placeholders. This approach enables the generation of identical \
cURL commands and underlying hashes for multiple unique requests, thereby \
enabling Wannabe to store a single record with one response for all those \
multiple unique requests in "proxy" mode, and to respond with an identical \
response for all those requests when in "server" or "mixed" mode.

For example, you can record responses for all possible requests to the \
Google Analytics Data API for a single "propertyId," but since you excluded \
the dynamic "propertyId" from request matching by replacing it with a static \
placeholder, different "propertyIds" in the request will result in identical \
cURL commands and underlying hashes, and Wannabe will respond with the responses \
recorded for a single "propertyId."

For a better understanding of how this works, refer to the "Usage of index wildcards", \
“Usage of key wildcards” and “Usage of regexes” sections and the explanations provided therein.

Important note: When configuring request matching to include a specific header \
in the generation of the cURL command and the underlying unique hash identifier \
for requests, you cannot exclude the same header from being stored in the request \
field of the records. This ensures that you can always regenerate existing records \
with a new request matching configuration, including this specific header. If headers \
to be included in request matching are not set, all of them are included in matching, \
and none of them can be excluded from being stored in the request field of records.

.. code-block:: javascript
    {
        "requestMatching": {
            "host": {
                "wildcards": [
                    {
                        "index": number, // required
                        "placeholder": string // optional; defaults to "{wannabe}"
                    }
                ],
                "regexes": [
                    {
                        "pattern": string, // required
                        "placeholder": string // optional; defaults to "{wannabe}"
                    }
                ]
            },
            "path": {
                "wildcards": [
                    {
                        "index": number, // required
                        "placeholder": string // optional; defaults to "{wannabe}"
                    }
                ],
                "regexes": [
                    {
                        "pattern": string, // required
                        "placeholder": string // optional; defaults to "{wannabe}"
                    }
                ]
            },
            "query": {
                "wildcards": [
                    {
                        "key": string, // required,
                        "placeholder": string // optional; defaults to "{wannabe}"
                    }
                ],
                "regexes": [
                    {
                        "pattern": string, // required
                        "placeholder": string // optional; defaults to "{wannabe}"
                    }
                ]
            },
            "body": {
                "regexes": [
                    {
                        "pattern": string, // required
                        "placeholder": string // optional; defaults to "{wannabe}"
                    }
                ]
            },
            "headers": {
                "include": array of strings, // if not set all headers are included
                "wildcards": [
                    {
                        "key": string, // required
                        "placeholder": string // optional; defaults to "{wannabe}"
                    }
                ]
            }
        }
    }

Usage of index wildcards
------------------------

.. code-block:: javascript
    {
        "host": {
            "wildcards": [
                {
                    "index": 0,
                    "placeholder": “placeholder”
                }
            ]
        }
    }

When generating cURL commands to be hashed as unique identifiers of \
requests, the host "https://analyticsdata.googleapis.com" and the given \
wildcard will result in the "https://{placeholder}.googleapis.com" host \
being included in the cURL command. After trimming the protocol prefix, \
the host is split using "." as a separator, and the value at the defined \
index is replaced with a defined placeholder, or “{wannabe}” placeholder by default.

Requests that differ only in the value at the first index of the host will \
result in the same cURL command and hash. Therefore, they will be stored as \
a single record with the underlying response in storage.

Usage of key wildcards
----------------------

.. code-block:: javascript
    {
        "query": {
            "wildcards": [
                {
                    "key": "userId",
                    "placeholder": "{placeholder}"
                }
            ]
        }
    }

When generating cURL commands to be hashed as unique identifiers of \
requests, the query "?status=completed&userId=123456" and the given \
wildcard will result in the "?status=completed&userId={placeholder}" \
query being included in the cURL command. After splitting the query \
string into an object of key-value pairs, the value of the key defined \
in the wildcard is replaced with a defined placeholder, or “{wannabe}” \
placeholder by default.

Requests that differ only in the value of the defined key in the query \
will result in the same cURL command and hash. Therefore, they will be \
stored as a single record with the underlying response in storage.

Usage of regexes
----------------

.. code-block:: javascript
    {
        "path": {
            "regexes": [
                {
                    "pattern": "(\\d+):runReport",
                    "placeholder": "{propertyId}:runReport"
                }
            ]
        }
    }

When generating cURL commands to be hashed as unique identifiers \
of requests, the path "/v1beta/properties/123456789:runReport" and \
the given regex will result in the "/v1beta/properties/{placeholder}:runReport" \
path being included in the cURL command. The regex pattern is replaced \
with the defined placeholder, or the "{wannabe}" placeholder by default.

Requests that differ only in the regex-defined pattern of the path \
will result in the same cURL command and hash. Therefore, they will \
be stored as a single record with the underlying response in storage.

Records
-------

.. code-block:: javascript
    {
        "records": {
            "headers": {
                "exclude": []string
            }
        }
    }

The "Records" field allows configuring headers to be excluded \
from the request field of the stored records. This allows exclusion \
of headers that might pose security risks, such as “Authorization” \
headers containing access tokens, API keys, or other credentials.

Important note: When configuring request matching to include a specific \
header in the generation of the cURL command and the underlying unique \
hash identifier for requests, you cannot exclude the same header from being \
stored in the request field of the records. This ensures that you can always \
regenerate existing records with a new request matching configuration, including \
this specific header. If headers to be included in request matching are not set, \
all of them are included in matching, and none of them can be excluded from being \
stored in the request field of records.

Defaults
--------

When the “mode” or “storageProvider” fields are not defined in the \
configuration, they default to the values below.

.. code-block:: javascript
    {
        "mode": "mixed",
        "storageProvider": {
            "type": "filesystem",
            "filesystemConfig": {
                "folder": "records",
                "regenerateFolder": "records/regenerated",
                "format": "json"
            }
        }
    }

Record entity
=============

After Wannabe retrieves a response for a specific request, it stores \
it in a record within the storage provider. The hash generated from the \
request's cURL command is used as the key for the stored record, and the \
record is added to the folder named after the host the request was made to.

For example, if the storage provider is the file system, and the default \
"records" folder is set for storing records, and the hash generated from \
the request's cURL command is "d050d9e39f…190b4037a", and the request was \
made to "api.github.com", the record would be stored at the path "records/api.github.com/d050d9e39f…190b4037a.json".

A record contains (1) a request object, which consists of hash, curl, \
httpMethod, host, path, query, headers, and body fields, (2) a response \
object, which consists of statusCode, headers, and body fields, and (3) a \
metadata object, which consists of generatedAt and regeneratedAt fields.

Regenerate records
==================

Wannabe supports the regeneration of existing records with new request-matching \
configurations. To prepare for the regeneration of existing records, follow these steps:

1. Prepare a new configuration file with updated “requestMatching” configurations \
for “wannabes” you would like to regenerate records for and set custom "regenerateFolder" \
when the file system is configured as the storage provider.
2. Restart the running Wannabe instance to load the new configuration file.
3. Execute the regeneration by calling the GET /wannabe/api/regenerate endpoint. \
Refer to the "API Reference" for details.
4. To use the newly regenerated records, copy them to a configured "folder", but ensure \
they are not mixed with previous records tied to different configuration files.

Important Notes:
- Use the regenerate records functionality with caution and always follow the described steps.
- Know which records correspond to which configuration file and ensure that configuration \
files are always used with relevant records. Regenerated records should not be used with \
unrelated configuration files.
- Mixing regenerated records with records used for regeneration in a configured "folder" \
can result in an inability to separate records.
- The "regenerateFolder" path should not be the same as the "folder" path. If it is, the \
folder will contain a mix of regenerated records and initial records used for regeneration, \
which could be impossible to separate, especially in cases with a large number of records.

API reference
=============

**GET /wannabe/api/records/<hash>?host=<host>**

Description: Retrieves all the records, all the records for a specified host, or a \
single record for a specified host.

Parameters:

<host> (string, optional): Host for which the records are stored. If the <hash> \
parameter is provided, <host> is required.

<hash> (string, optional): The unique identifier of the record. If the <hash> \
parameter is provided, <host> is required.

Response body

.. code-block:: javascript
  [
      {
          "request": {
              "hash": string,
              "curl": string,
              "httpMethod": string,
              "host": string,
              "path": string
              "query": {
                  “key”: string
                  …
              },
              "headers": {
                  “key”: string[],
                  …
              },
              "body": object                
          },
          "response": {
              "statusCode": integer,
              "headers": {
                  “key”: string[],
                  …
              },
              "body": object / string
          },
          "metadata": {
              "generatedAt": {
                  "unix": integer,
                  "utc": string
              },
              "regeneratedAt": {
                  "unix": integer,
                  "utc": string
              },
          }
      }
  ]

**POST /wannabe/api/records**

Description: Stores received records in the configured storage provider.

The “recordProcessingDetails” array in the response body contains the \
hash and message for each record posted in the request body in the same \
indexed order. This means that the record processing details for the first \
record posted in the request body are at index zero in the “recordProcessingDetails” \
array. In the case of a successfully stored record, the message equals "success", \
while in the case of inability to store the record, the message describes the error \
for why storing failed.

Request body

.. code-block:: javascript
  [
      {
          "request": {
              "httpMethod": string,
              "host": string,
              "path": string
              "query": {
                  “key”: string
                  …
              },
              "headers": {
                  “key”: string[],
                  …
              },
              "body": object                
          },
          "response": {
              "statusCode": integer,
              "headers": {
                  “key”: string[],
                  …
              },
              "body": object
          },
      }
  ]

Response body

.. code-block:: javascript
  {
      "insertedRecordsCount": integer,
      "notInsertedRecordsCount": integer,
      "recordProcessingDetails": [
          {
              "hash": string
              "message": string
          }
      ]
  }

**DELETE /wannabe/api/records/<hash>**

Description: Deletes all the records for a specified host or a \
single record for a specified host.

Parameters:

<host> (string, required): Host for which the records are stored.

<hash> (string, optional): The unique identifier of the record. \
If the <hash> parameter is provided, <host> is required.

Response body

.. code-block:: javascript
  {
      "message": string,
      "hashes": string[]
  }

**GET /wannabe/api/regenerate?host=<host>**

Description: Regenerates records for a specific host using the \
provided "wannabe" configuration. See the "Regenerate Records" \
section for details.

<host> (string, required): Host for which the records should be regenerated.

Response body

.. code-block:: javascript
  {
      "message": string,
      "regeneratedHashes": string[],
      "failedHashes": string[]
  }

Contributing
============

Thank you for considering contributing to Wannabe! Contributions from \
the community are more than welcome to help improve the project and make \
it even better.

How to Contribute
*****************

To contribute to Wannabe, follow these steps:

1. Fork the repository.
2. Create a branch.
3. Develop.
4. Commit changes.
5. Submit a pull request.

Your pull request will be reviewed, and you may be asked to make further \
changes or address feedback before your contribution is accepted. Adding \
and updating existing tests is mandatory for pull requests to enter the \
review process.

Where to start
**************

If you're eager to contribute to Wannabe but aren't sure where to begin, \
we've got you covered! You can dive right in by exploring our open issues \
or checking out our existing "next step" ideas. Simply head over to the \
Issues tab to get started!

Author
======

Uroš Trstenjak, https://www.linkedin.com/in/uros-trstenjak/.