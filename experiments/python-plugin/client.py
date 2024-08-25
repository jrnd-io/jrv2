import grpc
import producer_pb2
import producer_pb2_grpc
def run():
    with grpc.insecure_channel('localhost:1234') as channel:
        stub = producer_pb2_grpc.ProducerStub(channel)
        stub.Produce(
            producer_pb2.ProduceRequest(
                key=b"key",
                value=b"value",
                headers={"header": "value"}
            )
        )

        shstub = producer_pb2_grpc.GRPCControllerStub(channel)
        shstub.Shutdown(producer_pb2.Empty())

if __name__ == '__main__':
    run()