package bftsmart.communication.server;

import bftsmart.communication.SystemMessage;
import bftsmart.consensus.messages.ConsensusMessage;
import bftsmart.reconfiguration.ServerViewController;
import bftsmart.reconfiguration.util.NetrixConfiguration;
import bftsmart.tom.ServiceReplica;
import bftsmart.tom.leaderchange.LCMessage;

import com.google.gson.Gson;
import io.github.netrixframework.comm.Message;
import io.netty.util.CharsetUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import io.github.netrixframework.*;
import javax.crypto.SecretKey;
import java.io.*;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Vector;
import java.util.concurrent.LinkedBlockingQueue;

public class NetrixCommunicationLayer extends CommunicationLayer{
    private Logger logger = LoggerFactory.getLogger(this.getClass());
    private LinkedBlockingQueue<SystemMessage> inQueue;
    private NetrixClient client;

    private boolean doWork = true;

    public NetrixCommunicationLayer(ServerViewController controller,
                                    LinkedBlockingQueue<SystemMessage> inQueue,
                                    ServiceReplica replica) {
        this.inQueue = inQueue;

        NetrixConfiguration c = controller.getStaticConf();

        NetrixClientConfig config = new NetrixClientConfig(
                Integer.toString(replica.getId()),
                c.getNetrixAddr(),
                c.getNetrixClientServerAddr(),
                c.getNetrixClientServerPort(),
                c.getNetrixClientAdvAddr(),
                new HashMap<String, String>()
        );
        this.client = NetrixClientSingleton.init(config, replica);
    }

    public SecretKey getSecretKey(int id) {
        return null;
    }

    public void updateConnections() {}

    private SystemMessage unMarshalMessage(String type, byte[] data) throws Exception{
        SystemMessage sm;
        Gson gson = GsonHelper.gson;
        switch(type) {
            case "PROPOSE":
            case "ACCEPT":
            case "WRITE":
                String jsonString = new String(data, CharsetUtil.UTF_8);
                logger.info("Received consensus message: "+jsonString);
                sm = gson.fromJson(jsonString, ConsensusMessage.class);
                break;
            case "STOP":
            case "STOPDATA":
            case "SYNC":
            case "TRIGGER_LC_LOCALLY":
                jsonString = new String(data, CharsetUtil.UTF_8);
                logger.info("Received leader change message: "+jsonString);
                sm = gson.fromJson(jsonString, LCMessage.class);
                break;
            case "other":
            default:
                sm = (SystemMessage) (new ObjectInputStream(new ByteArrayInputStream(data))
                .readObject());
        }
        sm.authenticated = true;
        return sm;
    }

    private byte[] marshalMessage(SystemMessage sm) {
        if(sm instanceof ConsensusMessage) {
            logger.debug("Sending a consensus message");
            ConsensusMessage cm = (ConsensusMessage) sm;
            Gson gson = GsonHelper.gson;
            String jsonString = gson.toJson(cm);
            logger.info("Sending a consensus message: "+jsonString);
            return jsonString.getBytes(StandardCharsets.UTF_8);
        } else if(sm instanceof LCMessage) {
            logger.debug("Sending a leader change message");
            LCMessage cm = (LCMessage) sm;
            Gson gson = GsonHelper.gson;
            String jsonString = gson.toJson(cm);
            logger.info("Sending a leader change message: "+jsonString);
            return jsonString.getBytes(StandardCharsets.UTF_8);
        } else {
            ByteArrayOutputStream bOut = new ByteArrayOutputStream(248);
            try {
                new ObjectOutputStream(bOut).writeObject(sm);
            } catch (IOException e) {
                logger.error("Failed to serialize message", e);
            }
            return bOut.toByteArray();
        }
    }

    public final void send(int[] targets, SystemMessage sm, boolean useMAC) {
        byte[] data = marshalMessage(sm);

        String messageType = "other";
        if(sm.getClass() == ConsensusMessage.class) {
            messageType = ((ConsensusMessage) sm).getPaxosVerboseType();
        } else if (sm.getClass() == LCMessage.class) {
            messageType = ((LCMessage) sm).getTypeString();
        }

        for(int t: targets) {
            try {
                client.sendMessage(new Message(
                        Integer.toString(t),
                        messageType,
                        data
                ));
            } catch (IOException e) {
                logger.error("Failed to send message", e);
            }
        }
    }

    public void shutdown() {
        doWork = false;
        // client.stopClient();
    }

    public void joinViewReceived() {
    }

    @Override
    public void run() {
        try {
            if (!client.isAlive()) {
                client.start();
            }
            client.setReady();
        } catch (Exception ignored) {

        }
        while(doWork) {
            Vector<Message> messages = client.getMessages();
            for(Message m: messages) {
                try {
                    SystemMessage sm = unMarshalMessage(m.getType(), m.getData());
                    if(!inQueue.offer(sm)) {
                        logger.warn("Inqueue full (message from " + sm.getSender() + " discarded).");
                    }
                } catch (Exception e) {
                    logger.info("Error parsing message: "+e.getMessage());
                }
            }
        }
    }
}
