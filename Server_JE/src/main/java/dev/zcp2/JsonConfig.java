package dev.zcp2;

import com.google.gson.Gson;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import lombok.Getter;
import lombok.SneakyThrows;

import java.io.*;
import java.nio.charset.StandardCharsets;

public class JsonConfig {
    public static final String FOLDER = "config";
    @Getter
    private final String name;
    @Getter
    private final File file;
    @Getter
    private JsonObject json;

    static {
        //noinspection ResultOfMethodCallIgnored
        new File(System.getProperty("user.dir"), FOLDER).mkdirs();
    }

    @SneakyThrows
    public JsonConfig(String name, boolean outputDefault) {
        this.name = name;
        file = new File(System.getProperty("user.dir"), FOLDER + "/" + name + ".json");
        if(file.createNewFile()) {
            try(OutputStream out = new FileOutputStream(file); InputStream in = getClass().getClassLoader().getResourceAsStream(name + ".json")) {
                if(outputDefault && in != null)
                    in.transferTo(out);
            }
        }
    }

    @SneakyThrows
    public void load() {
        json = JsonParser.parseReader(
                new InputStreamReader(
                        new FileInputStream(file)
                )
        ).getAsJsonObject();
    }

    @SneakyThrows
    public void save() {
        FileOutputStream out = new FileOutputStream(file);
        out.write(new Gson().toJson(json).getBytes(StandardCharsets.UTF_8));
        out.close();
    }
}
