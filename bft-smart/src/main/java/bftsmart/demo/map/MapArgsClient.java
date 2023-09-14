package bftsmart.demo.map;

public class MapArgsClient {

    public static void main(String[] args) {
        if (args.length < 2) {
            System.out.println("Usage: demo.map.MapInteractiveClient <client id> <command>\n\n");
            System.out.println("Command help");
            System.out.println("get <key>");
            System.out.println("set <key> <value>");
            System.out.println("delete <key>");

            return;
        }

        int clientId = Integer.parseInt(args[0]);
        MapClient<String, String> client = new MapClient<>(clientId);

        switch (args[1]) {
            case "get":
                get(args, client);
                break;
            case "set":
                set(args, client);
                break;
            case "delete":
                delete(args, client);
                break;
            default:
                System.out.println("Invalid command. Use one of get, set or delete.");
        }
        client.close();
    }

    public static void get(String[] args, MapClient<String, String> client) {
        if (args.length != 3) {
            System.out.println("Usage: get <key>");
            return;
        }
        System.out.println(client.get(args[2]));
    }

    public static void set(String[] args, MapClient<String, String> client) {
        if (args.length != 4) {
            System.out.println("Usage: set <key> <value>");
            return;
        }
        System.out.println(client.put(args[2], args[3]));
    }

    public static void delete(String[] args, MapClient<String, String> client) {
        if (args.length != 3) {
            System.out.println("Usage: delete <key>");
            return;
        }
        System.out.println(client.remove(args[2]));
    }
}
