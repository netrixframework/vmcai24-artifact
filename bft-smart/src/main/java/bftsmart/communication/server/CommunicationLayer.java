package bftsmart.communication.server;

import bftsmart.communication.SystemMessage;

import javax.crypto.SecretKey;

public abstract class CommunicationLayer extends Thread {
    public abstract SecretKey getSecretKey(int id);
    public abstract void updateConnections();
    public abstract void send(int[] targets, SystemMessage sm, boolean useMAC);
    public abstract void shutdown();
    public abstract void joinViewReceived();
}
