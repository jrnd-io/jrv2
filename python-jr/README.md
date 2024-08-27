# Python Plugin package for JR, the data random generator
This package implements a framework to write a producer plugin for [JR](https://github.io/jrnd-io/jrv2) in python.

## Example
A simple plugin can be written as follows:

```python
import logging
import python_jr.jrplugin as jr
import python_jr.producer_pb2 as pb2


# Define a logger with the correct level
logger = jr.Logger(logging_level=logging.DEBUG)
log = logger.log

class MyProducer(jr.JRProducer):
    def Produce(self, request, context):
        key = request.key.decode("utf-8")
        value = request.value.decode("utf-8")
        # do something with the key and value

        # return a response
        response = pb2.ProduceResponse()
        response.bytes = len(request.value)
        response.message = "Some message"

        return response


if __name__ == "__main__":
    jr.serve(MyProducer(), logger)
```
