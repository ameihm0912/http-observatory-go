http-observatory-go
===================

Small library to assist with making calls to the Mozilla HTTP Observatory
from a go program.

For information on the HTTP observatory, see documentation found
`here <https://github.com/mozilla/http-observatory>`__.

Example
-------

.. code:: go

        results, err := httpobsgo.RunScan("www.mozilla.org", false, false)
        fmt.Println(results.Grade)
