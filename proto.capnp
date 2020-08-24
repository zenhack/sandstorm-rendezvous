@0xe24f05f123e395af;

using Go = import "/go.capnp";
$Go.package("main");
$Go.import("zenhack.net/go/sandstorm-rendezvous");

using Ip = import "/ip.capnp";
using Util = import "/util.capnp";

struct PortInfo {
  name @0 :Text;
}

interface LocalNetwork {
  bind @0 (info :PortInfo, port :Ip.TcpPort) -> (handle :Util.Handle);
  resolve @1 (name :Text) -> (port :Ip.TcpPort);
}
