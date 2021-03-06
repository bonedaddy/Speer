### Motivation and FAQ

**Speer** is a network discrete event Simulator for peer-to-peer network modeling.
It combines event-driven simulations with cycle-based concepts and
allows parallelization by taking advantage of Go’s concurrency features. It aims
to expose clean interfaces which hide the underlying complexity and ensures
correctness via progress and safety properties.

- Why peer-to-peer?

If you quickly google "peer-to-peer networking research" you will probably find
old links. Indeed, [file sharing](https://en.wikipedia.org/wiki/File_sharing)
traffic make up only of [3\% of Internet's upstream](https://torrentfreak.com/bittorrent-traffic-is-not-dead-its-making-a-comeback-180926/) and
system architectures like [Dynamo](https://en.wikipedia.org/wiki/Amazon_DynamoDB)
are less prevalent today. So, why would an open-source project center itself on
a research topic that is rapidly diminishing in popularity since 2008?

**These systems are astonishingly efficient and, arguably, elegant**;
[BitTorrent](https://en.wikipedia.org/wiki/BitTorrent) has been proved to be an
[optimal solution](https://dl.acm.org/citation.cfm?id=1064215) for distributing
data across multiple computers, while DynamoDB being a success project with more
than [100,000 AWS customers](https://aws.amazon.com/dynamodb/). Moreover, around
2003, a large number of not so mainstream overlay designs such as Koorde,
Hypercube, Viceroy, Gia or CAN showed interesting and challenging work, but
many of them produced ideas for a very narrow field of applications. We aim to change this!

- For who?

**Speer is made for students, researchers and hobbyists. It's goal is to allow
them to easily implement, simulate and study peer to peer networks.** This aligns
with the goal of making the internet more decentralized and bringing peer to peer
systems to their former glory.

- Why Go?

Go has memory safety, garbage collection, structural typing and supports a
large number of concurrent routines. This makes it a good candidate for
implementing a discrete event simulator for a large number of nodes. Moreover,
we take advantage of the ease of running on multiple cores.

- Other language support?

We plan to support the languages Python and Elixir.

- How does Speer model the network?

Speer aims to study peer-to-peer systems on top of overlays. An overlay network is a computer network that is built on top of another network. We place the networking interface at the application layer on the [OSI model](https://en.wikipedia.org/wiki/OSI_model). However, we model latency distribution at the Network layer, while the throughput assumes an optimal Transport layer.

- Where can I find more about the implementation of Speer?

The [architecture](architecture.md) section describes the general design and implementation details the most important components of Speer.
