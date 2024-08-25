import logging
import jrplugin


logger = jrplugin.Logger(logging_level=logging.DEBUG)
log = logger.log

class MyProducer(jrplugin.JRProducer):
    def Produce(self, request, context):
        log.info("Received request: %s", request)
        key = request.key.decode("utf-8")
        with open(f"py_{key}", "w", encoding="utf-8") as f:
            f.write(request.value.decode("utf-8"))
            for k, v in request.headers.items():
                f.write(f"{k}: {v}\n")
            f.close()

        return


if __name__ == "__main__":
    jrplugin.serve(MyProducer(), logger)
