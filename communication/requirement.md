
# What am I going to build?

A communication infrastructure serving for distributed systems.

# What feature/functionality need to support?

1. TCP and UDP transmission
2. Unicast and Multicast
3. FIFO order -> total order -> causal order
4. Asynchronize transmission

# What's the outside? Who is calling the service?

1. Replication part: multicast, tcp, total causal order
2. Gossip protocol: multicast, udp, no order requirement
3. Data handler part: unicast, tcp, total causal order

# What's the hierachal structual of service?

Unicast Module:

---------------
Basic Interface
---------------
Unicast
---------------
TCP/UDP
---------------

Multicast Module:

---------------
Basic Interface
---------------
Ordering
---------------
Unicast | Unicast | ...
---------------
TCP/UDP
---------------

# What's the Interface? Is it high coersion and low coupling?

## data field
tid: integer, timestamp id
lastTid: integer, last msg timestamp
kind: user-defined
data: interface{}

## Basic Interface
NewCommunication(ORDER_SENSITIVE, REALIBILITY, TCP ? UDP) -> may error if too many communication
  - if unicast: by default fifo+causal
  - if multicast: ORDER_SENSITIVE -> fifo+total order implies causal, otherwise -> fifo
Send(destination(s), data) -> instantly return, service handling communication asynchronizely
Receive() -> wait to get all return data, service doesn't guarantee response unique

## Ordering, three unicasts
prepare -> confirm
prepare: data: {kind=prepare} -> receive data {kind=ackPrepare}
confirm: data: {kind=confirm} -> receive data {kind=ackConfirm}

## Multicast
If not reliable, just iteratively unicast send
If reliable, use piggyback send or re-multicast
data:{kind=multicast, data={..., group member}}

## Unicast
TCP/UDP communication
Dail
Listen
Accept

# Implement. Focus on only necessary features.
