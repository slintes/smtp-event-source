# smtpServer

This is a smtpServer which
- accepts mails send to certain addresses
- extracts attachments
- posts attachments to a webhook 

## Context

Currently I use the IFTTT email trigger for doing something similar, but I'd like to implement and host it myself,
for more flexibilty (using my own target mail addresses, using any source mail address, making this a KNative event source)
and just the fun of it :)

## TODOs

- make this a KNative event source

# License

Copyright 2020 Marc Sluiter

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
