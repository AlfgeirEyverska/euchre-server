import socket
import time
import json
 
HOST = "localhost"
PORT = 8080

def encode_response(message_type, data):
    response = {"type": message_type, "data": {"response": data}}
    response = json.dumps(response)
    response = response.encode()
    return response

def play():
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:

        sock.connect((HOST, PORT))
        sock_file = sock.makefile("r")
        
        sock.send(b"hello")
        
        while True:
            line = sock_file.readline()
            
            if not line:
                break

            message = json.loads(line)

            match message["type"]:
                case "connectionCheck":
                    print(message["data"])
                    response = encode_response(message["type"], "Pong")
                    sock.send(response)
                case "pickUpOrPass":
                    print(message["data"])
                    response = encode_response(message["type"], 1)
                    sock.send(response)
                    # handlePickUpOrPass(message.Data)
                case "orderOrPass":
                    print(message["data"])
                    response = encode_response(message["type"], 2)
                    sock.send(response)
                    # handleOrderOrPass(message.Data)
                case "dealerDiscard":
                    print(message["data"])
                    response = encode_response(message["type"], 1)
                    sock.send(response)
                    # handleDealerDiscard(message.Data)
                case "playCard":
                    print(message["data"])
                    response = encode_response(message["type"], 1)
                    sock.send(response)                    # handlePlayCard(message.Data)
                case "goItAlone":
                    print(message["data"])
                    response = encode_response(message["type"], 2)
                    sock.send(response)                    # handleGoItAlone(message.Data)
                case "playerID":
                    print(message["data"])
                    # handlePlayerID(message.Data)
                case "dealerUpdate":
                    print(message["data"])
                    # handleDealerUpdate(message.Data)
                case "suitOrdered":
                    print(message["data"])
                    # handleSuitOrdered(message.Data)
                case "plays":
                    print(message["data"])
                    # handlePlays(message.Data)
                case "trickScore":
                    print(message["data"])
                    # handleTrickScore(message.Data)
                case "updateScore":
                    print(message["data"])
                    # handleUpdateScore(message.Data)
                case "error":
                    print(message["data"])
                    # handleError(message.Data)
                case "gameOver":
                    print(message["data"])
                    # res = handleGameOver(message.Data)
                case _:
                    print("Unknown message type: ", message["type"])


# for i in range(10):
#     play()
play()