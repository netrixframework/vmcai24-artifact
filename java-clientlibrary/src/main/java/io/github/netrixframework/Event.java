package io.github.netrixframework;

import com.google.gson.Gson;
import io.github.netrixframework.comm.GsonHelper;

import java.util.HashMap;

public class Event {
    private String type;
    private String replica;
    private HashMap<String, String> params;
    private long timestamp;

    public Event(String type, HashMap<String, String> params) {
        this.type = type;
        this.params = params;
        this.timestamp = System.currentTimeMillis();
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getReplicaID() {
        return replica;
    }

    public void setReplicaID(String replicaID) {
        this.replica = replicaID;
    }

    public HashMap<String, String> getParams() {
        return params;
    }

    public void setParams(HashMap<String, String> params) {
        this.params = params;
    }

    public long getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(long timestamp) {
        this.timestamp = timestamp;
    }

    public String toJsonString() {
        Gson gson = GsonHelper.gson;
        return gson.toJson(this);
    }

    public static Event fromJsonString(String jsonString) {
        Gson gson = GsonHelper.gson;
        return gson.fromJson(jsonString, Event.class);
    }
}
