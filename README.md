#SLYkey

##A Transparent Peer-to-Peer Public Key Directory

####Files:
- ca: 
    - Implements the webserver for the central authority
    - Includes logic to verify registration transaction POST requests
    - On successful request handling, returns a signature over the hash of the transaction data
- block:
    - The representation of a "block" in the SLYkey blockchain
    - Includes helper functions such as block hash calculation and verifications
- blockqueue:
    - A simple FIFO queue used as a communication buffer for nodeservers
- nodeserver:
    - Implements a "node" in the SLYkey network
    - Nodes accept transactions, calculate proof-of-work, and communication with each other to maintain the blockchain
    - Each node contains to threads: a block worker and a block processor. Block worker calculates proof-of-work to generate new blocks while the block processor handles communication and blockchain synchronization. The two threads communicate with each other with the blockqueue as well as a synchronized channel when new blocks are found.
- rpc:
    - RPC helper methods
- verifier:

Created by **S**erena Wang, **L**ily Tsai, **Y**ihe Huang

CS263, 2015
