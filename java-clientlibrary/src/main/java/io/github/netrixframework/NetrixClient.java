package io.github.netrixframework;

import io.github.netrixframework.comm.*;
import io.github.netrixframework.timeouts.Timer;
import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.*;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.codec.http.*;
import io.github.netrixframework.timeouts.Timeout;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;

import java.io.IOException;
import java.util.HashMap;
import java.util.Vector;

public class NetrixClient extends Thread {
    NettyServer server;
    NetrixCaller client;
    Timer timer;
    DirectiveExecutor executor;
    NetrixClientConfig netrixClientConfig;
    MessageHandler messageHandler;

    Channel serverChannel;

    Counter counter;

    public NetrixClient(NetrixClientConfig c, DirectiveExecutor executor) {
        this.counter = new Counter();
        this.netrixClientConfig = c;
        this.timer = new Timer();
        this.client = new NetrixCaller(c);
        this.executor = executor;
        this.messageHandler = new MessageHandler(this.client, netrixClientConfig.replicaID);
        initServer();

        try {
            this.client.register();
        } catch (IOException ignored) {

        }
    }

    @Override
    public void run() {
        EventLoopGroup bossGroup = new NioEventLoopGroup();
        EventLoopGroup workerGroup = new NioEventLoopGroup();

        try {
            ServerBootstrap b = new ServerBootstrap();
            b.group(bossGroup, workerGroup)
                .channel(NioServerSocketChannel.class)
                .handler(new LoggingHandler(LogLevel.INFO))
                .childHandler(new ChannelInitializer<SocketChannel>() {
                    @Override
                    public void initChannel(SocketChannel ch) throws Exception {
                        ChannelPipeline p = ch.pipeline();
                        p.addLast(new HttpRequestDecoder());
                        p.addLast(new HttpObjectAggregator(1048576));
                        p.addLast(new HttpResponseEncoder());
                        p.addLast(server);
                    }
                })
                .childOption(ChannelOption.SO_KEEPALIVE, true);;
            this.serverChannel = b.bind(netrixClientConfig.clientServerAddr, netrixClientConfig.clientServerPort).sync().channel();
            this.serverChannel.closeFuture().sync();
        } catch (Exception ignored) {
        } finally {
            workerGroup.shutdownGracefully();
            bossGroup.shutdownGracefully();
        }
    }

    public void stopClient() {
        try {
            this.serverChannel.close();
        } catch (Exception ignored) {

        }
    }

    private void initServer() {
        NettyRouter router = new NettyRouter();

        Route messagesRoute = new Route("/message");
        messagesRoute.post(messageHandler);
        router.addRoute(messagesRoute);

        Route directiveRoute = new Route("/directive");
        directiveRoute.post(new DirectiveHandler(this.executor));
        router.addRoute(directiveRoute);

        Route timeoutRoute = new Route("/timeout");
        timeoutRoute.post(new TimeoutHandler(timer, client));
        router.addRoute(timeoutRoute);

        this.server = new NettyServer(router);
    }

    private String nextID(String from, String to) {
        counter.incr();
        return String.format("%s_%s_%d", from, to, counter.getValue());
    }

    public Vector<Message> getMessages() {
        return messageHandler.getMessages();
    }

    public void sendMessage(Message message) throws IOException {
        String messageID = nextID(netrixClientConfig.replicaID, message.getTo());
        message.setFrom(netrixClientConfig.replicaID);
        message.setId(messageID);

        client.sendMessage(message);

        HashMap<String, String> params = new HashMap<String, String>();
        params.put("message_id", messageID);
        Event sendEvent = new Event("MessageSend", params);
        sendEvent.setReplicaID(netrixClientConfig.replicaID);
        client.sendEvent(sendEvent);
    }

    public void setReady() throws IOException {
        this.client.setReady();
    }

    public void unsetReady() throws IOException {
        this.client.unsetReady();
    }

    public void sendEvent(String type, HashMap<String, String> params) throws IOException{
        Event event = new Event(type, params);
        event.setReplicaID(netrixClientConfig.replicaID);
        client.sendEvent(event);
    }

    public void sendEvent(Event event) throws IOException {
        event.setReplicaID(netrixClientConfig.replicaID);
        client.sendEvent(event);
    }

    public void startTimeout(Timeout timeout) {
        if(timer.addTimeout(timeout)) {
            HashMap<String, String> params = new HashMap<String, String>();
            params.put("type", timeout.key());
            params.put("duration", timeout.getDuration().toString());
            Event event = new Event("TimeoutStart", params);
            // TODO: complete this
        }
    }
}
