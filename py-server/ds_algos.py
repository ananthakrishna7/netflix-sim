import socket
import threading
LEADER=0
ELECTION=False
'''Ports 9000 through 9004 are reserved for leader election'''
class Peer:
    def __init__(self, host, port, peers):
        self.host = host
        self.port = port
        self.pid = port - 9000
        self.peers = peers  # Stores known peer addresses
        self.leader = 0
        self.election = False
        # Create a socket
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        self.sock.bind((self.host, self.port))
        self.sock.settimeout(5)  # Set a timeout for the socket

    def receive_messages(self):
        """ Continuously listens for incoming messages. """
        while True:
            try:
                data, addr = self.sock.recvfrom(1024)
                message = data.decode()
                print(f"Received from {addr}: {message}")
                self.send_message("ACK", addr)
                spl = message.split("::")
                if spl[0] == "ELECTION": self.recv_election_message(spl[1])
                if spl[0] == "COORDINATOR": self.handle_coordinator_message(spl[1])
                # Add new peer if not in list
                if addr not in self.peers:
                    self.peers.append(addr)
            except socket.timeout:
               pass 
            except Exception as e:
                print(f"Error receiving message: {e}")
                break

    def broadcast_message(self, message):
        """ Sends a message to all known peers. """
        for peer in self.peers:
            try:
                self.sock.sendto(message.encode(), peer)
            except Exception as e:
                print(f"Error sending message to {peer}: {e}")

    def send_message(self, message, recipient=None):
        """ Sends a message to a specific recipient. """
        try:
            if recipient:
                self.sock.sendto(message.encode(), recipient)
                return self.wait_for_ack(recipient[1])
            else:
                return False
        except Exception as e:
            print(f"Error sending message: {e}")
            return False
        
    def wait_for_ack(self, port):
        try:
            data, addr = self.sock.recvfrom(1024)
            message = data.decode()
            if message == "ACK" and addr[1] == port:
                return True
            else:
                return False
        except socket.timeout:
            return False
        
    def initiate_leader_election(self):
        self.election = True
        print(f"{self.pid}:Leader election initiated. Sending ELECTION MESSAGES")
        for peer in self.peers:
            try:
                print(f"Attempting to send ELECTION message to {peer}")
                msg = f"ELECTION::{self.pid}"
                result = self.send_message(msg, peer)
                if result:
                    print(f"Message successfully sent to {peer}")
                    break
                else:
                    print("Message sending failed.")
                    continue
            except Exception as e:
                print(f"Failed to send message to {peer}: {e}")
                continue
        else: # worst case scenario
            print("No peers available to send message to.")
            print(f"Process {self.pid}: I am the leader")
            self.leader=self.pid
            # no need to send coordinator messages
            # print("Sending COORDINATOR message")

    def continue_leader_election(self, message, port):
        for peer in self.peers:
            try:
                print(f"Attempting to send ELECTION message to {peer}")
                msg = f"ELECTION::{message},{self.pid}"
                result = self.send_message(msg, peer)
                if result:
                    print(f"Message successfully sent to {peer}")
                    break
                else:
                    print("Message sending failed.")
                    continue
            except Exception as e:
                print(f"Failed to send message to {peer}: {e}")
                continue

    def recv_election_message(self, message):
        """ Receives an election message from another peer. """
        contents = message.split(",")
        contents = [int(i) for i in contents]
        if self.pid in contents:
            print(f"Process {self.pid}: I am the leader")
            self.leader=self.pid
            self.election=False
            print(f"{self.pid}: Sending COORDINATOR message")
            self.broadcast_message(f"COORDINATOR::{self.pid}")
        else:
            print(f"{self.pid}: Passing to next member")
            self.continue_leader_election(message, self.peers[0][1])

    def handle_coordinator_message(self, message):
        """ Receives a coordinator message from the leader. """
        print(f"Received COORDINATOR message from {message}")
        self.leader=int(message)
        self.election=False

    def start(self):
        """ Starts the peer's listening thread. """
        self.thread = threading.Thread(target=self.receive_messages, daemon=True)
        self.thread.start()
        print(f"Peer started at {self.host}:{self.port}")
        # if self.pid == 0: self.initiate_leader_election() # uncomment to quickly show election
        while not self.election: self.heartbeat()
        # leader should have been selected
        # thread.join()

    def heartbeat(self):
        self.send_message("HEARTBEAT", ('localhost', 9000 + self.leader))
        if self.wait_for_ack(self.leader):
            pass
        else:
            self.initiate_leader_election()

    def end(self):
        self.thread.join()

# Example Usage
if __name__ == "__main__":
    host = "127.0.0.1"  # Localhost for testing
    port = int(input("Enter your port: "))  # Choose a unique port for each peer
    peer = Peer(host, port)
    peer.start()
    
    while True:
        msg = input("Enter message: ")
        peer.send_message(msg)
