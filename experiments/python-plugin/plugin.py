import time
import sys
import logging
from concurrent import futures
from logging.handlers import QueueHandler, QueueListener
from queue import Queue
from io import StringIO
import grpc
import producer_pb2_grpc
import producer_pb2
import grpc_stdio_pb2_grpc
import grpc_stdio_pb2


from grpc_health.v1.health import HealthServicer
from grpc_health.v1 import health_pb2, health_pb2_grpc


class Logger:
    def __init__(self):
        self.stream = StringIO()  #
        que = Queue(-1)  # no limit on size
        self.queue_handler = QueueHandler(que)
        self.handler = logging.StreamHandler()
        self.listener = QueueListener(que, self.handler)
        self.log = logging.getLogger('python-plugin')
        self.log.setLevel(logging.DEBUG)
        self.logFormatter = logging.Formatter('%(asctime)s %(levelname)s  %(name)s %(pathname)s:%(lineno)d - %('
                                         'message)s')
        self.handler.setFormatter(self.logFormatter)
        for handler in self.log.handlers:
            self.log.removeHandler(handler)
        self.log.addHandler(self.queue_handler)
        self.listener.start()

    def __del__(self):
        self.listener.stop()

    def read(self):
        self.handler.flush()
        ret = self.logFormatter.format(self.listener.queue.get()) + "\n"
        return ret.encode("utf-8")

logger = Logger()
log = logger.log

class ProducerServicer(producer_pb2_grpc.ProducerServicer):
    def Produce(self, request, context):
        log.info(f"Received request: {request}")
        key = request.key.decode("utf-8")
        with open(f"py_{key}", "w") as f:
            f.write(request.value.decode("utf-8"))
            for k, v in request.headers.items():
                f.write(f"{k}: {v}\n")
            f.close()

        return producer_pb2.Empty()


class GRPCController(producer_pb2_grpc.GRPCControllerServicer):
  def Shutdown(self, request, context):
    # Execute what ever needs to be done to clean up and exit the plugin process gracefully
    log.info("Shutting down")

    pass

class GRPCStdioServicer(object):
    """GRPCStdio is a service that is automatically run by the plugin process
    to stream any stdout/err data so that it can be mirrored on the plugin
    host side.
    """
    def __init__(self, log):
        self.log = log

    def StreamStdio(self, request, context):


        """StreamStdio returns a stream that contains all the stdout/stderr.
        This RPC endpoint must only be called ONCE. Once stdio data is consumed
        it is not sent again.

        Callers should connect early to prevent blocking on the plugin process.
        """
        while True:
            sd = grpc_stdio_pb2.StdioData(channel=1, data=self.log.read())
            yield sd

def serve():
    # We need to build a health service to work with go-plugin
    health = HealthServicer()
    health.set("plugin", health_pb2.HealthCheckResponse.ServingStatus.Value('SERVING'))

    # Start the server.
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    producer_pb2_grpc.add_ProducerServicer_to_server(ProducerServicer(), server)
    producer_pb2_grpc.add_GRPCControllerServicer_to_server(GRPCController(), server)
    grpc_stdio_pb2_grpc.add_GRPCStdioServicer_to_server(GRPCStdioServicer(logger), server)
    health_pb2_grpc.add_HealthServicer_to_server(health, server)
    server.add_insecure_port('127.0.0.1:1234')
    server.start()

    # Output information
    print("1|1|tcp|127.0.0.1:1234|grpc",flush=True)
    sys.stdout.flush()

    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
