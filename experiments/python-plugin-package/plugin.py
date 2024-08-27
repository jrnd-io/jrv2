import logging
import python_jr.jrplugin as jr
import python_jr.producer_pb2 as pb2


logger = jr.Logger(logging_level=logging.DEBUG)
log = logger.log

class MyProducer(jr.JRProducer):
    def Produce(self, request, context):
        log.info("Received request: %s", request)
        key = request.key.decode("utf-8")
        with open(f"py_{key}", "w", encoding="utf-8") as f:
            f.write(request.value.decode("utf-8"))
            for k, v in request.headers.items():
                f.write(f"{k}: {v}\n")
            f.close()

        response = pb2.ProduceResponse()
        response.bytes = len(request.value)
        response.message = "Wrote to file"

        return response


if __name__ == "__main__":
    jr.serve(MyProducer(), logger)
