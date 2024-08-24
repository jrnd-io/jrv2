import time
import sys
from concurrent import futures
import grpc
import producer_pb2_grpc
import producer_pb2

from grpc_health.v1.health import HealthServicer
from grpc_health.v1 import health_pb2, health_pb2_grpc

class ProducerServicer(producer_pb2_grpc.ProducerServicer):
    def Produce(self, request, context):
        key = request.key.decode("utf-8")
        with open(f"py_{key}", "w") as f:
            f.write(request.value.decode("utf-8"))
            for k, v in request.headers.items():
                f.write(f"{k}: {v}\n")
            f.close()

        return producer_pb2.Empty()


def serve():
    # We need to build a health service to work with go-plugin
    health = HealthServicer()
    health.set("plugin", health_pb2.HealthCheckResponse.ServingStatus.Value('SERVING'))

    # Start the server.
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    producer_pb2_grpc.add_ProducerServicer_to_server(ProducerServicer(), server)
    health_pb2_grpc.add_HealthServicer_to_server(health, server)
    server.add_insecure_port('127.0.0.1:1234')
    server.start()

    # Output information
    print("1|1|tcp|127.0.0.1:1234|grpc")
    sys.stdout.flush()

    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
