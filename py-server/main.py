from ds_algos import *
from time import sleep
import sys

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 main.py <pid>")
        return
    pid = int(sys.argv[1])
    sock = 9000 + pid
    host = 'localhost'
    peerlist = []
    i = pid + 1
    while i != pid:
        peerlist.append(('localhost', 9000 + i))
        i = (i + 1)%5
    peer = Peer('localhost', sock, peerlist) # other joining peers will be added to end of list
    peer.start()
    # st = input("Press any key to initiate leader election: ")
    if pid == 0: peer.initiate_leader_election()
    print(LEADER)
    peer.end()

    if pid == LEADER:
        print(f"Process {pid}: Starting admin server...")

    

if __name__ == '__main__':
    main()