package bftsmart.communication.client.netty;

import bftsmart.communication.client.netty.NettyClientServerSession;
import bftsmart.communication.server.GsonHelper;
import bftsmart.tom.core.messages.TOMMessage;
import com.google.gson.Gson;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToByteEncoder;
import io.netty.util.CharsetUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.locks.ReentrantReadWriteLock;

public class NettyJsonMessageEncoder extends MessageToByteEncoder<TOMMessage> {

    private Logger logger = LoggerFactory.getLogger(this.getClass());
    private boolean isClient;
    private ConcurrentHashMap<Integer, NettyClientServerSession> sessionTable;
    private ReentrantReadWriteLock rl;

    public NettyJsonMessageEncoder(boolean isClient,
                                   ConcurrentHashMap<Integer, NettyClientServerSession> sessionTable,
                                   ReentrantReadWriteLock rl) {
        this.isClient = isClient;
        this.sessionTable = sessionTable;
        this.rl = rl;
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, TOMMessage msg, ByteBuf out) throws Exception {
        Gson gson = GsonHelper.gson;
        String jsonString = gson.toJson(msg);
        byte[] data = jsonString.getBytes(CharsetUtil.UTF_8);

        logger.info("serialized TOMMessage: "+jsonString);

        out.writeInt(data.length);
        out.writeBytes(data);

        ctx.flush();
    }
}
