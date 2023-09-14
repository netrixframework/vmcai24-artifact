package io.github.netrixframework.comm;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.handler.codec.http.*;
import io.netty.util.CharsetUtil;
import io.github.netrixframework.Event;

import java.io.IOException;
import java.util.HashMap;
import java.util.Vector;
import java.util.concurrent.locks.ReentrantLock;

import static io.netty.handler.codec.http.HttpHeaderNames.CONTENT_TYPE;
import static io.netty.handler.codec.http.HttpHeaderValues.APPLICATION_JSON;

public class MessageHandler implements Handler{

    private String replicaID;
    private Vector<Message> messages;
    private final ReentrantLock lock = new ReentrantLock();
    private NetrixCaller client;

    public MessageHandler(NetrixCaller client, String replicaID) {
        this.messages = new Vector<Message>();
        this.replicaID = replicaID;
        this.client = client;
    }

    @Override
    public FullHttpResponse handle(FullHttpRequest req) {
        try {
            Message m = getMessageFromReq(req);
            lock.lock();
            messages.add(m);
            lock.unlock();

            HashMap<String, String> params = new HashMap<String, String>();
            params.put("message_id", m.getId());
            Event receiveEvent = new Event(
                "MessageReceive",
                params
            );
            receiveEvent.setReplicaID(replicaID);
            client.sendEvent(receiveEvent);

            return new DefaultFullHttpResponse(
                    req.protocolVersion(),
                    HttpResponseStatus.OK,
                    Unpooled.copiedBuffer("Ok", CharsetUtil.UTF_8)
            );
        } catch (Exception e) {
            return new DefaultFullHttpResponse(
                    req.protocolVersion(),
                    HttpResponseStatus.INTERNAL_SERVER_ERROR
            );
        }
    }

    private Message getMessageFromReq(FullHttpRequest req) throws IOException {
        ByteBuf content = req.content();
        if(content == null || content.readableBytes() <= 0){
            throw new IOException("empty request");
        }
        if(!req.headers().get(CONTENT_TYPE).equals(APPLICATION_JSON.toString())) {
            throw new IOException("not a json request");
        }
        return Message.fromJsonString(content.toString(CharsetUtil.UTF_8));
    }

    public Vector<Message> getMessages() {
        Vector<Message> result = new Vector<Message>();
        lock.lock();
        result.addAll(messages);
        messages.clear();
        lock.unlock();
        return result;
    }
}
