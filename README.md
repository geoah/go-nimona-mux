# Deprecation Notice

__Warning__: This library is now deprecated and has been replaced by a [fabric](https://github.com/nimona/go-nimona-fabric) protocol.

## Stream Multiplexer

Mux allows splitting a single TCP connection to multiple bi-directional
streams. It does so by appending some basic stream information such as 
stream id, and content length.

