package dev.zcp2;

public final class Main {
    private static final JsonConfig config = new JsonConfig("config", true);
    public static void main(String[] args) {
        config.load();
    }
}
