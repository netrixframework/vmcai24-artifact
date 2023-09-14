package bftsmart.communication.client.netty;

import bftsmart.communication.client.netty.NettyClientServerSession;
import bftsmart.communication.server.GsonHelper;
import bftsmart.reconfiguration.ViewController;
import bftsmart.tom.core.messages.TOMMessage;
import com.google.gson.Gson;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.ByteToMessageDecoder;
import io.netty.util.CharsetUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.locks.ReentrantReadWriteLock;

public class NettyJsonMessageDecoder extends ByteToMessageDecoder {

    private Logger logger = LoggerFactory.getLogger(this.getClass());

    private boolean isClient;
    private ConcurrentHashMap<Integer, NettyClientServerSession> sessionTable;
    private ViewController controller;
    private ReentrantReadWriteLock rl;

    public NettyJsonMessageDecoder(boolean isClient,
                                   ConcurrentHashMap<Integer, NettyClientServerSession> sessionTable,
                                   ViewController controller,
                                   ReentrantReadWriteLock rl) {
        this.isClient = isClient;
        this.sessionTable = sessionTable;
        this.controller = controller;
        this.rl = rl;
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, ByteBuf in, List<Object> out) throws Exception {
        // read integer <size> in bytes and then the json blob
        if (in.readableBytes() < Integer.BYTES) {
            return;
        }

        int len = in.getInt(in.readerIndex());
        if (in.readableBytes() < len + Integer.BYTES) {
            return;
        }

        in.skipBytes(Integer.BYTES);
        byte[] data = new byte[len];
        in.readBytes(data);

        Gson gson = GsonHelper.gson;
        TOMMessage sm = gson.fromJson(new String(data, CharsetUtil.UTF_8), TOMMessage.class);

        if (!isClient) {
            rl.readLock().lock();
            boolean hasSender = sessionTable.containsKey(sm.getSender());
            rl.readLock().unlock();

            if (!hasSender) {
                NettyClientServerSession cs = new NettyClientServerSession(
                        ctx.channel(),
                        sm.getSender());
                rl.writeLock().lock();
                sessionTable.put(sm.getSender(), cs);
                logger.debug("Active clients: " + sessionTable.size());
                rl.writeLock().unlock();
            }
        }
        out.add(sm);
    }
}
