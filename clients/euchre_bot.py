import socket
import json
import time
from multiprocessing import Process
 
HOST = "localhost"
PORT = 8080

def encode_response(message_type, data):
    response = {"type": message_type, "data": {"response": data}}
    response = json.dumps(response) + '\n'
    response = response.encode()
    return response

def play(id):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:

        sock.connect((HOST, PORT))
        sock_file = sock.makefile("r", encoding="utf-8")
        
        sock.send(encode_response("hello", "hello"))
        
        while True:
            try:
                line = sock_file.readline()
            except ConnectionResetError as e:
                # print("Connection reset by peer")
                return
            
            if not line:
                break
            
            message = json.loads(line)

            if id == 0:
                print("#"*100)
                str_message = json.dumps(message, indent=4, ensure_ascii=False)
                # str_message.encode()
                print(str_message)
                print("#"*100)

            match message["type"]:
                case "connectionCheck":
                    # print(message["data"])
                    response = encode_response(message["type"], "Pong")
                    sock.send(response)
                case "pickUpOrPass":
                    # print(message["data"])
                    response = encode_response(message["type"], 1)
                    sock.send(response)
                    # handlePickUpOrPass(message.Data)
                case "orderOrPass":
                    # print(message["data"])
                    response = encode_response(message["type"], 2)
                    sock.send(response)
                    # handleOrderOrPass(message.Data)
                case "dealerDiscard":
                    # print(message["data"])
                    response = encode_response(message["type"], 1)
                    sock.send(response)
                    # handleDealerDiscard(message.Data)
                case "playCard":
                    # print(message["data"])
                    response = encode_response(message["type"], 1)
                    sock.send(response)                    # handlePlayCard(message.Data)
                case "goItAlone":
                    # print(message["data"])
                    response = encode_response(message["type"], 2)
                    sock.send(response)                    # handleGoItAlone(message.Data)
                case "playerID":
                    pass
                    # print(message["data"])
                    # handlePlayerID(message.Data)
                case "dealerUpdate":
                    pass
                    # print(message["data"])
                    # handleDealerUpdate(message.Data)
                case "suitOrdered":
                    pass
                    # print(message["data"])
                    # handleSuitOrdered(message.Data)
                case "plays":
                    pass
                    # print(message["data"])
                    # handlePlays(message.Data)
                case "trickScore":
                    pass
                    # print(message["data"])
                    # handleTrickScore(message.Data)
                case "updateScore":
                    pass
                    # print(message["data"])
                    # handleUpdateScore(message.Data)
                case "error":
                    pass
                    # print(message["data"])
                    # handleError(message.Data)
                case "gameOver":
                    pass
                    # print(message["data"])
                    # res = handleGameOver(message.Data)
                case _:
                    print("Unknown message type: ", message["type"])

            # time.sleep(1)
            # time.sleep(0.05)



if __name__=="__main__":
    players = []
    for i in range(4):
        p = Process(target=play, args=(i,))
        players.append(p)
        p.start()
    
    for player in players:
            player.join()
