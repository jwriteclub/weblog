Weblog
======

Weblog provides a web based panel which can be used to watch logs from a go
application in real time, including filtering logs.

Log Select Queries
=================

// TOOD query language and examples

Building
========

For ease of building, [`pkger`](https://github.com/markbates/pkger/) static
resources and the compiled [`pigeon`](https://github.com/mna/pigeon/) parser
are included in the source tree. When modifying the static resources or PEG,
please make sure to run `update-static.bat` before check-in.

License
=======

Licensed under the Apache 2.0 license:

       Copyright 2019 Christopher O'Connell
    
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    
        http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

For any questions, please contact jwriteclub@gmail.com

Contributing
============

Contributions are welcome, especially desired are integration with other
logging systems. Please just open a pull request.