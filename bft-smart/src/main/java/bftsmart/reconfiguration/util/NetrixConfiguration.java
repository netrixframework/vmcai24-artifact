package bftsmart.reconfiguration.util;

import bftsmart.tom.util.KeyLoader;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class NetrixConfiguration extends TOMConfiguration{

    private Logger logger = LoggerFactory.getLogger(this.getClass());
    private String netrixAddr;
    private int netrixAddrPort;
    private String netrixClientServerAddr;
    private int netrixClientServerPort;
    private String netrixClientAdvAddr;
    private int netrixClientAdvPort;
    private boolean useNetrix;

    public NetrixConfiguration(int processId, String configHome, KeyLoader loader) {
        super(processId, configHome, loader);
    }

    public NetrixConfiguration(int processId, KeyLoader loader) {
        super(processId, loader);
    }

    @Override
    protected void init() {
        super.init();
        try {

            String s = (String) this.configs.remove("system.communication.useNetrix");
            if (s == null) {
                this.useNetrix = false;
            } else {
                this.useNetrix = Boolean.parseBoolean(s);
            }

            s = (String) this.configs.remove("system.communication.netrixAddress");
            if (s == null) {
                this.netrixAddr = "";
            } else {
                this.netrixAddr = s;
            }

            s = (String) this.configs.remove("system.communication.netrixAddress");
            if (s == null) {
                this.netrixAddrPort = 7074;
            } else {
                this.netrixAddrPort = Integer.parseInt(s);
            }

            s = (String) this.configs.remove("system.communication.netrixClientServerAddress");
            if (s == null) {
                this.netrixClientServerAddr = "127.0.0.1";
            } else {
                this.netrixClientServerAddr = s;
            }

            s = (String) this.configs.remove( "system.communication.netrixClientServerPort");
            if (s == null) {
                this.netrixClientServerPort = 8080+processId;
            } else {
                this.netrixClientServerPort = Integer.parseInt(s);
            }

            s = (String) this.configs.remove("system.communication.netrixClientAdvAddress");
            if (s == null) {
                this.netrixClientAdvAddr = this.netrixClientServerAddr;
            } else {
                this.netrixClientAdvAddr = s;
            }

            s = (String) this.configs.remove( "system.communication.netrixClientAdvAddressPort");
            if (s == null) {
                this.netrixClientAdvPort = this.netrixClientServerPort;
            } else {
                this.netrixClientAdvPort = Integer.parseInt(s);
            }

        } catch (Exception e) {
            logger.error("Could not parse system configuration file",e);
        }
    }

    public String getNetrixAddr() {
        return String.format("%s:%d",this.netrixAddr, this.netrixAddrPort);
    }

    public String getNetrixClientServerAddr() {
        return this.netrixClientServerAddr;
    }

    public String getNetrixClientAdvAddr() {
        return String.format("%s:%d",this.netrixClientAdvAddr, this.netrixClientAdvPort);
    }

    public int getNetrixClientServerPort() {
        return this.netrixClientServerPort;
    }

    public boolean useNetrix() {
        return useNetrix;
    }
}
