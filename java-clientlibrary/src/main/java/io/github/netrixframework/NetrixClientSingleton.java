package io.github.netrixframework;

public class NetrixClientSingleton {
    private static NetrixClient client = null;

    public static NetrixClient init(NetrixClientConfig c, DirectiveExecutor executor) {
        if(client == null) {
            client = new NetrixClient(c, executor);
        }
        return client;
    }

    public static NetrixClient getClient() {
        return client;
    }
}
