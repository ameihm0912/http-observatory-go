http-observatory-go
===================

Small library to assist with making calls to the Mozilla HTTP Observatory
from a go program.

For information on the HTTP observatory, see documentation found
`here <https://github.com/mozilla/http-observatory>`__.

Documentation
-------------
Documentation can be found
`here <https://godoc.org/github.com/ameihm0912/http-observatory-go>`__.

Example
-------

.. code:: go

        results, err := httpobsgo.RunScan("www.mozilla.org", false, false)
        fmt.Println(results.Grade)
